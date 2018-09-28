package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	netIface  string
	dir       string
	localPort int
	ssl       bool
)

func getAddrFromIfaceName(iface string) string {
	var addr string

	if iface == "any" {
		addr = "0.0.0.0"
	} else {
		ief, err := net.InterfaceByName(iface)
		if err != nil {
			log.Fatal(err)
			log.Printf("Interface %s not found\n", iface)
			os.Exit(1)
		}
		addrs, err := ief.Addrs()
		if err != nil {
			log.Fatal(err)
			log.Printf("Error retrieving address for %s interface\n", iface)
			os.Exit(2)
		}
		addr = strings.Split(addrs[0].String(), "/")[0]
	}
	return addr
}

func init() {
	flag.StringVar(&netIface, "i", "any", "Listen interface")
	flag.StringVar(&dir, "d", ".", "Directory to expose")
	flag.IntVar(&localPort, "p", 9090, "Listen port")
	flag.BoolVar(&ssl, "ssl", false, "SSL support")
	flag.Parse()
}

func main() {
	localAddr := getAddrFromIfaceName(netIface)
	bindValue := fmt.Sprintf("%s:%d", localAddr, localPort)
	log.Printf("Listen to %s\n", bindValue)
	var err error

	http.Handle("/", http.FileServer(http.Dir(dir)))
	if ssl {
		err = http.ListenAndServeTLS(bindValue, "server.crt", "server.key", nil)
	} else {
		err = http.ListenAndServe(bindValue, nil)
	}
	log.Fatal(err)
}
