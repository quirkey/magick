// Package magick implements image manipulation routines based on the 
// ImageMagick MagickCore C library
//
// magick's goal is to provide quick processing of images in ways most
// commonly used by photo and other image based applications. It is not
// the intention to implement the very large number of methods that 
// MagickCore has to offer, rather just the most common needs for 
// basic applications
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
  DestroyExceptionInfo(exception);
  return GetMagickBlobSupport(magick_info);
}

Image *ReadImageFromBlob(ImageInfo *image_info, void *blob, size_t length)
{
  Image *image;
  ExceptionInfo *exception;
  exception = AcquireExceptionInfo();
  image_info->blob = blob;
  image_info->length = length;
  image = ReadImage(image_info, exception);
  CatchException(exception);
  DestroyExceptionInfo(exception);
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

Image *AddShadowToImage(Image *image, char *colorname, const double opacity,
  const double sigma,const ssize_t x_offset,const ssize_t y_offset,
  ExceptionInfo *exception)
{

  Image *shadow_image;
  if (QueryColorDatabase(colorname, &image->background_color, exception) == MagickFalse) {
    return MagickFalse;
  }
  shadow_image = ShadowImage(image, opacity, sigma, x_offset, y_offset, exception);
  AppendImageToList(&shadow_image, image);
  if (QueryColorDatabase("none", &shadow_image->background_color, exception) == MagickFalse) {
    return MagickFalse;
  }
  image = MergeImageLayers(shadow_image, MergeLayer, exception);
  DestroyImage(shadow_image);
  return image;
}

Image *FillBackgroundColor(Image *image, char *colorname, ExceptionInfo *exception)
{
    Image *new_image;
    new_image = CloneImage(image, 0, 0, MagickTrue, exception);
    if (QueryColorDatabase(colorname, &image->background_color, exception) == MagickFalse) {
      return MagickFalse;
    }
    if (SetImageBackgroundColor(image) == MagickFalse) {
      return MagickFalse;
    }
    AppendImageToList(&image, new_image);
    image = MergeImageLayers(image, MergeLayer, exception);
    DestroyImage(new_image);
    return image;
}

*/
import "C"
import (
	"log"
	"os"
	"strings"
	"unsafe"
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	c_wd := C.CString(wd)
	C.MagickCoreGenesis(c_wd, C.MagickFalse)
	defer C.free(unsafe.Pointer(c_wd))
}

type MagickImage struct {
	Image (*C.Image)
}

type MagickGeometry struct {
	Width, Height, Xoffset, Yoffset int
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
	defer C.DestroyExceptionInfo(exception)
	info := C.AcquireImageInfo()
	defer C.DestroyImageInfo(info)
	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))
	C.SetImageInfoFilename(info, c_filename)
	image := C.ReadImage(info, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	return &MagickImage{image}, nil
}

func NewFromBlob(blob []byte, extension string) (im *MagickImage, err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	info := C.AcquireImageInfo()
	defer C.DestroyImageInfo(info)
	c_filename := C.CString("image." + extension)
	defer C.free(unsafe.Pointer(c_filename))
	C.SetImageInfoFilename(info, c_filename)
	var success (C.MagickBooleanType)
	success = C.SetImageInfo(info, 1, exception)
	if success != C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	success = C.GetBlobSupport(info)
	if success != C.MagickTrue {
		return nil, &MagickError{"fatal", "", "image format " + extension + " does not support blobs"}
	}
	blob_copy := make([]byte, len(blob))
	copy(blob_copy, blob)
	length := (C.size_t)(len(blob_copy))
	blob_start := unsafe.Pointer(&blob_copy[0])
	image := C.ReadImageFromBlob(info, blob_start, length)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	return &MagickImage{image}, nil
}

func (im *MagickImage) Destroy() (err error) {
	C.DestroyImage(im.Image)
	im.Image = nil
	return
}

func (im *MagickImage) Width() int {
	return (int)(im.Image.columns)
}

func (im *MagickImage) Height() int {
	return (int)(im.Image.rows)
}

func (im *MagickImage) ParseGeometryToRectangleInfo(geometry string) (info (C.RectangleInfo), err error) {
	c_geometry := C.CString(geometry)
	defer C.free(unsafe.Pointer(c_geometry))
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	C.ParseRegionGeometry(im.Image, c_geometry, &info, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		err = ErrorFromExceptionInfo(exception)
	}
	return
}

func (im *MagickImage) ParseGeometry(geometry string) (info *MagickGeometry, err error) {
	rectangle, err := im.ParseGeometryToRectangleInfo(geometry)
	if err != nil {
		return nil, err
	}
	return &MagickGeometry{int(rectangle.width), int(rectangle.height), int(rectangle.x), int(rectangle.y)}, nil
}

func (im *MagickImage) Resize(geometry string) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	rect, err := im.ParseGeometryToRectangleInfo(geometry)
	if err != nil {
		return err
	}
	new_image := C.ThumbnailImage(im.Image, rect.width, rect.height, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.Destroy()
	im.Image = new_image
	return nil
}

func (im *MagickImage) Crop(geometry string) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	rect, err := im.ParseGeometryToRectangleInfo(geometry)
	if err != nil {
		return err
	}
	new_image := C.CropImage(im.Image, &rect, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.Destroy()
	im.Image = new_image
	return nil
}

func (im *MagickImage) Shadow(color string, opacity, sigma float32, xoffset, yoffset int) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	c_opacity := (C.double)(opacity)
	c_sigma := (C.double)(sigma)
	c_x := (C.ssize_t)(xoffset)
	c_y := (C.ssize_t)(yoffset)
	c_color := C.CString(color)
	defer C.free(unsafe.Pointer(c_color))
	new_image := C.AddShadowToImage(im.Image, c_color, c_opacity, c_sigma, c_x, c_y, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.Destroy()
	im.Image = new_image
	return nil
}

func (im *MagickImage) FillBackgroundColor(color string) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	c_color := C.CString(color)
	defer C.free(unsafe.Pointer(c_color))
	new_image := C.FillBackgroundColor(im.Image, c_color, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.Destroy()
	im.Image = new_image
	return nil
}

func (im *MagickImage) ToBlob(filetype string) (blob []byte, err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	image_info := C.AcquireImageInfo()
	defer C.DestroyImageInfo(image_info)
	c_outpath := C.CString("image." + filetype)
	defer C.free(unsafe.Pointer(c_outpath))
	C.SetImageInfoFilename(image_info, c_outpath)
	var outlength (C.size_t)
	outblob := C.ImageToBlob(image_info, im.Image, &outlength, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	char_pointer := unsafe.Pointer(outblob)
	return C.GoBytes(char_pointer, (C.int)(outlength)), nil
}

func (im *MagickImage) ToFile(filename string) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	c_outpath := C.CString(filename)
	defer C.free(unsafe.Pointer(c_outpath))
	write_info := C.AcquireImageInfo()
	defer C.DestroyImageInfo(write_info)
	C.SetImageInfoFilename(write_info, c_outpath)
	success := C.WriteImages(write_info, im.Image, c_outpath, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	if success != C.MagickTrue {
		return &MagickError{"fatal", "", "could not write to " + filename + " for unknown reason"}
	}
	return nil
}

func (im *MagickImage) Type() (t string) {
	return strings.Trim(string(C.GoBytes(unsafe.Pointer(&im.Image.magick), 4096)), "\x00")
}
