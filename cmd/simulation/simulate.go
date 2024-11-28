package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"sync"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

// Not intedent to be properly engineered
// Just a tool for testing manually
func sendRequest(mimeType string, data []byte, wg *sync.WaitGroup) {
  url := "http://localhost:3000/upload"

  var buf bytes.Buffer
  multipartWriter := multipart.NewWriter(&buf)

  h := make(textproto.MIMEHeader)
  h.Set("Content-Disposition", `form-data; name="image"; filename="image.jpeg"`)
  h.Set("Content-Type", mimeType)
  part, _ := multipartWriter.CreatePart(h)
  part.Write(data)
  multipartWriter.Close()

  req, err := http.NewRequest(http.MethodPost, url, &buf)
  if err != nil {
    log.Fatal(err)
  }
  req.Header.Add("Content-Type", multipartWriter.FormDataContentType())


	client := &http.Client{}
  r, err := client.Do(req)
  if err != nil {
    fmt.Println("Failed sending request", err)
    wg.Done()
    return
  } else {
    defer r.Body.Close()

    var response image.Response
    json.NewDecoder(r.Body).Decode(&response)

    if r.StatusCode == http.StatusOK {
      fmt.Printf("200 - %s\n", response.ImageID)
    } else {
      fmt.Printf("%d - %s\n", r.StatusCode, response.Error)
    }
    wg.Done()
  }

}

func main() {
  validPath := "./testdata/testimage_small.jpg"
  tooLargePath := "./testdata/testimage_big.jpg"
  wrongMimeTypePath := "./testdata/text.txt"

  imgMimeType := "image/jpeg"
  txtMimeType := "plain/text"

  validBytes, _ := os.ReadFile(validPath)
  tooLargeBytes, _ := os.ReadFile(tooLargePath)
  wrongMimeTypeBytes, _ := os.ReadFile(wrongMimeTypePath)

  wg := sync.WaitGroup{}

  for i := 0; i < 5000; i++ {
    wg.Add(1)
    go func() {
      sendRequest(imgMimeType, validBytes, &wg)
    }()
  }

  for i := 0; i < 500; i++ {
    wg.Add(1)
    go func() {
      sendRequest(imgMimeType, tooLargeBytes, &wg)
    }()
  }

  for i := 0; i < 500; i++ {
    wg.Add(1)
    go func() {
      sendRequest(txtMimeType, wrongMimeTypeBytes, &wg)
    }()
  }

  wg.Wait()
}
