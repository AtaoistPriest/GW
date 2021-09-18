package net

import (
	"github.com/google/gopacket/afpacket"
	"golang.org/x/sys/unix"
	"net"
	"time"
)

type AfpacketHandle struct {
	TPacket *afpacket.TPacket
}

//新建一个af_packet的socket，用于从网卡读取数据
func NewAfpacketHandle(device string, frame_size int, block_size int, num_blocks int,
	useVLAN bool, timeout time.Duration, afPktVersion int) (*AfpacketHandle, error) {

	h := &AfpacketHandle{}
	var err error

	var version = []afpacket.OptTPacketVersion{-1, afpacket.TPacketVersion1, afpacket.TPacketVersion2, afpacket.TPacketVersion3}

	if device == "any" {
		h.TPacket, err = afpacket.NewTPacket(
			afpacket.OptFrameSize(frame_size),
			afpacket.OptBlockSize(block_size),
			afpacket.OptNumBlocks(num_blocks),
			afpacket.OptAddVLANHeader(useVLAN),
			afpacket.OptPollTimeout(timeout),
			afpacket.SocketRaw,
			version[afPktVersion])
	} else {
		h.TPacket, err = afpacket.NewTPacket(
			afpacket.OptInterface(device),
			afpacket.OptFrameSize(frame_size),
			afpacket.OptBlockSize(block_size),
			afpacket.OptNumBlocks(num_blocks),
			afpacket.OptAddVLANHeader(useVLAN),
			afpacket.OptPollTimeout(timeout),
			afpacket.SocketRaw,
			version[afPktVersion])
	}
	return h, err
}

//二层sockekt， 指定网卡建立链路层raw packet. eg : unix.Write(fd, rd)
func RawSocket(ifaceName string) (fd int, err error) {
	protol := (unix.ETH_P_ALL<<8)&0xff00 | unix.ETH_P_ALL>>8
	fd, err = unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(protol))
	if err != nil {
		return -1, err
	}
	ifIndex := 0
	if ifaceName != "" {
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			return -1, err
		}
		ifIndex = iface.Index
	}
	s := &unix.SockaddrLinklayer{
		Protocol: uint16(protol),
		Ifindex:  ifIndex,
	}
	err = unix.Bind(fd, s)
	return fd, err
}
//三层socket，创建网络层 raw socekt，组包时不需要指定MAC地址，走内核协议栈
func NetRawSocket()(fd int, err error){
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)
	//进程自定义ip头部
	err = unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_HDRINCL, 1)
	return fd, err
}

//建立四层 ubp socket
func SocketIpv4Udp(ipLocal [4]byte, portLocal int) (fd int, err error) {
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return -1, err
	}
	//绑定源ip与源port
	srcAddr := &unix.SockaddrInet4{
		Port: portLocal,
		Addr: ipLocal,
	}
	err = unix.Bind(fd, srcAddr)
	return fd, err
}

func SocketIpv4Tcp(srcIp [4]byte, srcPort int, dstIp [4]byte, dstPort int) (fd int, err error) {
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}
	//绑定源ip与源port
	srcAddr := &unix.SockaddrInet4{
		Port: srcPort,
		Addr: srcIp,
	}
	dstAddr := &unix.SockaddrInet4{
		Port: dstPort,
		Addr: dstIp,
	}
	err = unix.Bind(fd, srcAddr)
	err = unix.Connect(fd, dstAddr)
	return fd, err
}

func SendTo(fd int, remoteIp [4]byte, remotePort int, data []byte) error {
	toAddress := &unix.SockaddrInet4{
		Port: remotePort,
		Addr: remoteIp,
	}
	err := unix.Sendto(fd, data, 0, toAddress)
	return err
}

func RecvFrom(fd int, recvBuff []byte, offset int) (int, *unix.SockaddrInet4, error){
	dataLen, from, err := unix.Recvfrom(fd, recvBuff[offset:], 0)
	fromAddr, _ := from.(*unix.SockaddrInet4)
	return dataLen, fromAddr, err
}
