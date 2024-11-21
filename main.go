package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func main() {
	panic(`not implemented`)
}

func MakeHandler(fsDirectory string, poolSize int) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_ = request.Body.Close()

		type Response struct {
			Error   string `json:"error,omitempty"`
			ImageID string `json:"image_id,omitempty"`
		}

		if !true {
			writer.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(writer).Encode(Response{
				Error:   "",
				ImageID: uuid.NewString(),
			})
		} else {
			const UnknownError = "unknown_error"
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(Response{
				Error:   UnknownError,
				ImageID: "",
			})
		}
	})
}
