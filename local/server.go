package main

import (
	"flag"
	"net"
	"fmt"
	"sync"
	"io"
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
	//connService, err := net.Dial("tcp", ":"+*flagServicePort)
	//if err != nil {
	//	panic(err)
	//}
	//_ = connService
	fmt.Println("Connected to service ", *flagServicePort)
	//bufService := make(chan []byte, 0)
	//bufRemote := make(chan []byte, 0)

	reader  := bufio.NewReader(connRemote)
	for {
		b, err := reader.ReadBytes('\x03')
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("+++++++ From remote: ", string(b))
		connService, err := net.Dial("tcp", ":"+*flagServicePort)
		if err != nil {
			fmt.Println(err)
			break
		}
		go func() {
			n, err := connService.Write(b)
			fmt.Println("+++++++ Write to service: ", n, err)
		}()

		receivedData, err := ioutil.ReadAll(connService)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("+++++++ Received from service: ", string(receivedData))
		connService.Close()
		//bufService <- b
		//receivedData := <- bufRemote
		connRemote.Write(append(receivedData, '\x03'))
		reader.Reset(connRemote)
	}
	//b, err := ioutil.ReadAll(connRemote)
	//if err != nil {
	//	fmt.Println(err)
	//}



	//wg := sync.WaitGroup{}
	//for {
	//	b, err := ioutil.ReadAll(connRemote)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println("From remote: ", string(b))
	//
	//}
	//connCopy(connService, connRemote, &wg)
	//connCopy(connRemote, connService, &wg)
	//connService.Close()
	//connRemote.Close()

	//for {
	//	wg := sync.WaitGroup{}
	//	//wg.Add(2)
	//	connCopy(connService, connRemote, &wg)
	//	connCopy(connRemote, connService, &wg)
	//	connService.Close()
	//	connRemote.Close()
	//	//wg.Wait()
	//
	//}
}

func connCopy(dst net.Conn, src net.Conn, wg *sync.WaitGroup) {

	written, err := io.Copy(dst, src)
	fmt.Println(written, err) // test only
	//wg.Done()
}