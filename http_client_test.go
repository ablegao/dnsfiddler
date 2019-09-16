package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {

	config = ReadConfig("./config.json")
	m.Run()
}

func Test_baidu(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.143 Mobile Safari/537.36")

	res, err := httpClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	buf, _ := ioutil.ReadAll(res.Body)
	t.Log(string(buf))
}

func Test_LocalIP(t *testing.T) {

}
