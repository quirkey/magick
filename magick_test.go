package magick

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestImageFromFile(t *testing.T) {
	filename := "test/heart_original.png"
	_, ok := NewFromFile(filename)
	assert.T(t, ok)

	bad_filename := "test/heart_whatwhat.png"
	_, ok = NewFromFile(bad_filename)
	assert.T(t, !ok)
}
