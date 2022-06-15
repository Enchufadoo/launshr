package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeFileReader struct {
	exists bool
}

func (f FakeFileReader) ReadFile(filename string) ([]byte, error) {
	if !f.exists {
		return []byte{}, errors.New("error")
	} else {
		return []byte("Hello world"), nil
	}

}

func TestOpenConfig(t *testing.T) {
	t.Run("Open succesfully a file", func(t *testing.T) {
		res, err := openConfig(FakeFileReader{exists: true})
		isEqual := string(res) == string([]byte("Hello world"))
		assert.True(t, isEqual, "True is true!")
		assert.Nil(t, err)
		if err != nil {
			t.Errorf("Shouldn't be an error")
		}
	})

	t.Run("Opening a file that doesn't exist", func(t *testing.T) {
		_, err := openConfig(FakeFileReader{exists: false})
		assert.NotNil(t, err)
	})

}
