package main

import (
	"github.com/quirkey/magick"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		log.Print("Please supply an input and output filename e.g. ./examples foo.jpg bar.jpg")
		os.Exit(1)
	}
	input := os.Args[1]
	output := os.Args[2]
	log.Printf("Reading from file %s, writing to file %s", input, output)

	start := time.Now()

	image, err := magick.NewFromFile(input)
	log.Printf("Loading image took %v\n", time.Now().Sub(start))
	start = time.Now()
	if err != nil {
		log.Printf("Error reading from file %s", err.Error())
		os.Exit(1)
	}

	image.Quality(100)
	_ = image.Progressive()
	log.Print("Transforming")
	log.Printf("size: %d %d", image.Width(), image.Height())
	err = image.Resize("2000x2000!")
	log.Printf("Transforming image took %v\n", time.Now().Sub(start))
	start = time.Now()
	if err != nil {
		log.Print("Problem with transforming")
		os.Exit(1)
	}

	log.Printf("Writing to %s", output)
	err = image.ToFile(output)
	log.Printf("Writing image took %v\n", time.Now().Sub(start))
	log.Printf("Wrote to %s", output)

	image.Destroy()
}
