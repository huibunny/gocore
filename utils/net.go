package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// GetHostIP get host ip address
func GetHostIP() string {
	hostAddress := ""
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				hostAddress = ipnet.IP.String()
				break
			}

		}
	}
	return hostAddress
}

func GetHostPort(addr string) (string, string) {
	host := "localhost"
	port := "0"
	hostPorts := strings.Split(addr, ":")
	if len(hostPorts) > 1 {
		host = hostPorts[0]
		port = hostPorts[1]
	} else {
		port = hostPorts[0]
	}
	if len(host) == 0 {
		host = GetHostIP()
	} else {
	}

	return host, port
}
