package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"proxy/encrypt"
)

var (
	token     string
	serverSet *flag.FlagSet
	clientSet *flag.FlagSet
	port      int
	address   string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	parseFlags()
	if serverSet.Parsed() {
		startServer(uint16(port))
	}
	if clientSet.Parsed() {
		startClient(uint16(port), address)
	}
}

func parseFlags() {
	//server
	serverSet = flag.NewFlagSet("server", flag.ExitOnError)
	serverSet.IntVar(&port, "port", 6688, "server port")
	serverSet.StringVar(&token, "token", "1234567890123456", "token")
	//client
	clientSet = flag.NewFlagSet("client", flag.ExitOnError)
	clientSet.IntVar(&port, "port", 6688, "client port")
	clientSet.StringVar(&token, "token", "1234567890123456", "token")
	clientSet.StringVar(&address, "address", "127.0.0.1:6688", "server address")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n%s server [options...]\n%s client [options...]\n\n", os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
		serverSet.Usage()
		fmt.Println()
		clientSet.Usage()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
	} else {
		encrypt.RestKey(token)
		switch args[0] {
		case "server":
			_ = serverSet.Parse(args[1:])
		case "client":
			_ = clientSet.Parse(args[1:])
		}
	}
}

func errPrint(err error) bool {
	if err != nil {
		_ = log.Output(2, fmt.Sprintln(err))
	}
	return err != nil
}
