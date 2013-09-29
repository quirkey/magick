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

	for i := 0; i < 10; i++ {
		start := time.Now()

		image, err := magick.NewFromFile(input)
		log.Printf("Loading image took %v\n", time.Now().Sub(start))
		if err != nil {
			log.Printf("Error reading from file %s", err.Error())
			os.Exit(1)
		}

		log.Print("Transforming")
		err = image.Resize("400x200")
		if err != nil {
			log.Print("Problem with transforming")
			os.Exit(1)
		}

		err = image.Shadow("#F00", 255, 5, 2, 2)
		if err != nil {
			log.Print("Problem with transforming")
			os.Exit(1)
		}

		err = image.FillBackgroundColor("#00F")
		if err != nil {
			log.Print("Problem setting background")
			os.Exit(1)
		}

		log.Printf("Writing to %s", output)
		err = image.ToFile(output)
		if err != nil {
			log.Printf("Error outputing file: %s", err.Error())
			os.Exit(1)
		}
		log.Printf("Wrote to %s %b", output)

		end := time.Now()
		log.Printf("Done. took %v\n", time.Now().Sub(start))
		image.Destroy()
	}
}
