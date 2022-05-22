package main

import (
	"fmt"
	"log"
	"net"
)

func startClient(clientPort uint16, serverAddress string) {
	address := fmt.Sprintf("127.0.0.1:%d", clientPort)
	log.Printf("start client on %s\n", address)
	listen, err := net.Listen("tcp", address)
	if errPrint(err) {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if errPrint(err) {
			continue
		}
		go client(conn, serverAddress)
	}

}

func client(client net.Conn, serverAddress string) {
	server, err := net.Dial("tcp", serverAddress)

	if errPrint(err) {
		return
	}
	go transformIoEncrypt(server, client)
	transformIoDecrypt(client, server)
	_ = client.Close()
	_ = server.Close()
}
