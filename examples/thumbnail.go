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
	new_image, err := image.Thumbnail(100, 100)
	if err != nil {
		log.Print("Problem with transforming")
		os.Exit(1)
	}
        new_image, err = new_image.FillBackgroundColor("#333")
        if err != nil {
                log.Print("Problem setting background")
                os.Exit(1)
        }
	// new_image, err = new_image.Shadow(0.8, 0, 2, 2)
	// if err != nil {
	// 	log.Print("Problem with transforming")
	// 	os.Exit(1)
	// }
	log.Printf("Writing to %s", output)
	ok, err := new_image.ToFile(output)
	log.Printf("Wrote to %s %b", output, ok)
	if err != nil {
		log.Printf("Error outputing file: %s", err.Error())
		os.Exit(1)
	}
	end := time.Now()
	log.Printf("Done. took %v\n", end.Sub(start))
}
