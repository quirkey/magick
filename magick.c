#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <magick/MagickCore.h>

void SetImageInfoFilename(ImageInfo *image_info, char *filename)
{
  (void) CopyMagickString(image_info->filename,filename,MaxTextExtent);
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

MagickBooleanType GetBlobSupport(ImageInfo *image_info)
{
  ExceptionInfo *exception;
  const MagickInfo *magick_info;
  MagickBooleanType supported;
  MagickBooleanType err;

  exception = AcquireExceptionInfo();
  magick_info = GetMagickInfo(image_info->magick,exception);
  if (magick_info == (const MagickInfo *) NULL) {
      return MagickFalse;
  }
  err = CheckException(exception);
  DestroyExceptionInfo(exception);
  if (err == MagickTrue) {
    return MagickFalse;
  }
  supported = GetMagickBlobSupport(magick_info);
  return supported;
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
