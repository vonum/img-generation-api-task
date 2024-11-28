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
const NJobsEnv = "N_JOBS"

func main() {
  nJobs := runtime.NumCPU()
  basePath, ok := os.LookupEnv(BasePathEnv)
  if !ok {
    panic("Base path not set, export BASE_PATH env var.")
  }

  nj, ok := os.LookupEnv(NJobsEnv)
  if ok {
    nJobs, _ = strconv.Atoi(nj)
  }

  sigChan := make(chan os.Signal, 1)
  wg := sync.WaitGroup{}
  c := make(chan image.Job)

  api.RunServiceAndWorkers(basePath, nJobs, &wg, c, sigChan)

  log.Println("Waiting for workers to finish processing.")
  wg.Wait()
  log.Println("Graceful shutdown complete.")
}
