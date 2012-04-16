package main

import (
  "log"
  "os"
  "time"
  "quirkey/magick"
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
  // log.Print("Transforming")
  // image.Transform("100>x", "100x100")
  log.Printf("Writing to %s", output)
  _, err = image.ToFile(output)
  if err != nil {
    log.Printf("Error outputing file: %s", err.Error())
    os.Exit(1)
  }
  end := time.Now()
  log.Printf("Done. took %v\n", end.Sub(start))
}
