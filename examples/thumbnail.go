package main

import (
	"log"
	"os"
	"quirkey/magick"
	"time"
)

func main() {
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
	var new_image *magick.MagickImage
	// new_image, err := image.Thumbnail(100, 100)
	// if err != nil {
	// 	log.Print("Problem with transforming")
	// 	os.Exit(1)
	// }
	new_image, err = image.Shadow("#000", 75, 5, 2, 2)
	if err != nil {
		log.Print("Problem with transforming")
		os.Exit(1)
	}
	new_image, err = new_image.FillBackgroundColor("#333")
	if err != nil {
		log.Print("Problem setting background")
		os.Exit(1)
	}
	log.Printf("Writing to %s", output)
	err := new_image.ToFile(output)
	if err != nil {
		log.Printf("Error outputing file: %s", err.Error())
		os.Exit(1)
	}
	log.Printf("Wrote to %s %b", output, ok)
	end := time.Now()
	log.Printf("Done. took %v\n", end.Sub(start))
}
