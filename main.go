package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var address = "localhost:"

func TCPscanner(ports, results chan int) {

	for p := range ports {
		start := time.Now()
		conn, err := net.Dial("tcp", address+strconv.Itoa(p))
		if err != nil {
			results <- 0
			continue
		}

		conn.Close()
		results <- p
		fmt.Println(p)
		fmt.Println("Working time:", time.Now().Sub(start))
	}

}

func TLSscanner(openports map[int]string, wg *sync.WaitGroup, mu *sync.Mutex) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	for r := range openports {
		connTLS, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", address+strconv.Itoa(r), conf)
		if err != nil {
			mu.Lock()
			openports[r] = "no TLS"
			mu.Unlock()
			continue
		} else {
			mu.Lock()
			fmt.Println("print cert:", connTLS.ConnectionState().PeerCertificates)
			openports[r] = "TLS"
			mu.Unlock()
		}
		connTLS.Close()
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	ports := make(chan int, 1024)
	results := make(chan int)

	openports := make(map[int]string)

	mu.Lock()

	start := time.Now()

	fmt.Println("Start TCP scanner")
	for i := 0; i < cap(ports); i++ {
		go TCPscanner(ports, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports[port] = ""
		}
	}

	fmt.Println("End TCP scanner")
	mu.Unlock()

	for i := 0; i < len(openports); i++ {
		wg.Add(1)
		go TLSscanner(openports, &wg, &mu)

	}
	wg.Wait()

	fmt.Println("Total operating time:", time.Since(start))
	close(ports)
	close(results)

	fmt.Println(openports)

}
