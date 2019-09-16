package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

var conn *net.UDPConn
var dnsLogs *log.Logger

func QuestInHosts(hosts []string, question []dns.Question) (string, bool) {
	for _, quest := range question {
		for _, host := range hosts {
			if strings.Index(quest.Name, host) > -1 {
				return quest.Name, true
			}
		}
	}
	return "", false
}

func DNSServer() {
	dns_log := config.DNSLog //os.TempDir() + "/debug.server.dns.log"
	log.Println("DNS log file:", dns_log)
	log.Println("http://", config.LocalIP)

	f, err := os.OpenFile(dns_log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dnsLogs = log.New(f, "DNS-->  ", log.Ltime|log.Lshortfile)
	myHosts := config.RedirectHost //strings.Split(, ",")
	log.Println("hosts", myHosts)
	client := new(dns.Client)
	ip := net.ParseIP(config.LocalIP)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53, IP: ip})
	if err != nil {
		dnsLogs.Panic(err)
		return
	}
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		dnsLogs.Println("DNS Search:", addr)
		var m dns.Msg
		err := m.Unpack(buf)
		if err != nil {
			dnsLogs.Println("ERROR 1", err)
			continue
		}
		m.RecursionDesired = true
		if _, ok := QuestInHosts(myHosts, m.Question); ok || len(myHosts) == 0 {
			if len(m.Question) > 0 {
				rr, err := dns.NewRR(fmt.Sprintf("%s 100 IN A "+config.LocalIP, m.Question[0].Name))
				if err != nil {
					log.Println(err)
					continue
				}
				m.Answer = append(m.Answer, rr)
				dnsLogs.Println(m.String())
				out, _ := m.Pack()
				conn.WriteToUDP(out, addr)
				continue
			}
		}
		// log.Println("QUESTION .... ", m.Question)
		var r *dns.Msg
		var dnsAddr string
		for _, dnsAddr = range config.DNS {
			r, _, err = client.Exchange(&m, net.JoinHostPort(dnsAddr, "53"))
			if err == nil && r != nil && r.Rcode == dns.RcodeSuccess {
				break
			}
		}
		if r != nil {
			out, _ := r.Pack()
			// dnsLogs.Println(m.Id, r.String())
			conn.WriteToUDP(out, addr)
		}
	}
}
