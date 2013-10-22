// Package magick implements image manipulation routines based on the
// ImageMagick MagickCore C library. It is an opinionated high level
// wrapper around the proven ImageMagick lib.
//
// magick's goal is to provide quick processing of images in ways most
// commonly used by photo and other image based applications. It is not
// the intention to implement the very large number of methods that
// MagickCore has to offer, rather just the most common needs for
// basic applications. It requires ImageMagick-devel libraries to
// be available in order to compile.
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

Image *SeparateAlphaChannel(Image *image, ExceptionInfo *exception){
  Image *new_image;
  new_image = CloneImage(image, 0, 0, MagickTrue, exception);
  if (SeparateImageChannel(new_image, 0x0008) == MagickFalse){
    return MagickFalse;
  }
  return new_image;
}

Image *Negate(Image *image, ExceptionInfo *exception){
  Image *new_image;
  new_image = CloneImage(image, 0, 0, MagickTrue, exception);
  if (NegateImage(new_image, MagickTrue) == MagickFalse){
    return MagickFalse;
  }
  return new_image;
}

*/
import "C"
import (
	"io/ioutil"
	"math"
	"os"
	"strings"
	"unsafe"
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	c_wd := C.CString(wd)
	C.MagickCoreGenesis(c_wd, C.MagickFalse)
	defer C.free(unsafe.Pointer(c_wd))
}

// A wrapper around an IM Image
type MagickImage struct {
	Image     (*C.Image)
	ImageInfo (*C.ImageInfo)
}

// Geometry is usually defined as a string of WxH+X+Y
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

// NewFromFile loads a file at filename into a MagickImage.
// Exceptions are returned as MagickErrors.
func NewFromFile(filename string) (im *MagickImage, err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	info := C.AcquireImageInfo()
	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))
	C.SetImageInfoFilename(info, c_filename)
	image := C.ReadImage(info, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		C.DestroyImageInfo(info)
		return nil, ErrorFromExceptionInfo(exception)
	}
	return &MagickImage{Image: image, ImageInfo: info}, nil
}

// NewFromBlob takes a byte slice of image data and an extension that defines the
// image type (e.g. "png", "jpg", etc). It loads the image data and returns a MagickImage.
// The extension is required so that Magick knows what processor to use.
func NewFromBlob(blob []byte, extension string) (im *MagickImage, err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	info := C.AcquireImageInfo()
	defer C.DestroyImageInfo(info)
	c_filename := C.CString("image." + extension)
	defer C.free(unsafe.Pointer(c_filename))
	C.SetImageInfoFilename(info, c_filename)
	cloned_info := C.CloneImageInfo(info)
	var success (C.MagickBooleanType)
	success = C.SetImageInfo(info, 1, exception)
	if success != C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	success = C.GetBlobSupport(info)
	if success != C.MagickTrue {
		// No blob support, lets try reading from a file
		file, err := ioutil.TempFile("", "image."+extension)
		if _, err = file.Write(blob); err != nil {
			return nil, &MagickError{"fatal", "", "image format " + extension + " does not support blobs and could not write temp file"}
		}
		file.Close()
		return NewFromFile(file.Name())
	}
	blob_copy := make([]byte, len(blob))
	copy(blob_copy, blob)
	length := (C.size_t)(len(blob_copy))
	blob_start := unsafe.Pointer(&blob_copy[0])
	image := C.ReadImageFromBlob(info, blob_start, length)

	if image == nil {
		return nil, &MagickError{"fatal", "", "corrupt image, not a " + extension}
	}

	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}

	return &MagickImage{Image: image, ImageInfo: cloned_info}, nil
}

// Destroy frees the C memory for the image. Should be called after processing is done.
func (im *MagickImage) Destroy() (err error) {
	C.DestroyImage(im.Image)
	C.DestroyImageInfo(im.ImageInfo)
	im.Image = nil
	im.ImageInfo = nil
	return
}

// ReplaceImage Replaces the underlying image, freeing the old one
func (im *MagickImage) ReplaceImage(new_image *C.Image) {
	C.DestroyImage(im.Image)
	im.Image = new_image
}

// Width returns the Width of the loaded image in pixels as an int
func (im *MagickImage) Width() int {
	return (int)(im.Image.columns)
}

// Height returns the Height of the loaded image in pixels as an int
func (im *MagickImage) Height() int {
	return (int)(im.Image.rows)
}

// Type returns the underlying encoding or "magick" of the image as a string
func (im *MagickImage) Type() (t string) {
	return strings.Trim(string(C.GoBytes(unsafe.Pointer(&im.Image.magick), 4096)), "\x00")
}

// ResizeRatio() returns the ratio that the size you want to resize to
// defined by width/height is over the size of the underlying Image
func (im *MagickImage) ResizeRatio(width, height int) float64 {
	return math.Abs((float64)(width*height) / (float64)(im.Width()*im.Height()))
}

// GetProperty() retreives the given attribute or freeform property
// string on the underlying Image
func (im *MagickImage) GetProperty(prop string) (value string) {
	c_prop := C.CString(prop)
	defer C.free(unsafe.Pointer(c_prop))
	c_value := C.GetImageProperty(im.Image, c_prop)
	defer C.free(unsafe.Pointer(c_value))
	return C.GoString(c_value)
}

// SetProperty() saves the given string value either to specific known
// attribute or to a freeform property string on the underlying Image
func (im *MagickImage) SetProperty(prop, value string) (err error) {
	c_prop := C.CString(prop)
	defer C.free(unsafe.Pointer(c_prop))
	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_value))
	ok := C.SetImageProperty(im.Image, c_prop, c_value)
	if ok == C.MagickFalse {
		return &MagickError{"error", "", "could not set property"}
	}
	return
}

// ParseGeometryToRectangleInfo converts from a geometry string (WxH+X+Y) into a Magick
// RectangleInfo that contains the individual properties
func (im *MagickImage) ParseGeometryToRectangleInfo(geometry string) (info C.RectangleInfo, err error) {
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

// ParseGeometry uses ParseGeometryToRectangleInfo to convert from a geometry string into a MagickGeometry
func (im *MagickImage) ParseGeometry(geometry string) (info *MagickGeometry, err error) {
	rectangle, err := im.ParseGeometryToRectangleInfo(geometry)
	if err != nil {
		return nil, err
	}
	return &MagickGeometry{int(rectangle.width), int(rectangle.height), int(rectangle.x), int(rectangle.y)}, nil
}

// Progessive() is a shortcut for making the underlying image a
// Plane interlaced Progressive JPG
func (im *MagickImage) Progressive() {
	im.ImageInfo.interlace = C.PlaneInterlace
}

func (im *MagickImage) Quality(quality int) {
	im.Image.quality = (C.size_t)(quality)
}

// Resize resizes the image based on the geometry string passed and stores the resized image in place
// For more info about Geometry see http://www.imagemagick.org/script/command-line-processing.php#geometry
func (im *MagickImage) Resize(geometry string) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	rect, err := im.ParseGeometryToRectangleInfo(geometry)
	if err != nil {
		return err
	}
	ratio := im.ResizeRatio(int(rect.width), int(rect.height))
	var new_image *C.Image
	if ratio > 0.4 {
		new_image = C.AdaptiveResizeImage(im.Image, rect.width, rect.height, exception)
	} else {
		new_image = C.ThumbnailImage(im.Image, rect.width, rect.height, exception)
	}
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.ReplaceImage(new_image)
	return nil
}

// Crop crops the image based on the geometry string passed and stores the cropped image in place
// For more info about Geometry see http://www.imagemagick.org/script/command-line-processing.php#geometry
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
	im.ReplaceImage(new_image)
	return nil
}

// Shadow adds a dropshadow to the current (transparent) image and stores the shadowed image in place
// For more information about shadow options see: http://www.imagemagick.org/Usage/blur/#shadow
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
	im.ReplaceImage(new_image)
	return nil
}

// FillBackgroundColor fills transparent areas of an image with a solid color and stores the filled image in place.
// color can be any color format that image magick understands, see: http://www.imagemagick.org/ImageMagick-7.0.0/script/color.php#models
func (im *MagickImage) FillBackgroundColor(color string) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	c_color := C.CString(color)
	defer C.free(unsafe.Pointer(c_color))
	new_image := C.FillBackgroundColor(im.Image, c_color, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.ReplaceImage(new_image)
	return nil
}

// SeparateAlphaChannel replaces the Image with grayscale data from the image's Alpha Channel values
func (im *MagickImage) SeparateAlphaChannel() (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	new_image := C.SeparateAlphaChannel(im.Image, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.ReplaceImage(new_image)
	return nil
}

// Negate inverts the colors in the image
func (im *MagickImage) Negate() (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	new_image := C.Negate(im.Image, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	im.ReplaceImage(new_image)
	return nil
}

// Strip strips the image of its extra meta (exif) data
func (im *MagickImage) Strip() (err error) {
	ok := C.StripImage(im.Image)
	if ok == C.MagickFalse {
		return &MagickError{"error", "", "could not strip image"}
	}
	return
}

// ToBlob takes a (transformed) MagickImage and returns a byte slice in the format you specify with extension.
// Magick uses the extension to transform the image in to the proper encoding (e.g. "jpg", "png")
func (im *MagickImage) ToBlob(extension string) (blob []byte, err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	c_outpath := C.CString("image." + extension)
	defer C.free(unsafe.Pointer(c_outpath))
	C.SetImageInfoFilename(im.ImageInfo, c_outpath)
	var outlength (C.size_t)
	outblob := C.ImageToBlob(im.ImageInfo, im.Image, &outlength, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return nil, ErrorFromExceptionInfo(exception)
	}
	char_pointer := unsafe.Pointer(outblob)
	defer C.free(char_pointer)
	return C.GoBytes(char_pointer, (C.int)(outlength)), nil
}

// ToFile writes the (transformed) MagickImage to the regular file at filename. Magick determines
// the encoding of the output file by the extension given to the filename (e.g. "image.jpg", "image.png")
func (im *MagickImage) ToFile(filename string) (err error) {
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)
	c_outpath := C.CString(filename)
	defer C.free(unsafe.Pointer(c_outpath))
	C.SetImageInfoFilename(im.ImageInfo, c_outpath)
	success := C.WriteImages(im.ImageInfo, im.Image, c_outpath, exception)
	if failed := C.CheckException(exception); failed == C.MagickTrue {
		return ErrorFromExceptionInfo(exception)
	}
	if success != C.MagickTrue {
		return &MagickError{"fatal", "", "could not write to " + filename + " for unknown reason"}
	}
	return nil
}
