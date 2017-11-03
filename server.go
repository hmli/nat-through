package main

import (
	"net"
	"fmt"
	"io"
	"sync"
	"flag"
)
// local:内网中的服务 				-sport: 要转发的服务的端口		-cport: Remote的地址(端口)
// remote:远端服务器上的中转服务 	-sport: 上面的cport, -cport: 给用户开的端口
func main() {
	flagServerPort := flag.String("sport", "8000", "Server port number")
	flagClientPort := flag.String("cport", "6999", "Client port number")
	flag.Parse()
	clientPort := *flagClientPort
	localServerPort := *flagServerPort

	listener, err := net.Listen("tcp", ":"+clientPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		connServer, err := net.Dial("tcp", ":"+localServerPort)
		if err != nil {
			panic(err)
		}
		wg := sync.WaitGroup{}
		wg.Add(2)
		go connCopy(connServer, conn, &wg)
		go connCopy(conn, connServer, &wg)
		wg.Wait()
		conn.Close()
		connServer.Close()
	}
}

func connCopy(dst net.Conn, src net.Conn, wg *sync.WaitGroup) {
	written, err := io.Copy(dst, src)
	fmt.Println(written, err) // test only
	wg.Done()
}
