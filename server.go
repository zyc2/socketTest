package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

func startServer(serverPort uint16) {
	address := fmt.Sprintf("0.0.0.0:%d", serverPort)
	log.Printf("start server on %s\n", address)
	listen, err := net.Listen("tcp", address)
	if errPrint(err) {
		log.Panic(err)
	}
	for {
		conn, err := listen.Accept()
		if errPrint(err) {
			continue
		}
		go serve(conn)
	}

}

func serve(client net.Conn) {
	defer func(client net.Conn) {
		_ = client.Close()
	}(client)
	buf, err := DecryptUnpackOne(client)
	if checkNetError(err) {
		return
	}
	var method, host, address string
	_, err = fmt.Sscanf(string(buf[:bytes.IndexByte(buf[:], '\n')]), "%s%s", &method, &host)
	if errPrint(err) {
		return
	}
	hostPortURL, err := url.Parse(host)
	if errPrint(err) {
		return
	}
	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	server, err := net.Dial("tcp", address)
	if errPrint(err) {
		return
	}

	if method == "CONNECT" {
		_, _ = client.Write(EncryptPack([]byte("HTTP/1.1 200 Connection established\r\n\r\n")))
	} else {
		cnt := 0
		for i, b := range buf {
			if b == '/' {
				cnt++
				if cnt == 3 {
					_, err = server.Write(append(buf[:len(method)+1], buf[i:]...))
					if errPrint(err) {
						return
					}
					break
				}
			}
		}
	}
	go transformIoDecrypt(server, client)
	transformIoEncrypt(client, server)
}
