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
