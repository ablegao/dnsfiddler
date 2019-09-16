package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

type Config struct {
	LocalIP      string   `json:"local_ip"`
	ProxyServer  string   `json:"proxy_server"`
	RedirectHost []string `json:"redirect_host"`
	V1Time       int      `json:"v1time"`
	ShowBody     bool     `json:"show_body"`
	DNSLog       string   `json:"dns_log"`
	DNS          []string `json:"dns_server"`
	RootCA       string   `json:"root_cert"`
	RootKey      string   `json:"root_key"`
	Graylog      string   `json:"graylog"`
}

func ReadConfig(p string) *Config {
	f, err := ioutil.ReadFile(p)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println(string(f))
	conf := new(Config)
	err = json.Unmarshal(f, conf)
	if err != nil {
		panic(err)
	}
	if conf.LocalIP == "" {

		conf.LocalIP = GetIP()
	}

	if conf.Graylog != "" {
		hook := graylog.NewGraylogHook(conf.Graylog, map[string]interface{}{"server": "dns"})
		Logger.AddHook(hook)
	}
	return conf
}
