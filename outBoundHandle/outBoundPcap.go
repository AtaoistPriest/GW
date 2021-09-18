package outBoundHandle

import (
	"GW/cfg"
	"GW/net"
	"GW/pcap"
	"github.com/sirupsen/logrus"
)

func intConfigs_(configs *cfg.Configs){
	frameSize = configs.GWCfgs.FrameSize
	blockNr = configs.GWCfgs.BlockNr
	blockSize = configs.GWCfgs.BlockSize
	localIp = configs.NetCfgs.LocalIp
	sendPort = configs.NetCfgs.SendPort
	received = 0
}

//收包、流量统计、信息管理（sa，session）,并将任务分发给worker线程

func OutBoundPcap(ifaceIn string, filter string, configs *cfg.Configs, loggerTmp *logrus.Logger) {
	intConfigs_(configs)
	logger = loggerTmp
	logger.Info("[Agent][OutBound]: OutBound is building Now base on Pcap")
	logger.Debug("[Agent][OutBound]: catch pkts on interface : ", ifaceIn, " frameSize : ", frameSize)
	//创建Pcap的socket，用于从网卡读取数据，同时创建内核与用户共享的内存
	pcapHandle, err := pcap.BuildHandle(ifaceIn, configs.GWCfgs.FrameSize, configs.GWCfgs.TimeOut, configs.GWCfgs.Promisc, filter)
	if err != nil {
		logger.Error("[Agent][OutBound]: Agent OutBound can't build Pcap.Detail Info is : ", err)
		return
	}
	//建立用于发送的udp socket
	fd, err := net.SocketIpv4Udp(localIp, sendPort)
	if err != nil {
		logger.Error("[Agent][OutBound]: Agent OutBound can't build UDP Socket.Detail Info is : ", err, " address : ", localIp, ":", sendPort)
		return
	}
	logger.Info("[Agent][OutBound]: OutBound is initialization Now")
	buff := make([]byte, frameSize)
	for {
		//使用pcap读取下一个数据包
		dataTmp, err := pcap.ReadNextPkt(pcapHandle)
		if dataTmp == nil || err != nil{
			continue
		}
		copy(buff[0:], dataTmp)
		dataLen := len(dataTmp)
		//过滤大于MTU的包
		if dataLen - 14 > configs.NetCfgs.MTU{
			logger.Warn("[Agent][OutBound]: the len of pkt caught is longer than MTU. len of pkt is : ", dataLen - 14)
			continue
		}
		err = net.SendTo(fd, configs.NetCfgs.RemoteIp, configs.NetCfgs.RecvPort, buff[:dataLen])
		if err != nil{
			logger.Error("[ME][OutBound]: when outbound send a packets by UDP, there is an error happening.Detail info is : ", err)
			continue
		}
	}
}
