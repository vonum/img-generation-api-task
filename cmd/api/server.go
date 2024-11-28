package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/tutti-ch/backend-coding-task-template/api"
	"github.com/tutti-ch/backend-coding-task-template/image"
)

const BasePathEnv = "BASE_PATH"
const NWorkersEnv = "N_WORKERS"

func main() {
  nWorkers := runtime.NumCPU()
  basePath, ok := os.LookupEnv(BasePathEnv)
  if !ok {
    log.Fatal("Base path not set, export BASE_PATH env var.")
  }

  if _, err := os.Stat(basePath); err != nil {
    log.Fatalf("%s does not exist.", basePath)
  }

  nw, ok := os.LookupEnv(NWorkersEnv)
  if ok {
    nWorkers, _ = strconv.Atoi(nw)
  }

  sigChan := make(chan os.Signal, 1)
  wg := sync.WaitGroup{}
  c := make(chan image.Job)

  api.RunServiceAndWorkers(basePath, nWorkers, &wg, c, sigChan)

  log.Println("Waiting for workers to finish processing.")
  wg.Wait()
  log.Println("Graceful shutdown complete.")
}
