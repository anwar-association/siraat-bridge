package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func handleConn(id int, newConn net.Conn) {
	defer newConn.Close()
	for {
		data, err := bufio.NewReader(newConn).ReadString('\n')
		if err != nil {
			log.Printf("Cannot read into buffer: %v", err)
		}
		fmt.Printf("ConnectionID %d says ==>\n%s \n", id, data)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	sigs := make(chan os.Signal)
	conns := make(chan net.Conn, 1000)
	done := make(chan struct{})
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- struct{}{}
	}()
	go func() {
		for {
			newConn := <-conns
			go handleConn(rand.Intn(1000000), newConn)
			fmt.Println("New Connection")
		}
	}()
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				fmt.Println("Error in accepting the connection")
				continue
			}
			conns <- conn
		}
	}()
	<-done
}
