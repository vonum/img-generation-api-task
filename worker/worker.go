package worker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

type Worker struct {
  id int
  basePath string
  c <- chan image.Job
}

func NewWorker(id int, basePath string, c <- chan image.Job) *Worker {
  return &Worker{id, basePath, c}
}

func InitWorkers(n int, basePath string, c <- chan image.Job) {
  for i := 0; i < n; i++ {
    fmt.Println("Starting worker: ", i + 1)
    worker := NewWorker(i + 1, basePath, c)
    go worker.Run()
  }
}

func (w *Worker) Run() {
  for job := range w.c {
    startTime := time.Now()
    ctx := context.Background()

    w.LogJobStarted(job.Id, len(job.Payload))

    rescaledBytes, err := image.Rescale(ctx, job.Payload)
    if err != nil {
      w.LogJobFailed(job, err)
      continue
    }

    filename := filepath.Join(w.basePath, job.Id + ".jpeg",)
    if err = os.WriteFile(filename, rescaledBytes, 0644); err != nil {
      w.LogJobFailed(job, err)
      continue
    }

    endTime := time.Now()
    dur := endTime.Sub(startTime).Milliseconds()
    w.LogJobFinished(job.Id, len(job.Payload), len(rescaledBytes), dur, filename)
  }
}

func (w *Worker) LogJobStarted(jobId string, nBytes int) {
  slog.Info(
    "Job Started",
    "id", w.id,
    "image", jobId,
    "n_bytes", nBytes,
  )
}

func (w *Worker) LogJobFinished(
  jobId string,
  nBytes int,
  nOutputBytes int,
  durationMS int64,
  imgPath string,
) {
  slog.Info(
    "Job Finished",
    "id", w.id,
    "image", jobId,
    "n_bytes", nBytes,
    "output_n_bytes", nOutputBytes,
    "duration_ms", durationMS,
    "image_path", imgPath,
  )
}

func (w *Worker) LogJobFailed(job image.Job, err error) {
  slog.Error(
    "Job Failed",
    "id", w.id,
    "image", job.Id,
    "reason", err,
  )
}
