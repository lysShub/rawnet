//go:build windows
// +build windows

package eth

import (
	"errors"
	"net"
	"rawnet/ip"
	"syscall"
	"unsafe"

	"github.com/google/gopacket/routing"
	route "github.com/libp2p/go-netroute"
)

func sendRawEthFrame(src, dst net.HardwareAddr, payload []byte) (err error) {
	return
}

func getMAC(localIP net.IP) (srcMAC, dstMAC net.HardwareAddr, err error) {
	if localIP == nil {
		localIP = ip.GetLocalIP()
	}
	var r routing.Router
	if r, err = route.New(); err != nil {
		return
	}

	var ifi *net.Interface
	var getway net.IP
	if ifi, getway, _, err = r.RouteWithSrc(nil, localIP, net.IPv4zero); err != nil {
		return
	} else if ifi == nil {
		if ifi, err = GetIfiByIP(localIP); err != nil {
			return
		} else if ifi == nil {
			return nil, nil, errors.New("no interface found")
		}
	}

	srcMAC = ifi.HardwareAddr
	dstMAC, err = sendARP(getway)

	return
}

var sendARPFn = syscall.MustLoadDLL("iphlpapi.dll").MustFindProc("SendARP")

func sendARP(ip net.IP) (net.HardwareAddr, error) {
	ip = ip.To16()
	dst := (uint32(ip[12])) | (uint32(ip[13]) << 8) | (uint32(ip[14]) << 16) | (uint32(ip[15]) << 24)

	mac := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	n := uint32(len(mac))
	ret, _, _ := sendARPFn.Call(
		uintptr(dst),
		0,
		uintptr(unsafe.Pointer(&mac[0])),
		uintptr(unsafe.Pointer(&n)),
	)
	if ret != 0 {
		return nil, syscall.Errno(ret)
	}
	return mac, nil
}
