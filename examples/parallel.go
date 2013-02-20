package main

import (
	"log"
	"os"
	"io/ioutil"
        "fmt"
	"quirkey/magick"
	"time"
)

func main() {
//	start := time.Now()
	input := os.Args[1]
	log.Printf("Reading from file %s", input)
	// times := os.Args[2]
        times := 100
	source, _ := ioutil.ReadFile(input)
        files := make(chan string)
        for i := 0; i < times; i++ {
          go MakeThumbnail(source, i, files)
        }
        for {
          <-files
        }
}

func MakeThumbnail(source []byte, num int, c chan string) {
	start := time.Now()
        output := fmt.Sprintf("tmp/out_%d.png", num)
        log.Printf("Working with %s", output)
	image, err := magick.NewFromBlob(source, "png")
        defer image.Destroy()
	if err != nil {
		log.Printf("Error reading from file %s", err.Error())
		os.Exit(1)
	}
	err = image.Crop("100x100")
	if err != nil {
		log.Print("Problem with transforming")
		os.Exit(1)
	}
	// new_image, err = new_image.Shadow("#000", 75, 5, 2, 2)
	if err != nil {
		log.Print("Problem with transforming")
		os.Exit(1)
	}
	err = image.FillBackgroundColor("#333")
	if err != nil {
		log.Print("Problem setting background")
		os.Exit(1)
	}
	err = image.ToFile(output)
	if err != nil {
		log.Printf("Error outputing file: %s", err.Error())
		os.Exit(1)
	}
	end := time.Now()
	log.Printf("done with %s. took %v", output, end.Sub(start))
        c <- output
}
