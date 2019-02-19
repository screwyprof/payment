package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type HttpTestClient struct {
	handler http.Handler
}

func (c *HttpTestClient) SendGetRequest(path string, resp interface{}) error {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	c.handler.ServeHTTP(w, req)

	if w.Body == nil {
		return fmt.Errorf("server response is nil")
	}

	return json.NewDecoder(w.Body).Decode(&resp)
}

func (c *HttpTestClient) SendPostRequest(path string, r interface{}, resp interface{}) error {
	body, _ := json.Marshal(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(body))
	c.handler.ServeHTTP(w, req)

	if w.Body == nil {
		return fmt.Errorf("server response is nil")
	}

	return json.NewDecoder(w.Body).Decode(&resp)
}
