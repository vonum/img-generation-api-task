package image_test

import (
	"testing"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

func TestImageMaxSizeExceededError(t *testing.T) {
  err := image.ImageMaxSizeExceededError{8192 * 1024}
  errMsg := "Image size exceeded - 8192 kB."
  if err.Error() != errMsg {
    t.Errorf("Expected %s but got %s", errMsg, err.Error())
  }
}

func TestMimeTypeError(t *testing.T) {
  err := image.ImageMimeTypeError{}
  errMsg := "Unsupported mime type - Only .jpeg is allowed."
  if err.Error() != errMsg {
    t.Errorf("Expected %s but got %s", errMsg, err.Error())
  }
}

func TestRequestTimeoutError(t *testing.T) {
  err := image.RequestTimeoutError{100}
  errMsg := "No workers available - max idle time 100ms."
  if err.Error() != errMsg {
    t.Errorf("Expected %s but got %s", errMsg, err.Error())
  }
}
