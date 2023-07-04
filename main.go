package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"
)

func TCPscanner(ports chan int, result chan int, wg *sync.WaitGroup) {

	for p := range ports {
		address := fmt.Sprintf("localhost:%d", p)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			result <- p
			fmt.Println(result)
			conn.Close()
			break
		} else {
			continue
		}

	}
	wg.Done()
}

func TLSscanner(result chan int, openports map[int]string) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	for r := range result {
		fmt.Println(r)
		address := fmt.Sprintf("localhost:%d", r)
		connTLS, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", address, conf)
		if err != nil {
			fmt.Println("1:", err)
			openports[r] = "no TLS"
			fmt.Println(openports)
			break
		} else {
			fmt.Println("print cert:", connTLS.ConnectionState().PeerCertificates)
			openports[r] = "TLS"
			fmt.Println(openports)
		}
		connTLS.Close()

	}
	fmt.Println("End of TLS scanner")
}

func main() {
	var wg sync.WaitGroup
	openports := make(map[int]string)
	ports := make(chan int, 100)
	result := make(chan int)

	go func() {
		fmt.Println("Taking ports")
		for i := 0; i <= 1024; i++ {
			ports <- i
		}
		fmt.Println("Close channel ports")
		close(ports)
	}()

	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go TCPscanner(ports, result, &wg)
		go TLSscanner(result, openports)
		//wg.Done()
	}
	wg.Wait()

	fmt.Println("goroutines out")
	//for i := 0; i < cap(ports); i++ {
	//	port := <-result
	//	if port == 0 {
	//		continue
	//	}
	//	return
	//}

	//defer close(result)
	//defer close(ports)

	fmt.Println(openports)
}
