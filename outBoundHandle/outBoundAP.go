package outBoundHandle

import (
	"GW/cfg"
	"GW/net"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
)

var (
	frameSize int
	blockSize int
	blockNr int
	localIp [4]byte
	sendPort int
	logger *logrus.Logger
	received uint64
)

func intConfigs(configs *cfg.Configs){
	frameSize = configs.GWCfgs.FrameSize
	blockNr = configs.GWCfgs.BlockNr
	blockSize = configs.GWCfgs.BlockSize
	localIp = configs.NetCfgs.LocalIp
	sendPort = configs.NetCfgs.SendPort
	received = 0
}

//收包、流量统计、信息管理（sa，session）,并将任务分发给worker线程

func OutBoundAP(ifaceIn string, filter string, configs *cfg.Configs, loggerTmp *logrus.Logger) {
	intConfigs(configs)
	logger = loggerTmp
	logger.Info("[GW][OutBound]: OutBound is building Now base on AF_Packet")
	logger.Debug("[GW][OutBound]: west-->east; west iface is: ", ifaceIn, " frameSize : ", frameSize, " blockSize " +
		": ", blockSize, " blockNr : ", blockNr)
	//创建AF_PACKET的socket，用于从网卡读取数据，同时创建内核与用户共享的内存
	afpacketHandle, err := net.NewAfpacketHandle(ifaceIn, frameSize, blockSize, blockNr, false, pcap.BlockForever, configs.GWCfgs.AF_PktVersion)
	if err != nil {
		logger.Error("[GW][OutBound]: GW OutBound can't build Af_packet.Detail Info is : ", err)
		return
	}
	//建立用于发送的udp socket
	fd, err := net.SocketIpv4Udp(localIp, sendPort)
	if err != nil {
		logger.Error("[GW][OutBound]: GW OutBound can't build UDP Socket.Detail Info is : ", err, " address : ", localIp, ":", sendPort)
		return
	}
	//关闭混杂模式
	logger.Info("[GW][OutBound]: close the net promiscuous")
	err = afpacketHandle.TPacket.SetPromiscuous(ifaceIn, false)
	if err != nil {
		logger.Error("[GW][OutBound]: SetPromiscuous ERROR.Detail info : ", err)
		return
	}
	//设置编译BPF（减少杂包）
	logger.Debug("[GW][OutBound]: set BPF filter : ", filter)
	err = afpacketHandle.SetBPFFilter(filter, frameSize)
	if err != nil {
		logger.Error("[GW][OutBound]: SetBPFFilter ERROR.Detail info : ", err)
		return
	}
	logger.Info("[GW][OutBound]: OutBound is initialization Now")
	buff := make([]byte, frameSize)
	for {
		dataLen, _, err := afpacketHandle.TPacket.ReadPacketData(0, buff)
		if err != nil {
			logger.Error("[GW][OutBound]: when af_apcket read packet,there is an error happening.Detail info : ", err)
		}
		//过滤大于MTU的包
		if dataLen - 14 > configs.NetCfgs.MTU{
			logger.Warn("[GW][OutBound]: the len of pkt caught is longer than MTU. len of pkt is : ", dataLen - 14)
			continue
		}
		err = net.SendTo(fd, configs.NetCfgs.RemoteIp, configs.NetCfgs.RecvPort, buff[:dataLen])
		if err != nil{
			logger.Error("[ME][OutBound]: when outbound send a packets by UDP, there is an error happening.Detail info is : ", err)
			continue
		}
	}
}
