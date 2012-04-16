package magick

/*
#cgo pkg-config: MagickCore
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
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

MagickBooleanType CheckException(ExceptionInfo *exception)
{
  register const ExceptionInfo
    *p;
  int haserr = 0;

  assert(exception != (ExceptionInfo *) NULL);
  assert(exception->signature == MagickSignature);
  if (exception->exceptions  == (void *) NULL)
    return MagickFalse;
  LockSemaphoreInfo(exception->semaphore);
  ResetLinkedListIterator((LinkedListInfo *) exception->exceptions);
  p=(const ExceptionInfo *) GetNextValueInLinkedList((LinkedListInfo *)
    exception->exceptions);
  while (p != (const ExceptionInfo *) NULL)
  {
    if ((p->severity >= WarningException) && (p->severity < ErrorException))
      haserr = 1;
    if ((p->severity >= ErrorException) && (p->severity < FatalErrorException))
      haserr = 1;
    if (p->severity >= FatalErrorException)
      haserr = 1;
    p=(const ExceptionInfo *) GetNextValueInLinkedList((LinkedListInfo *)
      exception->exceptions);
  }
  UnlockSemaphoreInfo(exception->semaphore);
  return haserr == 0 ? MagickFalse : MagickTrue;
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
	Image (*C.Image)

	exception (*C.ExceptionInfo)
	imageInfo (*C.ImageInfo)
}

type MagickError struct {
	Severity    string
	Reason      string
	Description string
}

func (err *MagickError) Error() string {
	return "MagickError " + err.Severity + ": " + err.Reason + "- " + err.Description
}

func ErrorFromExceptionInfo(exception *C.ExceptionInfo) (err error) {
	return &MagickError{string(exception.severity), C.GoString(exception.reason), C.GoString(exception.description)}
}

func NewFromFile(filename string) (im *MagickImage, err error) {
	exception := C.AcquireExceptionInfo()
	imageInfo := C.AcquireImageInfo()
	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))
	C.SetImageInfoFilename(imageInfo, c_filename)
	image := C.ReadImage(imageInfo, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	return &MagickImage{image, exception, imageInfo}, nil
}

func NewFromBlob(blob []byte, extension string) (im *MagickImage, err error) {
	imageInfo := C.AcquireImageInfo()
	c_filename := C.CString("image." + extension)
	defer C.free(unsafe.Pointer(c_filename))
	exception := C.AcquireExceptionInfo()
	C.SetImageInfoFilename(imageInfo, c_filename)
	var success (C.MagickBooleanType)
	success = C.SetImageInfo(imageInfo, 1, exception)
	if success != C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	success = C.GetBlobSupport(imageInfo)
	if success != C.MagickTrue {
		return nil, &MagickError{"fatal", "", "image format " + extension + " does not support blobs"}
	}
	length := (C.size_t)(len(blob))
	image := C.ReadImageFromBlob(imageInfo, unsafe.Pointer(&blob[0]), length)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	return &MagickImage{image, exception, imageInfo}, nil
}

// func (im *MagickImage) Transform(crop_geometry, image_geometry string) (ok bool) {

// }

func (im *MagickImage) ToBlob() (blob []byte, err error) {
	new_image_info := C.AcquireImageInfo()
	var outlength (C.size_t)
	outblob := C.ImageToBlob(new_image_info, im.Image, &outlength, im.exception)
	if failed := C.CheckException(im.exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(im.exception)
	}
	char_pointer := unsafe.Pointer(outblob)
	return C.GoBytes(char_pointer, (C.int)(outlength)), nil
}

func (im *MagickImage) ToFile(filename string) (ok bool, err error) {
	c_outpath := C.CString(filename)
	defer C.free(unsafe.Pointer(c_outpath))
	success := C.ImageToFile(im.Image, c_outpath, im.exception)
	C.CatchException(im.exception)
	if failed := C.CheckException(im.exception); failed == C.MagickTrue {
		return false, ErrorFromExceptionInfo(im.exception)
	}
	if success == C.MagickTrue {
		ok = true
	}
	return
}
