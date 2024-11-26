package image

import (
  "io"
  "net/http"
)

func ReadBytes(w http.ResponseWriter, r *http.Request, maxBytes int64) ([]byte, error) {
  // this limit includes the header as well
  r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
  img, header, err := r.FormFile("image")

  if err != nil {
    return nil, &ImageMaxSizeExceededError{
      MaxSizeBytes: maxBytes,
    }
  }

  // this would allow unset mime type to be rescalled
  // clients set these headers, so the other option
  // is too check bytes manually for validity
  mimeType, ok := header.Header["Content-Type"]
  if ok && mimeType[0] != JPEGMimeType {
    return nil, &ImageMimeTypeError{}
  }

  imgBytes, _ := io.ReadAll(img)

  return imgBytes, nil
}
