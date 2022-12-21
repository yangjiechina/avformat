package utils

import (
	"net"
)

const PortMaximum = 65535

func Used(port int, tcp bool) bool {
	if tcp {
		listenTCP, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
		if err == nil {
			listenTCP.Close()
		}
		return err != nil
	} else {
		udp, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
		if err == nil {
			udp.Close()
		}
		return err != nil
	}
}

func AllocPort(startPort int, tcp bool) (int, bool) {
	for PortMaximum >= startPort {
		if !Used(startPort, tcp) {
			return startPort, true
		}

		startPort++
	}
	return 0, false
}

func AllocPairPort(startPort int, tcp bool) (int, int, bool) {
	for PortMaximum-startPort >= 2 {
		if !Used(startPort, tcp) && !Used(startPort+1, tcp) {
			return startPort, startPort + 1, true
		}
		startPort += 2
	}
	return 0, 0, false
}
