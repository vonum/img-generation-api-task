package image

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const MaxImageSize = 8192 * 1024
const JPEGMimeType = "image/jpeg"
const MaxIdleTimeMS = 100

type ImageHandler struct {
  MaxFileSize int64
  c chan <- Job
}

func NewImageHandler(maxFileSize int64, c chan <- Job) *ImageHandler {
  return &ImageHandler{maxFileSize, c}
}

type Sizer interface {
  Size() int64
}

func (h *ImageHandler) Rescale(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
    return
  }

  imgBytes, err := ReadBytes(w, r, h.MaxFileSize)
  if err != nil {
    if maxSizeErr, ok := err.(*ImageMaxSizeExceededError); ok {
      w.WriteHeader(http.StatusRequestEntityTooLarge)
      _ = json.NewEncoder(w).Encode(Response{
        Error: maxSizeErr.Error(),
        ImageID: "",
      })
    }

    if mimeTypeErr, ok := err.(*ImageMimeTypeError); ok {
      w.WriteHeader(http.StatusBadRequest)
      _ = json.NewEncoder(w).Encode(Response{
        Error: mimeTypeErr.Error(),
        ImageID: "",
      })
    }

    return
  }

  imgId := uuid.NewString()
  select {
  case h.c <- Job{Id: imgId, Payload: imgBytes}:
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(Response{
      Error:   "",
      ImageID: imgId,
    })
  case <- time.After(MaxIdleTimeMS * time.Millisecond):
    timeoutErr := RequestTimeoutError{MaxIdleTimeMS}
    w.WriteHeader(http.StatusTooManyRequests)
    _ = json.NewEncoder(w).Encode(Response{
      Error: timeoutErr.Error(),
      ImageID: "",
    })
  }
}
