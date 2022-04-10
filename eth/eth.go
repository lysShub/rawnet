package eth

import (
	"net"
)

func SendRawEthFrame(src, dst net.HardwareAddr, payload []byte) (err error) {
	return sendRawEthFrame(src, dst, payload)
}

// GetSrcMAC 获取默认网卡MAC地址和下一跳MAC地址
// localIP指定本机网卡, nil表示默认网卡
func GetMAC(localIP net.IP) (srcMAC, dstMAC net.HardwareAddr, err error) {
	return getMAC(localIP)
}

// GetIfiByIP 获取指定IP的网卡
func GetIfiByIP(localIP net.IP) (ifi *net.Interface, err error) {
	var ifis []net.Interface
	if ifis, err = net.Interfaces(); err != nil {
		return
	}

	for _, i := range ifis {
		var addrs []net.Addr
		if addrs, err = i.Addrs(); err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			if ip, _, err = net.ParseCIDR(addr.String()); err != nil {
				continue
			}

			if ip.Equal(localIP) {
				return &i, nil
			}
		}
	}

	return
}

func packEthFrameV4(src, dst net.HardwareAddr, payload []byte) (frame []byte) {
	if 14+len(payload) < 64 {
		frame = make([]byte, 64)
	} else {
		frame = make([]byte, 14+len(payload))
	}
	frame = append(frame, dst...)
	frame = append(frame, src...)
	frame = append(frame, payload...)
	if len(frame) < 64 {
		for len(frame) < 64 {
			frame = append(frame, 0)
		}
	}
	return frame
}
