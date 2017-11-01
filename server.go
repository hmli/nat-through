package main

import (
	"net"
	"fmt"
	"time"
	"io"
	"sync"
	"flag"
)
// Server side: -sport: dst, -cport:  给client side开的端口
// Client side: -sport: 上面的cport, -cport: 给用户开的端口
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
	//defer listenServer.Close()
	defer listener.Close()
	//go test(clientPort)
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
	}
}

func connCopy(dst net.Conn, src net.Conn, wg *sync.WaitGroup) {
	written, err := io.Copy(dst, src)
	fmt.Println(written, err)
	wg.Done()
	dst.Close()
}

func test(port string){

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <- ticker.C:
			conn, err := net.Dial("tcp", ":"+port)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(conn.Write([]byte("dddd")))

			//n, err := conn.Read([]byte("12bbb"))
			//fmt.Println(n, err)
		}
	}

}


//func pipe(connRev net.Conn, connSed net.Conn) {
//	var buf []byte
//	_, err := connRev.Write(buf)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	connSed.Read(buf)
//}

