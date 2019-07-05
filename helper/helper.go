package helper

import (
	"net"
	"os"
	"path/filepath"
	"sea_log/logs"
	"strings"
)

func LocalIPv4s() (string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips[0], nil
}

func GetRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.ERROR(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}
