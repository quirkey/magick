package main

import (
  "log"
  "os"
  "time"
  "fmt"
  "quirkey/magick"
)

func main() {
  start := time.Now()
  input := os.Args[1]
  log.Printf("Reading from file %s", input)
  output := os.Args[2]

  image, ok := magick.NewFromFile(input)
  if !ok {
    log.Printf("Error reading from file")
    os.Exit(1)
  }
  // log.Print("Transforming")
  // image.Transform("100>x", "100x100")
  log.Printf("Writing to %s", output)
  blob, ok := image.ToBlob()
  fmt.Print(string(blob))
  if !ok {
    log.Printf("Error reading from file")
    os.Exit(1)
  }
  end := time.Now()
  log.Printf("Done. took %v\n", end.Sub(start))
}
