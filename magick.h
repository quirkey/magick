#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <magick/MagickCore.h>

extern void SetImageInfoFilename(ImageInfo *image_info, char *filename);
extern MagickBooleanType CheckException(ExceptionInfo *exception);
extern MagickBooleanType GetBlobSupport(ImageInfo *image_info);
extern Image *ReadImageFromBlob(ImageInfo *image_info, void *blob, size_t length);
extern Image *AddShadowToImage(Image *image, char *colorname, const double opacity,
  const double sigma,const ssize_t x_offset,const ssize_t y_offset,
  ExceptionInfo *exception);
extern Image *FillBackgroundColor(Image *image, char *colorname, ExceptionInfo *exception);
extern Image *SeparateAlphaChannel(Image *image, ExceptionInfo *exception);
extern Image *Negate(Image *image, ExceptionInfo *exception);
