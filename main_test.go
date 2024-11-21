package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func makeMultipartRequest(t *testing.T, data []byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer

	mw := multipart.NewWriter(&buf)
	if err := mw.SetBoundary(`xYzZY`); err != nil {
		t.Fatal(err)
	}

	w, err := mw.CreateFormFile(`image`, `image.jpg`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	err = mw.Close()
	if err != nil {
		t.Fatal(err)
	}

	return &buf, mw.FormDataContentType()
}

// TestMakeHandler is a very basic test that does not cover all of the specification and should
// only serve as a starting point.
func TestMakeHandler(t *testing.T) {
	fsDirectory, err := os.MkdirTemp(``, `test-make-handler`)
	if err != nil {
		t.Fatalf(`failed to create temporary directory`)
	}

	defer os.RemoveAll(fsDirectory)

	handler := MakeHandler(fsDirectory, 1)

	t.Run(`small image`, func(t *testing.T) {
		inputData, err := os.ReadFile(`testdata/testimage_small.jpg`)
		if err != nil {
			t.Fatal(`failed to read testdata file`)
		}

		buf, contentType := makeMultipartRequest(t, inputData)
		req := httptest.NewRequest(`POST`, `/upload`, buf)
		req.Header.Set(`Content-Type`, contentType)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf(`unexpected status code %v`, w.Code)
		}
	})

	t.Run(`too big image`, func(t *testing.T) {
		inputData, err := os.ReadFile(`testdata/testimage_big.jpg`)
		if err != nil {
			t.Fatal(`failed to read testdata file`)
		}

		buf, contentType := makeMultipartRequest(t, inputData)
		req := httptest.NewRequest(`POST`, `/upload`, buf)
		req.Header.Set(`Content-Type`, contentType)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf(`unexpected status code %v`, w.Code)
		}
	})
}
