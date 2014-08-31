package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
 *  API Helpers
 */

func return404(w http.ResponseWriter) {
	http.Error(w, "404 page not found", 404)
}

/*
 *  API Testing Helpers
 */

type HandleTester func(
	method string,
	params string,
) *httptest.ResponseRecorder

// Given the current test runner and an http.Handler, generate a
// HandleTester which will test its the handler against its input

func GenerateHandleTester(
	t *testing.T,
	handleFunc http.Handler,
) HandleTester {

	// Given a method type ("GET", "POST", etc) and
	// parameters, serve the response against the handler
	// and return the ResponseRecorder.

	return func(
		method string,
		params string,
	) *httptest.ResponseRecorder {

		req, err := http.NewRequest(
			method,
			"",
			strings.NewReader(params),
		)
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set(
			"Content-Type",
			"application/json",
		)
		req.Body.Close()
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}

// From Matt Aimonetti's blog post:
// matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
// Creates a new file upload http request with optional extra params
func newFileUploadRequest(
	paramName string,
	path string,
	params map[string]string,
) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
