//go:build linux || netbsd || openbsd || darwin || freebsd
// +build linux netbsd openbsd darwin freebsd

package eth

import (
	"errors"
	"net"
	"time"

	"github.com/google/gopacket/routing"
	"github.com/mdlayher/arp"
	"golang.org/x/sys/unix"
)

// sendRawEthFrame 发送一个以太帧
func sendRawEthFrame(src, dst net.HardwareAddr, payload []byte) (err error) {
	var fd int
	fd, err = unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, unix.IPPROTO_RAW)
	// fd, err = unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, unix.ETH_P_ALL)
	// fd, err = unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, unix.ETH_P_IP)
	if err != nil {
		return err
	}

	var addr unix.SockaddrL2

	if copy(addr.Addr[:], dst) != 6 {
		return errors.New("invalid dst HardwareAddr")
	}

	err = unix.Sendto(fd, payload, 0, &addr)
	return
}

func getMAC(localIP net.IP) (srcMAC, dstMAC net.HardwareAddr, err error) {
	r, err := routing.New()
	if err != nil {
		panic(err)
	}

	var ifi *net.Interface
	var getway net.IP
	ifi, getway, _, err = r.RouteWithSrc(nil, localIP, net.IPv4zero)
	if err != nil {
		return
	} else if ifi == nil {
		if ifi, err = GetIfiByIP(localIP); err != nil {
			return
		} else if ifi == nil {
			return nil, nil, errors.New("no interface found")
		}
	}
	srcMAC = ifi.HardwareAddr

	var arpc *arp.Client
	if arpc, err = arp.Dial(ifi); err != nil {
		return
	} else {
		arpc.SetDeadline(time.Now().Add(time.Millisecond * 100))

		dstMAC, err = arpc.Resolve(getway)
		return
	}
}
