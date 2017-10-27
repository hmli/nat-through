package main

import (
	"net"
	"fmt"
	"time"
	"io"
	"sync"
)

func main() {
	clientPort := "6999"
	localServerPort := "8000"

	connServer, err := net.Dial("tcp", ":"+localServerPort)
	if err != nil {
		panic(err)
	}
	//listenServer, err := net.Listen("tcp", ":"+localServerPort)
	//if err != nil {
	//	panic(err)
	//}


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
		//conn2, err := listenServer.Accept()
		//if err != nil {
		//	fmt.Println("err-2:", err)
		//	continue
		//}
		//b, err := ioutil.ReadAll(conn)
		//fmt.Println(string(b), err)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go connCopy(connServer, conn, &wg)
		go connCopy(conn, connServer, &wg)
		wg.Wait()
		//written, err := io.Copy(conn2, conn)
		//fmt.Println(written, err)
		//w2, err := io.Copy(conn, conn2)
		//
		//fmt.Println(w2, err)
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

