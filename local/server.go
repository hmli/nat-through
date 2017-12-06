package main

import (
	"flag"
	"net"
	"fmt"
	"bufio"
	"io/ioutil"
)

func main() {
	flagServicePort := flag.String("sport", "8000", "Server port number")
	flagRemotePort := flag.String("rport", "6999", "Client port number")
	flag.Parse()
	connRemote, err := net.Dial("tcp", ":"+*flagRemotePort)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to remote ", *flagRemotePort)
	fmt.Println("Connected to service ", *flagServicePort)

	reader  := bufio.NewReader(connRemote)
	for {
		b, err := reader.ReadBytes('\x03')
		if err != nil {
			fmt.Println(err)
		}
		connService, err := net.Dial("tcp", ":"+*flagServicePort)
		if err != nil {
			fmt.Println(err)
			break
		}
		go func() {
			n, err := connService.Write(b)
			fmt.Println("Write to service: ", n, err)
		}()

		receivedData, err := ioutil.ReadAll(connService)
		if err != nil {
			fmt.Println(err)
		}
		connService.Close()
		connRemote.Write(append(receivedData, '\x03'))
		reader.Reset(connRemote)
	}

}

