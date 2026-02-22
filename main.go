//go:build linux
// +build linux

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"kerbeus/src/banner"
	"kerbeus/src/capture"
)

func init() {
	flag.Usage = banner.ShowHelp
}

func main() {
	iface := flag.String("i", "", "network interface to use (e.g. eth0)")
	flag.Parse()

	dev, err := getInterface(*iface)
	if err != nil {
		log.Fatalf("\r\033[KError: %v", err)
	}

	banner.Print(dev)

	if err := capture.Start(dev); err != nil {
		log.Fatalf("\r\033[KError starting capture: %v", err)
	}
}

func getInterface(iface string) (string, error) {
	if iface != "" {
		return iface, nil
	}

	ifs, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("\r\033[Kcould not list interfaces: %v", err)
	}

	for _, ifi := range ifs {
		if (ifi.Flags&net.FlagUp) != 0 && (ifi.Flags&net.FlagLoopback) == 0 {
			return ifi.Name, nil
		}
	}

	return "", fmt.Errorf("\r\033[Kno valid interface found")
}
