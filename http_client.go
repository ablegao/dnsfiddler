package main

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

func initClient() {

	// 配置使用的代理和dns
	tran := http.Transport{
		Proxy: http.ProxyFromEnvironment,

		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
			//优先使用Go内置DNS查询，替换DNS
			Resolver: &net.Resolver{
				PreferGo: true,
				Dial:     UserDNSDialer,
			},
		}).DialContext,
		// MaxIdleConns:          20,
		// IdleConnTimeout:       90 * time.Second,
		// TLSHandshakeTimeout:   10 * time.Second,
		// ExpectContinueTimeout: 1 * time.Second,
	}

	if config != nil && config.ProxyServer != "" {
		url, err := url.Parse(config.ProxyServer)
		if err != nil {
			panic(err)
		}
		tran.Proxy = http.ProxyURL(url)
	}
	httpClient = &http.Client{
		Transport: &tran,
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	log.Println(via, req)
		//	return http.ErrUseLastResponse
		//},
	}

}
