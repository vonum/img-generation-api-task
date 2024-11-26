package image

import (
	"bytes"
	"mime/multipart"
	"net/http"
  "net/textproto"
	"net/http/httptest"
	"os"
	"testing"
)

func createMultiPart(mimeType string, data []byte, buf *bytes.Buffer) *multipart.Writer {
  multipartWriter := multipart.NewWriter(buf)

  h := make(textproto.MIMEHeader)
  h.Set("Content-Disposition", `form-data; name="image"; filename="image.jpeg"`)
  h.Set("Content-Type", mimeType)
  part, _ := multipartWriter.CreatePart(h)
  part.Write(data)
  multipartWriter.Close()

  return multipartWriter
}

func TestReadBytesAboveThreshold(t *testing.T) {
  imgBytes, _ := os.ReadFile("../testdata/testimage_big.jpg")

  var buf bytes.Buffer
  multipartWriter := createMultiPart(
    "image/jpeg",
    imgBytes,
    &buf,
  )

  r := httptest.NewRequest(http.MethodPost, "/upload", &buf)
  w := httptest.NewRecorder()
  r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

  _, err := ReadBytes(w, r, MaxImageSize)

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
  multipartWriter := createMultiPart(
    "plain/text",
    textBytes,
    &buf,
  )

  r := httptest.NewRequest(http.MethodPost, "/upload", &buf)
  w := httptest.NewRecorder()
  r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

  _, err := ReadBytes(w, r, MaxImageSize)

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
  multipartWriter := createMultiPart(
    "image/jpeg",
    imgBytes,
    &buf,
  )

  r := httptest.NewRequest(http.MethodPost, "/upload", &buf)
  w := httptest.NewRecorder()
  r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

  resultBytes, err := ReadBytes(w, r, MaxImageSize)

  if err != nil {
    t.Error("Expected to read bytes successfully")
  }

  if len(imgBytes) != len(resultBytes) {
    t.Error("Expected result bytes to be equal to the request bytes")
  }
}
