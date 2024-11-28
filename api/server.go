package api

import (
	"net/http"

	"github.com/tutti-ch/backend-coding-task-template/image"
)

type ImageServer struct {
  Server *http.Server
  basePath string
  c chan <- image.Job
}

func NewImageServer(port string, basePath string, c chan <- image.Job) *ImageServer {
  Server := &http.Server{Addr: port}

  return &ImageServer{Server, basePath, c}
}

func (is *ImageServer) Run() error {
  imageHandler := image.NewImageHandler(image.MaxImageSize, is.c)

  http.HandleFunc("/health", health)
  http.HandleFunc("/upload", imageHandler.Rescale)

  err := is.Server.ListenAndServe()
  return err
}

func health(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}
