package worker_test

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/tutti-ch/backend-coding-task-template/image"
	"github.com/tutti-ch/backend-coding-task-template/worker"
)

func runWorker(imgId string, imgBytes []byte, outputPath string, buf bytes.Buffer) string {
  log.SetOutput(&buf)

  wg := sync.WaitGroup{}
  c := make(chan image.Job)

  w := worker.NewWorker(1, outputPath, &wg, c)
  go w.Run()

  c <- image.Job{Id: imgId, Payload: imgBytes}
  close(c)
  wg.Wait()

  return buf.String()
}

func TestInitWorkers(t *testing.T) {
  nStart := runtime.NumGoroutine()

  wg := sync.WaitGroup{}
  c := make(chan image.Job)
  defer close(c)

  worker.InitWorkers(3, "", &wg, c)

  n := runtime.NumGoroutine()

  if nStart + 3 != n {
    t.Error("Expected 3 additional go routines to spawn")
  }
}

func TestRun(t *testing.T) {
  imgId := "testimage"

  imgPath := "../testdata/testimage_small.jpg"
  imgBytes, err := os.ReadFile(imgPath)

  if err != nil {
    t.Fatal(err)
  }

  t.Run("Successful Run", func(t *testing.T) {
    outputPath:= t.TempDir()
    outputImagePath := filepath.Join(outputPath, imgId + ".jpeg")

    var buf bytes.Buffer
    log.SetOutput(&buf)

    logs := runWorker(imgId, imgBytes, outputPath, buf)

    resultBytes, err := os.ReadFile(outputImagePath)
    if err != nil || len(resultBytes) == 0 {
      t.Error("Expected rescaled image to be generated: ", err)
    }

    if !strings.Contains(logs, "Job Started id=1 image=testimage") {
      t.Error("Expected Job started log to be present")
    }

    if !strings.Contains(logs, "Job Finished id=1 image=testimage") {
      t.Error("Expected Job finished log to be present")
    }

    if !strings.Contains(logs, "image_path=" + outputImagePath) {
      t.Error("Expected output path to be present")
    }
  })

  // There are more potential causes for rescaling to fail
  // This covers worker behaviour when rescaling returns an error
  // By the empty byte slice case
  // Rescaling is tested separately
  t.Run("Failed Rescaling", func(t *testing.T) {
    outputPath:= t.TempDir()
    outputImagePath := filepath.Join(outputPath, imgId + ".jpeg")

    var buf bytes.Buffer
    log.SetOutput(&buf)

    imgBytes := []byte{}

    logs := runWorker(imgId, imgBytes, outputPath, buf)

    _, err := os.ReadFile(outputImagePath)
    if err == nil {
      t.Error("Expected image to not be generated")
    }

    if !strings.Contains(logs, "Job Started id=1 image=testimage") {
      t.Error("Expected Job started log to be present")
    }

    if !strings.Contains(logs, "Job Failed id=1 image=testimage") {
      t.Error("Expected Job finished log to be present")
    }

    if !strings.Contains(logs, "failed to decode image config") {
      t.Error("Expected reason to be present")
    }
  })

  t.Run("Missing Directory", func(t *testing.T) {
    outputPath:= t.TempDir()
    outputImagePath := filepath.Join(outputPath, imgId + ".jpeg")

    var buf bytes.Buffer
    log.SetOutput(&buf)

    logs := runWorker(imgId, imgBytes, "/tmp/smgtestmissing/", buf)

    _, err := os.ReadFile(outputImagePath)
    if err == nil {
      t.Error("Expected image to not be generated")
    }

    if !strings.Contains(logs, "Job Started id=1 image=testimage") {
      t.Error("Expected Job started log to be present")
    }

    if !strings.Contains(logs, "Job Failed id=1 image=testimage") {
      t.Error("Expected Job Failed log to be present")
    }

    if !strings.Contains(logs, "no such file or directory") {
      t.Error("Expected reason to be present")
    }
  })
}
