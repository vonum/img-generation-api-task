package worker_test

import (
	"bytes"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/tutti-ch/backend-coding-task-template/image"
	"github.com/tutti-ch/backend-coding-task-template/worker"
)

func TestInitWorkers(t *testing.T) {
  nStart := runtime.NumGoroutine()

  worker.InitWorkers(3, "", make(<-chan image.Job))

  n := runtime.NumGoroutine()


  if nStart + 3 != n {
    t.Error("Expected 3 additional go routines to spawn")
  }
}

func TestSuccessfulJob(t *testing.T) {
  var buf bytes.Buffer
  log.SetOutput(&buf)

  c := make(chan image.Job)
  defer close(c)

  w := worker.NewWorker(1, "/tmp/smgtest/", c)
  go w.Run()

  imgBytes, _ := os.ReadFile("../testdata/testimage_small.jpg")
  id := "testimage"

  c <- image.Job{Id: id, Payload: imgBytes}
  <- time.After(10 * time.Millisecond)

  resultBytes, err := os.ReadFile("/tmp/smgtest/testimage.jpeg")
  if err != nil || len(resultBytes) == 0 {
    t.Error("Expected rescaled image to be generated: ", err)
  }

  logs := buf.String()

  if !strings.Contains(logs, "Job Started id=1 image=testimage") {
    t.Error("Expected Job started log to be present")
  }

  if !strings.Contains(logs, "Job Finished id=1 image=testimage") {
    t.Error("Expected Job finished log to be present")
  }

  if !strings.Contains(logs, "image_path=/tmp/smgtest/testimage.jpeg") {
    t.Error("Expected output path to be present")
  }
}

// There are more potential causes for rescaling to fail
// This covers only the empty byte slice case
func TestFailedJobRescaling(t *testing.T) {
  var buf bytes.Buffer
  log.SetOutput(&buf)

  c := make(chan image.Job)
  defer close(c)

  w := worker.NewWorker(1, "/tmp/smgtest/", c)
  go w.Run()

  imgBytes := []byte{}
  id := "testimage"
  c <- image.Job{Id: id, Payload: imgBytes}
  <- time.After(10 * time.Millisecond)

  logs := buf.String()

  if !strings.Contains(logs, "Job Started id=1 image=testimage") {
    t.Error("Expected Job started log to be present")
  }

  if !strings.Contains(logs, "Job Failed id=1 image=testimage") {
    t.Error("Expected Job finished log to be present")
  }

  if !strings.Contains(logs, "failed to decode image config") {
    t.Error("Expected reason to be present")
  }
}

func TestFailedJobMissingDir(t *testing.T) {
  var buf bytes.Buffer
  log.SetOutput(&buf)

  c := make(chan image.Job)
  defer close(c)

  w := worker.NewWorker(1, "/tmp/smgtestmissing/", c)
  go w.Run()

  imgBytes, _ := os.ReadFile("../testdata/testimage_small.jpg")
  id := "testimage"
  c <- image.Job{Id: id, Payload: imgBytes}
  <- time.After(10 * time.Millisecond)

  logs := buf.String()

  if !strings.Contains(logs, "Job Started id=1 image=testimage") {
    t.Error("Expected Job started log to be present")
  }

  if !strings.Contains(logs, "Job Failed id=1 image=testimage") {
    t.Error("Expected Job Failed log to be present")
  }

  if !strings.Contains(logs, "no such file or directory") {
    t.Error("Expected reason to be present")
  }
}
