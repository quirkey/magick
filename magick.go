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

MagickBooleanType SetBackgroundColor(Image *image, char *colorname, ExceptionInfo *exception) {
    return QueryColorDatabase(colorname, &image->background_color, exception);
}

Image *FillBackgroundColor(Image *image, char *colorname, ExceptionInfo *exception) {
    Image *new_image;
    new_image = CloneImage(image, 0, 0, MagickTrue, exception);
    if (SetBackgroundColor(new_image, colorname, exception) == MagickFalse) {
      return MagickFalse;
    }
    if (SetImageBackgroundColor(new_image) == MagickFalse) {
      return MagickFalse;
    }
    AppendImageToList(&new_image, image);    
    return MergeImageLayers(new_image, FlattenLayer, exception);
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

func (im *MagickImage) Transform(crop_geometry, image_geometry string) (ok bool) {
	c_crop_geometry := C.CString(crop_geometry)
	c_image_geometry := C.CString(image_geometry)
	defer C.free(unsafe.Pointer(c_crop_geometry))
	defer C.free(unsafe.Pointer(c_image_geometry))
	success := C.TransformImage(&im.Image, c_crop_geometry, c_image_geometry)
	if success == C.MagickTrue {
		ok = true
	}
	C.CheckException(im.exception)
	return
}

func (im *MagickImage) Thumbnail(width, height int) (resized *MagickImage, err error) {
	c_cols := (C.size_t)(width)
	c_rows := (C.size_t)(height)
	new_image := C.ThumbnailImage(im.Image, c_cols, c_rows, im.exception)
	if failed := C.CheckException(im.exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(im.exception)
	}
	return &MagickImage{new_image, im.exception, C.AcquireImageInfo()}, nil
}

func (im *MagickImage) Shadow(color string, opacity, sigma float32, xoffset, yoffset int) (shadowed *MagickImage, err error) {
        c_opacity := (C.double)(opacity)
        c_sigma := (C.double)(sigma)
	c_x := (C.ssize_t)(xoffset)
	c_y := (C.ssize_t)(yoffset)
        c_color := C.CString(color)
	defer C.free(unsafe.Pointer(c_color))
        C.SetBackgroundColor(im.Image, c_color, im.exception)
	new_image := C.ShadowImage(im.Image, c_opacity, c_sigma, c_x, c_y, im.exception)
	if failed := C.CheckException(im.exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(im.exception)
	}
	return &MagickImage{new_image, im.exception, C.AcquireImageInfo()}, nil
}

func (im *MagickImage) FillBackgroundColor(color string) (flattened *MagickImage, err error) {
        c_color := C.CString(color)
	defer C.free(unsafe.Pointer(c_color))
        new_image := C.FillBackgroundColor(im.Image, c_color, im.exception)
	if failed := C.CheckException(im.exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(im.exception)
	}
	return &MagickImage{new_image, im.exception, C.AcquireImageInfo()}, nil
}

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
	write_info := C.AcquireImageInfo()
	C.SetImageInfoFilename(im.imageInfo, c_outpath)
	success := C.WriteImages(write_info, im.Image, c_outpath, im.exception)
	if failed := C.CheckException(im.exception); failed == C.MagickTrue {
		return false, ErrorFromExceptionInfo(im.exception)
	}
	if success == C.MagickTrue {
		ok = true
	}
	return
}
