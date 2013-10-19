Algo
====

As a means of learning Go, I'm writing a few small programs to create interesting mathematical animations. Currently, each frame is rendered individually and written to a file, and ImageMagick is called to composite the frames into a .gif animation. Hopefully this will be replaced by a less ridiculous approach in the near future.

Current programs:
* `rotation.go`: rotate a three-dimensional figure (i.e. an array of vector pairs) along a given axis
* `parametric.go`: animate the effect of changing the t-step value of a given parametric function
