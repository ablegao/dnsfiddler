package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

var (
	httpLog          *log.Logger
	httpClient       *http.Client
	HandlerMap       = map[string]HandlerFunc{}
	HandlerDomainMap = map[string]HandlerFunc{}
)

type HandlerFunc func(*http.Response, *http.Request, *Output)

type serverHandler struct {
	ssl bool
}

func homePage(w http.ResponseWriter, r *http.Request, outlog *Output) {
	if r.URL.Path == "/ca" {
		buf, err := ioutil.ReadFile(config.RootCA)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("content-type", "application/octet-stream")
		w.Header().Set("content-disposition", "attachment; filename=rootCA.pem")
		w.Write(buf)
		return
	}
	w.Write([]byte(`<a href="/ca">Install CA certificate</a><br/>`))
	w.Write([]byte(`<a href="https://` + config.LocalIP + `">Validation CA</a>`))
	w.Write([]byte("Client IP: \t\t" + r.RemoteAddr))
	w.Write([]byte("Server IP: \t\t" + config.LocalIP))
	w.Write([]byte("Proxy: \t\t" + config.ProxyServer))
}

func (sh *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	outlog := &Output{bytes.NewBuffer(nil), "\t"}
	outlog.Println("--------------------------------------------------------------------------------- START")
	defer func() {
		outlog.Println("--------------------------------------------------------------------------------- END")
		Logger.WithField("domain", r.Host).WithField("ip", config.LocalIP).WithField("path", r.URL.Path).Info(outlog.String() + "\n\n")
	}()
	if r.Host == config.LocalIP {
		homePage(w, r, outlog)
		return
	}
	httpPrefix := "http"
	if sh.ssl {
		// TODO: https 有个死循环...
		httpPrefix = "https"
	}

	outlog.NonPrefixPrintln("URL:\t", httpPrefix+"://"+r.Host+r.RequestURI)
	outlog.NonPrefixPrintln("SERVER-SSL:", sh.ssl)
	input := bytes.NewBuffer(nil)
	io.Copy(input, r.Body)

	req, _ := http.NewRequest(r.Method, fmt.Sprintf("%s://%s%s", httpPrefix, r.Host, r.RequestURI), input)
	outlog.NonPrefixPrintln("HEAD:S\t ================================")
	for k, m := range r.Header {
		outlog.NonPrefixPrintln(k, ":\t", strings.Join(m, ","))
		req.Header.Set(k, strings.Join(m, "; "))
	}

	// trace := &httptrace.ClientTrace{
	// 	GotConn: func(connInfo httptrace.GotConnInfo) {
	// 		outlog.NonPrefixPrintln("resolved to: %s", connInfo.Conn.RemoteAddr())
	// 	},
	// }
	// req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	outlog.NonPrefixPrintln("HEAD:E\t================================")
	res, err := httpClient.Do(req)
	if err != nil {
		outlog.Println("ERROR", err)
		return
	}
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, "; "))
	}
	outlog.Println("remote address:", res.StatusCode)
	if fun, ok := HandlerMap[r.URL.Path]; ok {
		fun(res, r, outlog)
		io.Copy(w, res.Body)
		return
	}
	if fun, ok := HandlerDomainMap[r.Host]; ok {
		fun(res, r, outlog)
		io.Copy(w, res.Body)
		return
	}
	if config.ShowBody {

		outlog.NonPrefixPrintln("BODY:")
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, res.Body)
		outlog.NonPrefixPrintln(buf.String())
		if r.Host == "www.baidu.com" {
			return
		}
		w.Write(buf.Bytes())
	} else {
		io.Copy(w, res.Body)
	}
}

func UserDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	//设置自定义DNS地址和端口
	log.Println("httpDNSSSSSSSSS;", config.DNS[0])
	return d.DialContext(ctx, "udp", config.DNS[0]+":53")
}

func HTTPServer() {

	log.Println("https://" + config.LocalIP + ":443")
	log.Println("http://" + config.LocalIP + ":80")
	go func() {
		http.ListenAndServe(":11111", nil)
	}()
	go func() {
		cfg := &tls.Config{
			InsecureSkipVerify: true,
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return fetchCert(info.ServerName, strings.Split(info.Conn.LocalAddr().String(), ":")[0])
			},
		}

		srv := &http.Server{
			Addr:      ":443",
			Handler:   &serverHandler{true},
			TLSConfig: cfg,
		}
		err := srv.ListenAndServeTLS("", "")
		httpLog.Panic(err)
	}()

	httpLog = log.New(os.Stdout, "HTTP --> ", log.Ltime|log.Lshortfile)
	err := http.ListenAndServe(config.LocalIP+":80", &serverHandler{false})
	if err != nil {
		httpLog.Panic(err)
	}
}
