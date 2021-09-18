package inBoundHandle

import (
	"GW/cfg"
	"GW/net"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var (
	frameSize int
	recvPort int
	mtu int
	localIp [4]byte
	logger *logrus.Logger
)

func intConfigs(configs *cfg.Configs){
	frameSize = configs.GWCfgs.FrameSize
	localIp = configs.NetCfgs.LocalIp
	recvPort = configs.NetCfgs.RecvPort
	mtu = configs.NetCfgs.MTU

}

//东向网卡，接收对端ME数据后进行解密、解封装后通过西向网卡转发给目的端用户

func InBound(configs *cfg.Configs, loggerTmp *logrus.Logger) {
	intConfigs(configs)
	logger = loggerTmp
	logger.Info("[GW][InBound]: InBound is building Now")

	fd, err := net.NetRawSocket()
	if err != nil{
		logger.Error("[GW][InBound]: when the udp socket is building, an error is happening.Detail info : ", err)
	}
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_RCVBUF, configs.GWCfgs.RecvBuff)
	if err != nil{
		logger.Error("[GW][InBound]: when set socket recv buff, an error is happening.Detail info : ", err)
	}
	recvFd, err := net.SocketIpv4Udp(localIp, recvPort)
	if err != nil {
		logger.Error("[GW][InBound]: GW InBound can't build Udp socket with address : ", localIp, ":", recvPort ,". Detail Info is : ", err)
		return
	}

	logger.Info("[GW][InBound]: InBound is initialization Now")
	buff := make([]byte, frameSize)
	for {
		dataLen, _, err:= net.RecvFrom(recvFd, buff, 0)
		if err != nil {
			logger.Error("[GW][InBound]: when udp-socket read packet,there is an error happening.Detail info : ", err)
		}
		if dataLen > mtu{
			logger.Error("[GW][InBound]: the len of payload is longer than MTU")
			continue
		}
		dstIp := [4]byte{buff[16], buff[17], buff[18], buff[19]}
		err = net.SendTo(fd, dstIp,0, buff[:dataLen])
		if err != nil {
			logger.Error("[GW][InBound][Worker]: when InBound worker send a packets by UDP, there is an error happening.Detail info is : ", err)
		}
	}
}
