package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/tutti-ch/backend-coding-task-template/api"
	"github.com/tutti-ch/backend-coding-task-template/image"
	"github.com/tutti-ch/backend-coding-task-template/worker"
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

  c := make(chan image.Job)

  fmt.Println(basePath)
  fmt.Println(nJobs)

  worker.InitWorkers(nJobs, basePath, c)

  server := api.NewImageServer(basePath, c)
  server.Run(":3000")

  fmt.Println(basePath)
}
