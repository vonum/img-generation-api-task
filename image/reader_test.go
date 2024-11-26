package image_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

func TestReadBytesAboveThreshold(t *testing.T) {
  imgBytes, _ := os.ReadFile("../testdata/testimage_big.jpg")

  var buf bytes.Buffer
  multipartWriter := CreateMultiPart(
    "image/jpeg",
    imgBytes,
    &buf,
  )

  r := httptest.NewRequest(http.MethodPost, "/upload", &buf)
  w := httptest.NewRecorder()
  r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

  _, err := image.ReadBytes(w, r, image.MaxImageSize)

  errMsg := "Image size exceeded - 8192 kB."
  if err == nil {
    t.Errorf("Expected to return error: %s", errMsg)
  }

  if err.Error() != errMsg {
    t.Errorf("Expected error: %s, but got %s", errMsg, err.Error())
  }
}

func TestReadBytesWrongMimeType(t *testing.T) {
  textBytes, _ := os.ReadFile("../testdata/text.txt")

  var buf bytes.Buffer
  multipartWriter := CreateMultiPart(
    "plain/text",
    textBytes,
    &buf,
  )

  r := httptest.NewRequest(http.MethodPost, "/upload", &buf)
  w := httptest.NewRecorder()
  r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

  _, err := image.ReadBytes(w, r, image.MaxImageSize)

  errMsg := "Unsupported mime type - Only .jpeg is allowed."
  if err == nil {
    t.Errorf("Expected to return error: %s", errMsg)
  }

  if err.Error() != errMsg {
    t.Errorf("Expected error: %s, but got %s", errMsg, err.Error())
  }
}

func TestReadBytesOk(t *testing.T) {
  imgBytes, _ := os.ReadFile("../testdata/testimage_smalljpg")

  var buf bytes.Buffer
  multipartWriter := CreateMultiPart(
    "image/jpeg",
    imgBytes,
    &buf,
  )

  r := httptest.NewRequest(http.MethodPost, "/upload", &buf)
  w := httptest.NewRecorder()
  r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

  resultBytes, err := image.ReadBytes(w, r, image.MaxImageSize)

  if err != nil {
    t.Error("Expected to read bytes successfully")
  }

  if len(imgBytes) != len(resultBytes) {
    t.Error("Expected result bytes to be equal to the request bytes")
  }
}
