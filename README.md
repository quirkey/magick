# magick

magick implements image manipulation routines based on the 
ImageMagick MagickCore C library in Go. It is an opinionated high level 
wrapper around the proven ImageMagick lib.

## Why

### ImageMagick is Magic

Though Go's stdlib includes utilities for working with images, ImageMagick is one of the industry standards, not only in its relative ease-of-use for simple operations (like thumbnailing) but also in it's smart anti-aliasing and other proven techniques. We wanted to utilize the fast and conccurent environment that Go provides with the known and reliable output of what we were already seeing with ImageMagick.

### Simple, High-level Operations

There are other libraries that wrap ImageMagick or related libraries, but we were looking/aiming for something that had simple functions that handled the most common operations, that we wanted needed in our web applications: Thumbnailing (Resizing/Cropping), adding shadows, converting to jpg, etc.

### Works with BLOBs

One of the major bottlenecks we saw in our previous image pipeline was the reading and writing from disk (IO). We wanted to build our new pipeline around the idea that images could be read and manipulated from memory (or non-disk storage i.e. a database). This allows an image to be transformed without ever touching disk.

## Usage

With files:

``` go
image, err := magick.NewFromFile("input.png")
defer image.Destroy()
err = image.Resize("400x200")
err = image.Shadow("#F00", 255, 5, 2, 2)
err = image.FillBackgroundColor("#00F")
err = image.ToFile("output.jpg")
```

With BLOBs:

``` go
image, err := magick.NewFromBlob(a_byteslice, "png")
defer image.Destroy()
err = image.Resize("400x200")
err = image.Shadow("#F00", 255, 5, 2, 2)
err = image.FillBackgroundColor("#00F")
a_new_byteslice, err = image.ToBlob("jpg")
```

For full API see the [API docs](http://godoc.org/github.com/quirkey/magick)

## Gotchas/Known Issues

magick has been thorougly tested and is memory-leak free as long as you always `Destroy()` MagickImage's after you no longer need them.

Internally, MagickCore can be used concurrently without issues, though weve observed crashes/issues with concurrent usage when ImageMagick is compiled with OpenMP on OS X (this happens to be the default with homebrew). Mileage may vary.

## TODO

- Text rendering
- compositing two images

## Dependencies

magick depends on ImageMagick and specifically the [MagickCore](http://www.imagemagick.org/script/magick-core.php) C library. In most linux environments this is included in the ImageMagick-devel packages (e.g. `yum install ImageMagick-devel` or `sudo aptitude install ImageMagick-devel`).

## More

* [API Docs](http://godoc.org/github.com/quirkey/magick)
* [Examples](examples/)

## Who

magick is written and maintained by Aaron Quint ([@aq](http://twitter.com/aq)) and Mike Bernstein ([@mrb_bk](http://twitter.com/mrb_bk))

