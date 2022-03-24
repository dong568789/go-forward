package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	debug   = true
	version string
)

func ParseArgs() (string, string) {
	listenAddr := flag.String("l", "", "listen address")
	forwardAddr := flag.String("f", "", "forwarding address")
	flagVersion := flag.Bool("v", false, "print version")
	flag.Parse()
	if *flagVersion {
		fmt.Println("version:", version)
		os.Exit(0)
	}
	if *forwardAddr == "" || *listenAddr == "" {
		flag.Usage()
		os.Exit(0)
	}

	return *listenAddr, *forwardAddr
}

func listenAndServer(listenAddr, forwardAddr string) {
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("accept %s to %s\n", listenAddr, forwardAddr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		log.Printf("HandleRequest\n")
		go HandleRequest(conn, forwardAddr)
	}
}

func HandleRequest(conn net.Conn, forwardAddr string) {
	d := net.Dialer{Timeout: time.Second * 5}

	proxy, err := d.Dial("tcp", forwardAddr)
	log.Printf("proxy\n")

	if err != nil {
		log.Printf("%s -> %s failed\n", conn.RemoteAddr(), forwardAddr)
		conn.Close()
		return
	}

	pipe(conn, proxy)
}

func pipe(src, dst net.Conn) {
	var readBytes int64
	var writeBytes int64
	ts := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)

	closeFun := func(err error) {
		dst.Close()
		src.Close()
	}
	go func() {
		wg.Done()
		n, err := io.Copy(dst, src)
		readBytes += n
		closeFun(err)
	}()

	n, err := io.Copy(src, dst)
	writeBytes += n
	closeFun(err)
	wg.Wait()
	log.Printf("connection %s -> %s closed: readBytes %d, writeBytes %d, duration %s", src.RemoteAddr(), dst.RemoteAddr(), readBytes, writeBytes, time.Now().Sub(ts))
}

func init() {
	if !debug {
		file := "logs/" + time.Now().Format("2006-01-02") + ".log"
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0766)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logFile)
	}
	return
}

func main() {
	listenAddr, forwardAddr := ParseArgs()
	listenAndServer(listenAddr, forwardAddr)
}
