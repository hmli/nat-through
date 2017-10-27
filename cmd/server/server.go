package main

import (
	"net"
	"fmt"
	"io"
)


func main() {
	clientPort := "6999"
	serverPort := "7000"

	connClient, err := net.Dial("tcp", ":"+clientPort)
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		panic(err)
	}
	defer connClient.Close()
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(conn, connClient)
		io.Copy(connClient, conn)
		conn.Close()

	}



}


