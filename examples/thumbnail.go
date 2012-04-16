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
	ok := image.Transform("", "100x100")
        if !ok {
          log.Print("Problem with transforming")
          os.Exit(1)
        }
	log.Printf("Writing to %s", output)
	ok, err = image.ToFile(output)
	log.Printf("Wrote to %s %b", output, ok)
	if err != nil {
		log.Printf("Error outputing file: %s", err.Error())
		os.Exit(1)
	}
	end := time.Now()
	log.Printf("Done. took %v\n", end.Sub(start))
}
