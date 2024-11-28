package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tutti-ch/backend-coding-task-template/image"
	"github.com/tutti-ch/backend-coding-task-template/worker"
)

func RunServiceAndWorkers(
  basePath string,
  nJobs int,
  wg *sync.WaitGroup,
  c chan image.Job,
  sigChan chan os.Signal,
) {
  worker.InitWorkers(nJobs, basePath, wg, c)
  server := NewImageServer(":3000", basePath, c)

  go func() {
    err := server.Run()

    if !errors.Is(err, http.ErrServerClosed) {
      log.Fatalf("Http server error: %v\n", err)
    }
    log.Println("Not serving new requests.")
  }()

  signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
  v := <- sigChan
  fmt.Println("\n", v)

  shutdownCtx, cancel := context.WithTimeout(
    context.Background(),
    2 * time.Second,
  )
  defer cancel()

  if err := server.Server.Shutdown(shutdownCtx); err != nil {
      log.Fatalf("HTTP shutdown error: %v", err)
  }

  close(c)
}
