package image

import (
	"io"
	"net/http"
)

// very tightly coupled to the handler
// extracted just for simplicity
// improvement could be to use interfaces if possible for writer and request
func ReadBytes(w http.ResponseWriter, r *http.Request, maxBytes int64) ([]byte, error) {
  // this limit includes the header as well
  r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
  img, header, err := r.FormFile("image")

  // this should handle different types of errors
  // example: empty byte slice
  if err != nil {
    return nil, &ImageMaxSizeExceededError{
      MaxSizeBytes: maxBytes,
    }
  }

  mimeType, ok := header.Header["Content-Type"]
  if !ok || mimeType[0] != JPEGMimeType {
    return nil, &ImageMimeTypeError{}
  }

  imgBytes, _ := io.ReadAll(img)

  return imgBytes, nil
}
