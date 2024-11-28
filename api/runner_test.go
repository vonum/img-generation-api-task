package api_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/tutti-ch/backend-coding-task-template/api"
	"github.com/tutti-ch/backend-coding-task-template/image"
)

// Ideally the wait would be done inside RunServiceAndWorkers
// This would be difficult to test
// Since multiple go routines would need to be synced
// One that is making the requests
// And one that is sending the SIGTERM or SIGINT signal
// This could lead to flaky tests
// To avoid this issue, it is assumed that the entrypoint would do the wait
// In this case, the entrypoint is /cmd/api/server.go
func TestRunServiceAndWorkers(t *testing.T) {
  basePath := "/tmp/smgtest/"
  nJobs := 2
  wg := sync.WaitGroup{}

  c := make(chan image.Job)
  sigChan := make(chan os.Signal, 1)
  imgIds := []string{"shutdowntest1", "shutdowntest2", "shutdowntest3"}

  imgBytes, err := os.ReadFile("../testdata/testimage_small.jpg")
  if err != nil {
    t.Fatal("Failed reading image bytes.")
  }

  go api.RunServiceAndWorkers(basePath, nJobs, &wg, c, sigChan)
  for _, imgId := range imgIds {
    c <- image.Job{Id: imgId, Payload: imgBytes}
  }
  sigChan <- os.Interrupt
  wg.Wait()

  for _, imgId := range imgIds {
    fpath := filepath.Join(basePath, imgId + ".jpeg")
    if _, err := os.Stat(fpath); err != nil {
      t.Errorf("Expected file with id %s to be generated.\n", imgId)
    }
  }
}
