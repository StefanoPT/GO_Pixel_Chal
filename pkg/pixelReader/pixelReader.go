package pixelReader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const PixelSize = 3

const NotEnoughEquals = -1

var CheckForEnough = false

var ReadingFactor = 1

var ReadingChunck int

var MatrixSize int

type Matrix struct {
	Pixels       []byte
	NumberPixels int
	FileName     string
	FileSize     int64
}

type Result struct {
	File                string
	Percent             float64
	NumberOfEqualPixels int
}

var MainImage Matrix
var Best Result
var Second Result
var Third Result

func updateResults(numEqual int, fname *string) {
	////HereGlobal.inc()
	percent := (float64(numEqual) / float64(MainImage.NumberPixels)) * 100
	if percent <= Third.Percent {
		return
	}
	res := Result{File: *fname, Percent: percent, NumberOfEqualPixels: numEqual}
	if percent > Best.Percent {
		Third = Second
		Second = Best
		Best = res
		return
	}
	if percent > Second.Percent {
		Third = Second
		Second = res
		return
	}
	Third = res
}

func parseFiles(file fs.DirEntry, dir string, ch chan struct{}) {

	defer func(ch chan struct{}) {
		ch <- struct{}{}
	}(ch)

	fn := file.Name()
	fileExtension := filepath.Ext(fn)

	if fileExtension != ".raw" {
		//HereGlobal.inc()
		return
	}

	fname := filepath.Join(dir, fn)

	if fname == MainImage.FileName {
		//HereGlobal.inc()
		return
	}

	fi, err := os.Stat(fname)
	if err != nil {
		//HereGlobal.inc()
		return
	}

	// get the size
	size := fi.Size()
	if size != int64(MainImage.FileSize) {
		return
	}

	data, err := os.ReadFile(fname)

	if err != nil {
		//HereGlobal.inc()
		return
	}

	numEqual := 0

	resultChannel := make(chan int, ReadingFactor)
	readFrom := 0
	readTo := ReadingChunck

	for i := 0; i < ReadingFactor; i++ {
		go parseAndCompareMatrixes(readFrom, readTo, data[readFrom:readTo], resultChannel)
		readFrom += ReadingChunck
		readTo += ReadingChunck
	}

	for i := 0; i < ReadingFactor; i++ {
		numEqual += <-resultChannel
	}

	if numEqual != NotEnoughEquals {
		updateResults(numEqual, &fname)
	}
}

/*func checkIfStillPossible(numberOfEqualPixels, startingPos int) bool {
	availablePixels := (MainImage.NumberPixels - startingPos) / PixelSize
	possibleMaxNPixels := availablePixels + numberOfEqualPixels
	possibleMaxPercent := float64(possibleMaxNPixels) / float64(MainImage.NumberPixels)

	if possibleMaxPercent > Third.Percent {
		return true
	}

	return false
}*/

func comparePixels(r, g, b byte, matrixPos int) bool {
	if r != MainImage.Pixels[matrixPos] {
		return false
	}
	if g != MainImage.Pixels[matrixPos+1] {
		return false
	}
	if b != MainImage.Pixels[matrixPos+2] {
		return false
	}
	return true
}

func checkIfStillPossible(mainMatrixPos, numberOfEqualPixels int) bool {
	leftOverPixels := (MainImage.NumberPixels - mainMatrixPos/PixelSize)
	return (leftOverPixels + numberOfEqualPixels) > Third.NumberOfEqualPixels
}

//parseAndCompareMatrixes takes the starting point in the Main image (the X coordinate) and a few bytes to read and reads them 3 by 3, comparing the resulting RGB to the Main Image
func parseAndCompareMatrixes(mainMatrixPos, readTo int, filePiece []byte, ch chan int) {
	numberOfEqualPixels := 0
	//lenght := readTo - startingPoint
	for i := 0; i < len(filePiece); i += PixelSize {
		if equals := comparePixels(filePiece[i], filePiece[i+1], filePiece[i+2], mainMatrixPos); equals {
			numberOfEqualPixels += 1
		}

		if CheckForEnough && ReadingFactor == 1 {
			if possible := checkIfStillPossible(mainMatrixPos, numberOfEqualPixels); !possible {
				ch <- NotEnoughEquals
				return
			}
		}
		mainMatrixPos += PixelSize

	}
	ch <- numberOfEqualPixels

}

func parseMainMatrix(startingPoint int, filePiece []byte, ch chan struct{}) {
	for i := 0; i < len(filePiece); i += 1 {
		MainImage.Pixels[startingPoint] = filePiece[i]
		startingPoint += 1
	}
	ch <- struct{}{}
}

func makeMainImage(data *[]byte, fn *string) error {
	fileStat, err := os.Stat(*fn)

	if err != nil {
		return err
	}

	size := fileStat.Size()

	MatrixSize = len(*data)
	MainImage = Matrix{Pixels: make([]byte, MatrixSize), NumberPixels: MatrixSize / PixelSize, FileName: *fn, FileSize: size}
	ReadingChunck = MatrixSize / ReadingFactor

	return nil
}

func parseMainImage(filename *string) error {

	data, err := os.ReadFile(*filename)

	if err != nil {
		return err
	}

	err = makeMainImage(&data, filename)

	if err != nil {
		fmt.Println(err)
		return err
	}

	readCh := make(chan struct{}, ReadingFactor)

	readFrom := 0
	readTo := ReadingChunck
	for i := 0; i < ReadingFactor; i++ {
		go parseMainMatrix(readFrom, data[readFrom:readTo], readCh)
		readFrom += ReadingChunck
		readTo += ReadingChunck
	}

	for i := 0; i < ReadingFactor; i++ {
		<-readCh
	}

	return nil
}

func parseImageFiles(directory *string) error {
	imageFiles, err := os.ReadDir(*directory)

	if err != nil {
		return err
	}

	numberOfImages := len(imageFiles)
	parseCh := make(chan struct{}, numberOfImages)

	for _, el := range imageFiles {
		go parseFiles(el, *directory, parseCh)
	}

	for i := 0; i < numberOfImages; i++ {
		<-parseCh
	}

	return nil
}

func PrintTopThreeString() {
	fmt.Printf("BEST: %v\nSECOND: %v\nTHIRD: %v\n", Best, Second, Third)
}

func setUp() {
	Best = Result{Percent: 0.0}
	Second = Result{Percent: 0.0}
	Third = Result{Percent: 0.0}
}

func Run(directory, mainImage *string) {

	setUp()

	imagePath := filepath.Join(*directory, *mainImage)

	err := parseMainImage(&imagePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = parseImageFiles(directory)

	if err != nil {
		return
	}
}
