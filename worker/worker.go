package worker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

type Worker struct {
  Id int
  basePath string
  c <- chan image.Job
}

func InitWorkers(n int, basePath string, c <- chan image.Job) {
  for i := 0; i < n; i++ {
    fmt.Println("Starting worker: ", i + 1)
    worker := Worker{Id: i + 1, basePath: basePath, c: c}
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

    filename := fmt.Sprintf("%s%s.jpeg", w.basePath, job.Id)
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
    "id", w.Id,
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
    "id", w.Id,
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
    "id", w.Id,
    "image", job.Id,
    "reason", err,
  )
}
