package main

import (
	"github.com/quirkey/magick"
	"io/ioutil"
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

	source, _ := ioutil.ReadFile(input)
	image, err := magick.NewFromBlob(source, "jpg")
	log.Printf("Loading image took %v\n", time.Now().Sub(start))
	start = time.Now()
	if err != nil {
		log.Printf("Error reading from file %s", err.Error())
		os.Exit(1)
	}

	image.Quality(50)
	image.Strip()
	image.Progressive()
	image.SetProperty("jpeg:sampling-factor", "4:4:4")
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
	if err != nil {
		log.Printf("Problem with writing %v", err)
		os.Exit(1)
	}
	log.Printf("Writing image took %v\n", time.Now().Sub(start))
	log.Printf("Wrote to %s", output)

	image.Destroy()
}
