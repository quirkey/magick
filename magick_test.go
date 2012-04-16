package magick

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"testing"
)

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
