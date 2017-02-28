package magick

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

func setupImage(t *testing.T) (image *MagickImage) {
	filename := "test/heart_original.png"
	image, error := NewFromFile(filename)
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)
	return
}

func TestImageFromFile(t *testing.T) {
	filename := "test/heart_original.png"
	image, error := NewFromFile(filename)
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)

	bad_filename := "test/heart_whatwhat.png"
	image, error = NewFromFile(bad_filename)
	assert.T(t, error != nil)
}

func TestImageFromBlob(t *testing.T) {
	filename := "test/heart_original.png"
	source, _ := ioutil.ReadFile(filename)
	image, error := NewFromBlob(source, "png")
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)
	image, error = NewFromBlob([]byte{}, "png")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
}

func TestBadDataFromBlob(t *testing.T) {
	filename := "test/heart_original.png"
	source, _ := ioutil.ReadFile(filename)
	image, error := NewFromBlob(source, "")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte(""), "png")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte("blah"), "jpg")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte("   "), "image/jpg")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte(""), " ")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil
	image, error = NewFromBlob([]byte(""), ":")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil
}

func TestPDFFromBlob(t *testing.T) {
	filename := "test/heart_original.pdf"
	source, _ := ioutil.ReadFile(filename)
	image, error := NewFromBlob(source, "pdf")
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)
}

func TestParseGeometry(t *testing.T) {
	image := setupImage(t)
	geometry, err := image.ParseGeometry("100x100>")
	assert.T(t, err == nil)
	assert.T(t, geometry != nil)
	assert.Equal(t, 100, geometry.Width)
}

func TestResizeRatio(t *testing.T) {
	image := setupImage(t)
	ratio := image.ResizeRatio(300, 300)
	assert.T(t, ratio > 0.27)
	assert.T(t, ratio < 0.28)
}

func TestStrip(t *testing.T) {
	image := setupImage(t)
	err := image.Strip()
	assert.T(t, err == nil)
}

func TestProgressive(t *testing.T) {
	image := setupImage(t)
	image.Progressive()
	_, err := image.ToBlob("jpg")
	assert.T(t, err == nil)
}

func TestDestroy(t *testing.T) {
	image := setupImage(t)
	assert.T(t, image.Destroy() == nil)
	assert.T(t, image.Image == nil)
	assert.T(t, image.ImageInfo == nil)
}

func TestResize(t *testing.T) {
	image := setupImage(t)
	err := image.Resize("100x100!")
	assert.T(t, err == nil)
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 100, image.Height())

	image = setupImage(t)
	err = image.Resize("blurgh")
	assert.T(t, err != nil)
}

func TestPDFResize(t *testing.T) {
	filename := "test/heart_original.pdf"
	source, _ := ioutil.ReadFile(filename)
	image, err := NewFromBlob(source, "pdf")
	assert.T(t, err == nil)
	err = image.Resize("100x100!")
	assert.T(t, err == nil)
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 100, image.Height())
	if err == nil {
		filename = "test/test_from_pdf.jpg"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
}

func TestCrop(t *testing.T) {
	image := setupImage(t)
	err := image.Crop("100x100!+10+10")
	assert.T(t, err == nil)
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 100, image.Height())

	image = setupImage(t)
	err = image.Crop("blurgh")
	assert.T(t, err != nil)
}

// The use case for +repage is that we need to reset an image's canvas so only
// the cropped portion is present. Cropping and not using +repage leaves the
// canvas the same size, but with only the cropped portion's pixels remaining.
// I believe this is only the case for certain image formats such as GIF and
// PNG.
//
// Test to demonstrate the problem and show how PlusRepage() solves it.
func TestPlusRepage(t *testing.T) {
	filename := "test/heart_original.png"
	image, err := NewFromFile(filename)
	if err != nil {
		t.Errorf("NewFromFile(%s) error = %s", filename, err)
		return
	}

	origWidth := image.Width()
	origHeight := image.Height()
	origVirtualCanvasWidth := image.VirtualCanvasWidth()
	origVirtualCanvasHeight := image.VirtualCanvasHeight()

	wantedWidth := 100
	wantedHeight := 100
	geom := fmt.Sprintf("%dx%d!+10+10", wantedWidth, wantedHeight)
	if err := image.Crop(geom); err != nil {
		t.Errorf("image.Crop(%s) error = %s", geom, err)
		_ = image.Destroy()
		return
	}

	if image.Width() != wantedWidth || image.Height() != wantedHeight {
		t.Errorf("image width x height after cropping is %dx%d, wanted %dx%d",
			image.Width(), image.Height(), origWidth, origHeight)
		_ = image.Destroy()
		return
	}

	// The virtual canvas width x height should currently be as it was. The
	// cropped region will be surrounded by empty space if we save it now.
	if image.VirtualCanvasWidth() != origVirtualCanvasWidth ||
		image.VirtualCanvasHeight() != origVirtualCanvasHeight {
		t.Errorf("virtual canvas width x height after cropping is %dx%d, wanted %dx%d",
			image.VirtualCanvasWidth(), image.VirtualCanvasHeight(),
			origVirtualCanvasWidth, origVirtualCanvasHeight)
		_ = image.Destroy()
		return
	}

	// Reset the virtual canvas using +repage. This removes virtual canvas
	// information. The effect of removing it means the image, once saved, will
	// have the dimensions of the image region we cropped.
	image.PlusRepage()

	// Check both sets of dimensions again. The image size should still be the
	// same. The virtual canvas size should be 0x0 since we removed information
	// about it.

	if image.Width() != wantedWidth || image.Height() != wantedHeight {
		t.Errorf("image width x height after cropping is %dx%d, wanted %dx%d",
			image.Width(), image.Height(), origWidth, origHeight)
		_ = image.Destroy()
		return
	}

	if image.VirtualCanvasWidth() != 0 || image.VirtualCanvasHeight() != 0 {
		t.Errorf("virtual canvas width x height after +repage is %dx%d, wanted %dx%d",
			image.VirtualCanvasWidth(), image.VirtualCanvasHeight(), 0, 0)
		_ = image.Destroy()
		return
	}

	if err := image.Destroy(); err != nil {
		t.Errorf("image.Destroy() error = %s", err)
		return
	}
}

func TestShadow(t *testing.T) {
	image := setupImage(t)
	err := image.Shadow("#000", 75, 2, 0, 0)
	assert.T(t, err == nil)
}

func TestFillBackgroundColor(t *testing.T) {
	image := setupImage(t)
	err := image.FillBackgroundColor("#CCC")
	assert.T(t, err == nil)
}

func TestSeparateAlphaChannel(t *testing.T) {
	image := setupImage(t)
	err := image.SeparateAlphaChannel()
	assert.T(t, err == nil)
	assert.T(t, image != nil)
}

func TestNegateImage(t *testing.T) {
	image := setupImage(t)
	err := image.Negate()
	assert.T(t, err == nil)
	assert.T(t, image != nil)
}

func TestToBlob(t *testing.T) {
	image := setupImage(t)
	bytes, err := image.ToBlob("png")
	assert.T(t, err == nil)
	assert.T(t, bytes != nil)
	assert.T(t, len(bytes) > 0)
	assert.Equal(t, 437198, len(bytes))
}

func TestToFile(t *testing.T) {
	image := setupImage(t)
	filename := "test/test_out.png"
	os.Remove(filename)
	err := image.ToFile(filename)
	assert.T(t, err == nil)
	file, err := os.Open(filename)
	assert.T(t, err == nil)
	defer file.Close()
	stat, err := file.Stat()
	assert.T(t, stat != nil)
	assert.Equal(t, (int64)(437198), stat.Size())
}

func TestType(t *testing.T) {
	image := setupImage(t)
	assert.Equal(t, "PNG", image.Type())
}

func TestWidth(t *testing.T) {
	image := setupImage(t)
	assert.Equal(t, 600, image.Width())
}

func TestHeight(t *testing.T) {
	image := setupImage(t)
	assert.Equal(t, 552, image.Height())
}

func TestSetProperty(t *testing.T) {
	image := setupImage(t)
	err := image.SetProperty("jpeg:sampling-factor", "4:4:4")
	assert.T(t, err == nil)
	factor := image.GetProperty("jpeg:sampling-factor")
	assert.Equal(t, "4:4:4", factor)
}

func TestFullStack(t *testing.T) {
	var err error
	var filename string
	var image *MagickImage
	// resize
	image = setupImage(t)
	err = image.Resize("100x100")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_resize.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// crop
	image = setupImage(t)
	err = image.Crop("100x100+10+10")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_crop.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// shadow
	image = setupImage(t)
	err = image.Shadow("#000", 90, 10, 0, 0)
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_shadow.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// fill
	image = setupImage(t)
	err = image.FillBackgroundColor("red")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_fill.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// combination
	image = setupImage(t)
	err = image.Resize("100x100>")
	assert.T(t, err == nil)
	err = image.Shadow("#000", 90, 10, 0, 0)
	assert.T(t, err == nil)
	err = image.FillBackgroundColor("#CCC")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_combo.jpg"
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// alpha mask
	image = setupImage(t)
	err = image.SeparateAlphaChannel()
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_alpha.jpg"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// alpha mask + negate
	image = setupImage(t)
	err = image.SeparateAlphaChannel()
	assert.T(t, err == nil)

	err = image.Negate()
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_alpha_negative.jpg"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
}

func TestAutoOrientNeedsToRotate(t *testing.T) {
	// This file is set so the heart appears to need to be rotated 90 degrees
	// clockwise. It has EXIF orientation set to indicate that.
	//
	// I created it from test/heart_original.png this way:
	// convert -rotate 270 heart_original.png heart_rotated.jpg
	// exiftool -Orientation=6 -n heart_rotated.jpg
	filename := "test/heart_rotated.jpg"

	image, err := NewFromFile(filename)
	if err != nil {
		t.Errorf("NewFromFile(%s) = %s", filename, err)
		return
	}

	err = image.AutoOrient()
	if err != nil {
		_ = image.Destroy()
		t.Errorf("AutoOrient() = %s", err)
		return
	}

	err = image.Destroy()
	if err != nil {
		t.Errorf("Destroy() = %s", err)
		return
	}
}

func TestAutoOrientNoNeedToRotate(t *testing.T) {
	// This file does not need to change.
	filename := "test/heart_original.png"

	image, err := NewFromFile(filename)
	if err != nil {
		t.Errorf("NewFromFile(%s) = %s", filename, err)
		return
	}

	err = image.AutoOrient()
	if err != nil {
		_ = image.Destroy()
		t.Errorf("AutoOrient() = %s", err)
		return
	}

	err = image.Destroy()
	if err != nil {
		t.Errorf("Destroy() = %s", err)
		return
	}
}
