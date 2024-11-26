package image

import "testing"

func TestImageMaxSizeExceededError(t *testing.T) {
  err := ImageMaxSizeExceededError{8192 * 1024}
  errMsg := "Image size exceeded - 8192 kB."
  if err.Error() != errMsg {
    t.Errorf("Expected %s but got %s", errMsg, err.Error())
  }
}

func TestMimeTypeError(t *testing.T) {
  err := ImageMimeTypeError{}
  errMsg := "Unsupported mime type - Only .jpeg is allowed."
  if err.Error() != errMsg {
    t.Errorf("Expected %s but got %s", errMsg, err.Error())
  }
}

func TestRequestTimeoutError(t *testing.T) {
  err := RequestTimeoutError{100}
  errMsg := "No workers available - max idle time 100ms"
  if err.Error() != errMsg {
    t.Errorf("Expected %s but got %s", errMsg, err.Error())
  }
}
