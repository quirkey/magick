package main

import (
	"log"
	"os"
	"github.com/quirkey/magick"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		log.Print("Please supply an input and output filename e.g. ./examples foo.jpg bar.jpg")
		os.Exit(1)
	}
	start := time.Now()
	input := os.Args[1]
	log.Printf("Reading from file %s", input)
	output := os.Args[2]

	image, err := magick.NewFromFile(input)
	if err != nil {
		log.Printf("Error reading from file %s", err.Error())
		os.Exit(1)
	}
	log.Print("Transforming")
	err = image.Crop("100x100")
	if err != nil {
		log.Print("Problem with transforming")
		os.Exit(1)
	}
	err = new_image.Shadow("#000", 75, 5, 2, 2)
	if err != nil {
		log.Print("Problem with transforming")
		os.Exit(1)
	}
	err = image.FillBackgroundColor("#333")
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
	log.Printf("Done. took %v\n", end.Sub(start))
}
