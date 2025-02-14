package api_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/tutti-ch/backend-coding-task-template/api"
	"github.com/tutti-ch/backend-coding-task-template/image"
)

// Ideally the wait would be done inside RunServiceAndWorkers
// To simulate the behavour where requests and SIGTERM are coming from the outside
// This would be difficult to test
// Since multiple go routines would need to be synced
// One that is making the requests
// And one that is sending the SIGTERM or SIGINT signal
// This could lead to flakyness without proper testing
// Could look something like this
// go func() {SEND_IMAGE_JOBS}()
// go func() {SEND_SIGTERM}()
// RunServiceAndWorkers()
// DoTheCheks()
// To avoid this issue, it is assumed that the entrypoint would do the wait
// In this case, the entrypoint is /cmd/api/server.go
func TestRunServiceAndWorkers(t *testing.T) {
  basePath := t.TempDir()
  nJobs := 2
  wg := sync.WaitGroup{}

  var imgIds []string
  for i := 0; i < 10000; i++ {
    imgIds = append(imgIds, fmt.Sprintf("shutdowntest%d", i))
  }

  c := make(chan image.Job)
  sigChan := make(chan os.Signal, 1)

  imgBytes, err := os.ReadFile("../testdata/testimage_small.jpg")
  if err != nil {
    t.Fatal("Failed reading image bytes.")
  }

  go api.RunServiceAndWorkers(basePath, nJobs, &wg, c, sigChan)
  for _, imgId := range imgIds {
    c <- image.Job{Id: imgId, Payload: imgBytes}
  }
  sigChan <- os.Interrupt
  // assumed the caller would wait
  wg.Wait()

  for _, imgId := range imgIds {
    fpath := filepath.Join(basePath, imgId + ".jpeg")
    if _, err := os.Stat(fpath); err != nil {
      t.Errorf("Expected file with id %s to be generated.\n", imgId)
    }
  }
}
