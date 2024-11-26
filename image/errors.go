package image

import "fmt"

type ImageMaxSizeExceededError struct {
  MaxSizeBytes int64
}
type ImageMimeTypeError struct {}
type RequestTimeoutError struct {
  MaxIdleTimeMS int
}

func (e *ImageMaxSizeExceededError) Error() string {
  kb := e.MaxSizeBytes / 1024
  return fmt.Sprintf("Image size exceeded - %d kB.", kb)
}

func (e *ImageMimeTypeError) Error() string {
  return "Unsupported mime type - Only .jpeg is allowed."
}

func (e *RequestTimeoutError) Error() string {
  return fmt.Sprintf("No workers available - max idle time %dms", e.MaxIdleTimeMS)
}
