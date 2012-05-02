package magick

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"testing"
)

func setupImage(t *testing.T) (image *MagickImage) {
	filename := "test/heart_original.png"
	image, error := NewFromFile(filename)
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	return
}

func TestImageFromFile(t *testing.T) {
	filename := "test/heart_original.png"
	image, error := NewFromFile(filename)
	assert.T(t, error == nil)
	assert.T(t, image != nil)

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
}

func TestParseGeometry(t *testing.T) {
	image := setupImage(t)
	geometry, err := image.ParseGeometry("100x100>")
	assert.T(t, err == nil)
	assert.T(t, geometry != nil)
	assert.Equal(t, 100, geometry.Width)
}

func TestDestroy(t *testing.T) {
	image := setupImage(t)
	assert.T(t, image.Destroy() == nil)
	assert.T(t, image.Image == nil)
}

func TestResize(t *testing.T) {
	image := setupImage(t)
	err := image.Resize("100x100!")
	assert.T(t, err == nil)
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 100, image.Height())
}
