package main

import (
  "log"
  "os"
  magick "../magick"
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
  log.Print("Transforming")
  image.Transform("100>x", "100x100")
  log.Printf("Writing to %s", output)
  ok = magick.ToFile(output)
  if !ok {
    log.Printf("Error reading from file")
    os.Exit(1)
  }
  log.Printf("Done. %f")
  end := time.Now()
  log.Printf("took %v\n", end.Sub(start))
}
