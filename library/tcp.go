package library

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/dong568789/go-forward/library/util"
)

func ListenAndServer(listenAddr, forwardAddr string, wg *sync.WaitGroup) {
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		util.Log().Error("listen fail: %v", err)
		return
	}
	util.Log().Info("accept %s to %s\n", listenAddr, forwardAddr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			util.Log().Error("accept failed, err:", err)
			continue
		}
		go HandleRequest(conn, forwardAddr)
	}
	return
}

func HandleRequest(conn net.Conn, forwardAddr string) {
	d := net.Dialer{Timeout: time.Second * 5}
	proxy, err := d.Dial("tcp", forwardAddr)
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
	util.Log().Info("connection %s -> %s closed: readBytes %d, writeBytes %d, duration %s", src.RemoteAddr(), dst.RemoteAddr(), readBytes, writeBytes, time.Now().Sub(ts))
}
