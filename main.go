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

func getIPFromIfaceName(iface string) net.IP {
	var result net.IP

	if iface == "any" {
		result = net.ParseIP("0.0.0.0")
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
		// Retrieving the first IPv4 or IPv6 address
		for k := range addrs {
			s := strings.Split(addrs[k].String(), "/")[0]
			ip := net.ParseIP(s)
			if ip.To4() != nil { // IPv4 address
				result = ip
				break
			}
		}
		if result == nil {
			log.Printf("Error retrieving a valid IPv4 address for %s interface\n", iface)
			os.Exit(3)
		}
	}
	return result
}

func main() {
	netIface := flag.String("i", "any", "Listen interface")
	dir := flag.String("d", ".", "Directory to expose")
	localPort := flag.Int("p", 9090, "Listen port")
	ssl := flag.Bool("ssl", false, "SSL support")
	flag.Parse()

	localIP := getIPFromIfaceName(*netIface)
	bindValue := fmt.Sprintf("%s:%d", localIP.String(), *localPort)
	log.Printf("Listen to %s\n", bindValue)
	var err error

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	if *ssl {
		err = http.ListenAndServeTLS(bindValue, "server.crt", "server.key", nil)
	} else {
		err = http.ListenAndServe(bindValue, nil)
	}
	log.Fatal(err)
}
