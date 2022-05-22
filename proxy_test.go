package main

import (
	"sync"
	"testing"
)

func TestSC(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		startServer(7788)
		wg.Done()
	}()
	go func() {
		startClient(7887, "127.0.0.1:7788")
		wg.Done()
	}()
	wg.Wait()
}

func TestC(t *testing.T) {
	startClient(7887, "127.0.0.1:7788")
}

func TestS(t *testing.T) {
	startServer(7788)
}
