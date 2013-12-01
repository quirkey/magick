package main

import (
	"fmt"
	//"runtime/debug"
	"github.com/quirkey/magick"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	input := os.Args[1]
	quit := make(chan int)
	log.Printf("Reading from file %s", input)
	// times := os.Args[2]
	times := 50
	source, _ := ioutil.ReadFile(input)
	files := make(chan string)

	for i := 0; i < times; i++ {
		go MakeThumbnail(source, i, files)
	}
	for i := 0; i < times; i++ {
		<-files
	}
	WriteHeapProfile()
	time.Sleep(10000 * time.Millisecond)
	//debug.FreeOSMemory()
	runtime.GC()
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	log.Printf("Alloc: %d", m.Alloc/1024)
	log.Printf("Total Alloc: %d", m.TotalAlloc/1024)
	log.Printf("Sys: %d", m.Sys/1024)
	log.Printf("Heap Sys: %d", m.HeapSys/1024)
	log.Printf("Heap Inuse: %d", m.HeapInuse/1024)
	log.Printf("Heap Idle: %d", m.HeapIdle/1024)
	log.Printf("Heap Released: %d", m.HeapReleased/1024)
	log.Printf("Heap Objects: %d", m.HeapObjects/1024)
	log.Printf("NumGC: %d", m.NumGC)
	<-quit
}

func WriteHeapProfile() {
	f, err := os.Create("memprofile.latest.mprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
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
		//	os.Exit(1)
	}
	end := time.Now()
	log.Printf("done with %s. took %v", output, end.Sub(start))
	c <- output
}
