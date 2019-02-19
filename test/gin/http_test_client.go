package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type HttpError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e HttpError) Error() string {
	return e.Message
}

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

	bodyBytes, _ := ioutil.ReadAll(w.Body)
	if err := c.tryUnmarshalAsError(bodyBytes); err != nil {
		return err
	}

	//reset the response body to the original unread state
	w.Body = bytes.NewBuffer(bodyBytes)

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

func (c *HttpTestClient) tryUnmarshalAsError(body []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.DisallowUnknownFields()

	resp := HttpError{}
	if err := dec.Decode(&resp); err == nil {
		return resp
	}

	return nil
}
