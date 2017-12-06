package main

import (
	"net"
	"fmt"
	"io"
	"flag"
	"bufio"
	"bytes"
)
 //local:内网中的服务 				-sport: 要转发的服务的端口		-cport: Remote的地址(端口)
  // remote:远端服务器上的中转服务 	-sport: 上面的cport, -cport: 给用户开的端口
func main() {
	flagUserPort := flag.String("sport", "6998", "Server port number")
	flagServicePort := flag.String("cport", "6999", "Client port number")
	flag.Parse()
	bufUser :=  make(chan []byte, 0)
	bufLocal := make(chan []byte, 0)
	sListener, err := net.Listen("tcp", ":"+ *flagServicePort)
	if err != nil {
		panic(err)
	}
	defer sListener.Close()
	go listenLocal(sListener, bufUser, bufLocal)

	listener, err := net.Listen("tcp", ":"+*flagUserPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	listenUser(listener, bufLocal, bufUser)
}

func listenUser(listener net.Listener, bufLocal chan []byte, bufUser chan []byte) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		fmt.Println("conn accepted ", listener.Addr().String())
		go func(conn net.Conn){
			b := readFromConn(conn)
			bufLocal <- b
			receivedData := <-bufUser
			fmt.Println("======= received from service: \n", string(receivedData))
			reader := bytes.NewReader(receivedData)
			n, err := io.Copy(conn, reader)
			fmt.Println("conn writed", n, err)
			conn.Close()
		}(conn)

	}
}
func listenLocal(listener net.Listener, bufUser chan []byte, bufLocal chan []byte) {

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		fmt.Println("connected : ", listener.Addr().String())
		go func(){
			for {
				receivedData := <-bufLocal
				receivedData = append(receivedData, '\x03')
				n, err := conn.Write(receivedData)
				//fmt.Println("+++++writed: ", n, err)
				reader  := bufio.NewReader(conn)
				data, err := reader.ReadBytes('\x03')
				if err != nil {
					fmt.Println(err)
				}
				newData := make([]byte, 0)
				for _, b := range data {
					if b != '\x03' {
						newData = append(newData, b)
					}
				}
				bufUser <- newData
			}

		}()

	}
}

func readFromConn(conn net.Conn) (data []byte) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF{
			fmt.Println("err:", err)
			return
		}
		data = append(data, buf[:n]...)

		if n <1024  || err == io.EOF {
			break
		}
	}
	return
}

