package image_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tutti-ch/backend-coding-task-template/image"
	"github.com/tutti-ch/backend-coding-task-template/worker"
)

func CreateMultiPart(mimeType string, data []byte, buf *bytes.Buffer) *multipart.Writer {
  multipartWriter := multipart.NewWriter(buf)

  h := make(textproto.MIMEHeader)
  h.Set("Content-Disposition", `form-data; name="image"; filename="image.jpeg"`)
  h.Set("Content-Type", mimeType)
  part, _ := multipartWriter.CreatePart(h)
  part.Write(data)
  multipartWriter.Close()

  return multipartWriter
}

func sendRequest(imgPath, mimeType string, c chan image.Job) *http.Response {
  imgBytes, _ := os.ReadFile(imgPath)

  var buf bytes.Buffer
  multipartWriter := CreateMultiPart(
    mimeType,
    imgBytes,
    &buf,
  )

  r := httptest.NewRequest(http.MethodPost, "/upload", &buf)
  w := httptest.NewRecorder()
  r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

  h := image.NewImageHandler(image.MaxImageSize, c)
  h.Rescale(w, r)

  return w.Result()
}

func TestRescaling(t *testing.T) {
  outputPath := "/tmp/smgtest/"
  c := make(chan image.Job)
  defer close(c)

  t.Run("Successful Rescaling", func(t *testing.T) {
    imgPath := "../testdata/testimage_small.jpg"
    go worker.InitWorkers(1, outputPath, c)

    res := sendRequest(imgPath, "image/jpeg", c)
    defer res.Body.Close()

    var response image.Response
    json.NewDecoder(res.Body).Decode(&response)

    if res.StatusCode != 200 {
      t.Errorf("Expected response with status code %d, but got %d.", 200, res.StatusCode)
    }

    if err := uuid.Validate(response.ImageID); err != nil {
      t.Errorf("Expected a valid uuid but got %s.", response.ImageID)
    }
  })

  t.Run("Image Too Large", func(t *testing.T) {
    imgPath := "../testdata/testimage_big.jpg"
    errMsg := "Image size exceeded - 8192 kB."
    go worker.InitWorkers(1, outputPath, c)

    res := sendRequest(imgPath, "image/jpeg", c)
    defer res.Body.Close()

    var response image.Response
    json.NewDecoder(res.Body).Decode(&response)

    if res.StatusCode != 413 {
      t.Errorf("Expected response with status code %d, but got %d.", 413, res.StatusCode)
    }

    if response.Error != errMsg {
      t.Errorf("Expected error message %s but got %s.", errMsg, response.Error)
    }
  })

  t.Run("Wrong Mime Type", func(t *testing.T) {
    imgPath := "../testdata/text.txt"
    errMsg := "Unsupported mime type - Only .jpeg is allowed."
    go worker.InitWorkers(1, outputPath, c)

    res := sendRequest(imgPath, "plain/text", c)
    defer res.Body.Close()

    var response image.Response
    json.NewDecoder(res.Body).Decode(&response)

    if res.StatusCode != 400 {
      t.Errorf("Expected response with status code %d, but got %d.", 400, res.StatusCode)
    }

    if response.Error != errMsg {
      t.Errorf("Expected error message %s but got %s.", errMsg, response.Error)
    }
  })

  t.Run("Timeout exceeded", func(t *testing.T) {
    imgPath := "../testdata/testimage_small.jpg"
    errMsg := "No workers available - max idle time 100ms."

    c := make(chan image.Job)
    defer close(c)

    startTime := time.Now()
    res := sendRequest(imgPath, "image/jpeg", c)
    defer res.Body.Close()
    endTime := time.Now().Sub(startTime)

    var response image.Response
    json.NewDecoder(res.Body).Decode(&response)

    if res.StatusCode != 429 {
      t.Errorf("Expected response with status code %d, but got %d.", 429, res.StatusCode)
    }

    if response.Error != errMsg {
      t.Errorf("Expected error message %s but got %s.", errMsg, response.Error)
    }

    if endTime.Milliseconds() < 100 {
      t.Error("Expected timeout to be after 100 milliseconds.")
    }
  })
}
