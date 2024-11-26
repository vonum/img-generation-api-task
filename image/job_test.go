package image_test

import (
	"bytes"
	"testing"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

func TestJob(t *testing.T) {
  payload := []byte{1, 2, 3, 4}
  job := image.Job{Id: "id", Payload: payload}

  if job.Id != "id" || !bytes.Equal(job.Payload, payload) {
    t.Error("Expected job to be initialized properly")
  }
}
