Algo
====

As a means of learning Go, I'm writing a few small programs to create interesting mathematical animations. This requires learning Go's type system, navigating the standard library, and playing with concurrency, so overall it's a nice tour of the language.

The shared drawing functions are in the algo folder. Each other folder contains a program that uses the drawing functions to produce an animation. The current programs are:

* `polyfill`: rotate an array of three-dimensional polygons around a given axis
* `rotation`: rotate a three-dimensional wireframe (i.e. an array of vector pairs) around a given axis
* `parametric`: animate the effect of changing the t-step value of a given parametric function

Installation
------------

Just run `go get github.com/lukechampine/algo/...` and binaries will be installed in your GOPATH.

This project was inspired by http://www.pheelicks.com/2013/10/intro-to-images-in-go-part-1
