package utils

import (
	"net"
	"os"
)

func GetLocalIp() []string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	var result []string

	for _, a := range address {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				result = append(result, ipNet.IP.String())
			}
		}
	}

	result = append(result, "127.0.0.1")
	return result
}
