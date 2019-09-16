package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpPost(t *testing.T) {
	server := httptest.NewServer(&serverHandler{})

	client := server.Client()
	req, _ := http.NewRequest("POST", "http://r1.snnd.co/v3/sdk-api", bytes.NewBufferString("xxxxxxxx fuck hhhh"))
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	buf, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == http.StatusBadRequest {
		t.Log(string(buf))
		t.Log(res.StatusCode)
	} else {
		t.Error(res.StatusCode)
	}
}
