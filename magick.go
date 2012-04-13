package magick

/*
#cgo pkg-config: MagickCore
#include <stdio.h>
#include <stdlib.h>
#include <magick/MagickCore.h>

void SetImageInfoFilename(ImageInfo *image_info, char *filename) 
{
  (void) CopyMagickString(image_info->filename,filename,MaxTextExtent);
}

MagickBooleanType GetBlobSupport(ImageInfo *image_info) 
{
  ExceptionInfo *exception;
  const MagickInfo *magick_info;

  exception = AcquireExceptionInfo();
  magick_info = GetMagickInfo(image_info->magick,exception);
  CatchException(exception);
  return GetMagickBlobSupport(magick_info);
}

Image *ReadImageFromBlob(ImageInfo *image_info, void *blob, size_t length) {
  Image *image;
  ExceptionInfo *exception;
  exception = AcquireExceptionInfo();
  *image_info->filename='\0';
  *image_info->magick='\0';
  image_info->blob = blob;
  image_info->length = length;
  image = ReadImage(image_info, exception);
  CatchException(exception);
  return image;
}
*/
import "C"
import (
	"log"
	"os"
	"unsafe"
)

func init() {
	wd, err := os.Getwd()
	log.Printf("Working dir %s", wd)
	if err != nil {
		log.Fatal(err)
	}
	c_wd := C.CString(wd)
	defer C.free(unsafe.Pointer(c_wd))
	C.MagickCoreGenesis(c_wd, C.MagickTrue)
}

type MagickImage struct {
	image     (*C.Image)
	exception (*C.ExceptionInfo)
	imageInfo (*C.ImageInfo)
}

func NewFromFile(filename string) (im *MagickImage, ok bool) {
	exception := C.AcquireExceptionInfo()
	imageInfo := C.AcquireImageInfo()
	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))
	C.SetImageInfoFilename(imageInfo, c_filename)
	image := C.ReadImage(imageInfo, exception)
	C.CatchException(exception)
	im = &MagickImage{image, exception, imageInfo}
	return
}

// func NewFromBlob(blob []byte) (im *MagickImage, ok bool) {

// }

// func (im *MagickImage) Transform(crop_geometry, image_geometry string) (ok bool) {

// }

// func (im *MagickImage) ToBlob() (blob []byte, ok bool) {

// }

func (im *MagickImage) ToFile(filename string) (ok bool) {
	c_outpath := C.CString(filename)
	defer C.free(unsafe.Pointer(c_outpath))
	success := C.ImageToFile(im.image, c_outpath, im.exception)
	C.CatchException(im.exception)
	if success == C.MagickTrue {
		ok = true
	}
	return
}
