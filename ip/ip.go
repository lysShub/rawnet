package ip

import (
	"net"
	"strings"
)

var err error

// GetLocalIP 获取内网IP
func GetLocalIP() net.IP {
	var conn net.Conn
	if conn, err = net.Dial("udp", "114.114.114.114:80"); err != nil {
		return nil
	}
	defer conn.Close()
	return net.ParseIP(strings.Split(conn.LocalAddr().String(), ":")[0])
}
