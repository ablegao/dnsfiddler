package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	// "golang.org/x/net/dns/dnsmessage"
)

func init() {
	if len(os.Args) > 1 {
		config = ReadConfig(os.Args[1])
	} else {
		config = ReadConfig("./config.json")
	}

	initClient()
}

var (
	config *Config
)

func listenSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sigs:
		fmt.Println("exitapp,sigs:", sigs)
		os.Exit(0)
	}

}

func main() {
	cmdPath := os.Args[0]
	log.Println(cmdPath)
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(i)
	if i == 0 {
		go HTTPServer()
		go DNSServer()
		listenSignal()
	} else {
		log.Println("This program must be run as root! (sudo)")
		log.Println("Need to exit should press Ctrl + c !!!")
		args := append([]string{cmdPath}, flag.Args()...)
		startByRoot := exec.Command("sudo", args...)
		startByRoot.Stdout = os.Stdout
		startByRoot.Stderr = os.Stderr
		startByRoot.Start()
		startByRoot.Wait()
	}
}
