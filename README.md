# goinyourface

Detect faces in images using a library called [CCV](https://github.com/liuliu/ccv).

## Note

This code is **most definitely not production ready by any stretch of the imagination.** This code was simply a quick exercise to learn Go and how to bind a C library in two days. Was fun, lots of trial and error.. especially on a mac osx. I really have no idea how to reinstall this... on that note -- when you run it -- you may get a syscall error but the program itself actually worked (used to not be the case, and not worth the effort in repairing since this was educational only).

Don't think it works? I don't blame you, the results returned from the **test commands** are in the `results` folder.

## Requirements

* [fftw](https://github.com/FFTW/fftw3)
* [ccv](https://github.com/liuliu/ccv)
* `go get github.com/lazywei/go-opencv/opencv`

## Test commands

`go run image.go ./nfl.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m`

`go run image.go ./face1.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian`

`go run image.go ./face2.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m`

`go run image.go ./face3.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m`

`go run image.go ./face4.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m`
