package magick

/*
#cgo pkg-config: MagickCore
#include <stdio.h>
#include <stdlib.h>
#include <magick/MagickCore.h>
#include "magick_wrappers.c"
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"unsafe"
)

type MagickImage struct {
  image (*C.Image)
}

func NewFromFile(filename string) (im *MagickImage, ok bool) {

}

func NewFromBlob(blob []byte) (im *MagickImage, ok bool) {

}

func (*MagickImage) Transform(crop_geometry, image_geometry string) (ok bool) {

}

func (*MagickImage) ToBlob() (blob []byte, ok bool) {

}

func (*MagickImage) ToFile(filename string) (ok bool) {

}
// func main() {
// 
// 	wd, err := os.Getwd()
// 	// log.Printf("Working dir %s", wd)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	c_wd := C.CString(wd)
// 	defer C.free(unsafe.Pointer(c_wd))
// 
// 	C.MagickCoreGenesis(c_wd, C.MagickTrue)
// 	defer C.MagickCoreTerminus()
// 	// "%s" \( +clone -background black -shadow 75x5+0+0 \) +swap -background "#F8F9F3" -layers merge +repage -format jpeg "%s"
// 
// 	filename := os.Args[1]
// 	source, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		log.Printf("Error reading from file %s: %d", filename, err)
// 		os.Exit(1)
// 	}
// 	start := time.Now()
// 	image_info := C.AcquireImageInfo()
// 	c_filename := C.CString(filename)
// 	defer C.free(unsafe.Pointer(c_filename))
// 	exception := C.AcquireExceptionInfo()
// 	C.SetImageInfoFilename(image_info, c_filename)
// 	var success (C.MagickBooleanType)
// 	success = C.SetImageInfo(image_info, 1, exception)
// 	if success != C.MagickTrue {
// 		C.CatchException(exception)
// 		os.Exit(1)
// 	}
// 	success = C.GetBlobSupport(image_info)
// 	if success != C.MagickTrue {
// 		log.Print("Does not support blobs")
// 		os.Exit(1)
// 	}
// 	length := (C.size_t)(len(source))
// 	// log.Printf("Reading blob of size %d", length)
// 	image := C.ReadImageFromBlob(image_info, unsafe.Pointer(&source[0]), length)
// 	// log.Printf("Read blob of size %d", image_info.length)
// 	C.CatchException(exception)
// 	if image == nil {
// 		os.Exit(1)
// 	}
// 	crop_geometry := C.CString("")
// 	image_geometry := C.CString("x100")
// 	defer C.free(unsafe.Pointer(crop_geometry))
// 	defer C.free(unsafe.Pointer(image_geometry))
// 	success = C.TransformImage(&image, crop_geometry, image_geometry)
// 	if success != C.MagickTrue {
// 		log.Printf("Transform Fail %d", success)
// 		os.Exit(1)
// 	} else {
// 		// log.Printf("Transform Success %d", success)
// 	}
// 	// c_path := C.CString(os.Args[1])
// 	// var success (C.MagickBooleanType)
// 	//C.CopyMagickString(image_info.filename, c_path, C.MaxTextExtent)
// 	//images := C.ReadImage(image_info,exception);
// 	// image := C.ReadImageFromFile(c_path)
// 	new_image_info := C.AcquireImageInfo()
// 	// success = C.WriteImage(new_image_info, image)
// 	// if success != C.MagickTrue {
// 	//   C.CatchException(exception)
// 	//   os.Exit(1)
// 	// }
// 	var outlength (C.size_t)
// 	outblob := C.ImageToBlob(new_image_info, image, &outlength, exception)
// 	C.CatchException(exception)
// 	log.Printf("Write Success %d", outlength)
// 	char_pointer := (*C.char)(unsafe.Pointer(outblob))
// 	fmt.Print(C.GoStringN(char_pointer, (C.int)(outlength)))
// 	end := time.Now()
// 	log.Printf("took %v\n", end.Sub(start))
// }
