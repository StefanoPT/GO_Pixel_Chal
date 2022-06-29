To run the application, run:

$ go run ./cmd/main.go -dir ./Gold -img main.raw

-dir flag points to the directory of the images to compare.
-img is the name of the main image inside that directory. Note, in the example above ./Gold/main.raw is an incorrect value for the img flag

Pixel Challange
Learning Outcomes
Using a non-trivial exercise, we aim to accomplish the following:
Unit + Benchmark Testing to a high coverage %
Using pprof to profile CPU and Memory Usage
Iteratively improving performance by identifying bottlenecks
Good, clean Git behaviour
All whilst maintaining good, clean code that can be easily read and conforms to Idiomatic Go standards
The Challenge
Given a selection of images, and choosing a particular reference one, find the 3 “closest” images to that reference image.
How “close” an image is determined by the number of pixels that match precisely between the two images - both in value and location.
% Closeness = (Number of Pixels Matching) / (Number of Pixels Total) * 100
The output should be a list of images, and their to the reference either as full percentage (97%) or decimal (0.97)
The Images
Images are purely random noise. They are 1024 pixels wide and 1024 pixels high.
Each pixel is encoded as 3 bytes of data, representing the Red, Green and Blue Channels.
For example: the first byte could be represented as [255, 255, 255] which would be a white pixel. Another pixel could be [0, 0, 0] which would be a black pixel.