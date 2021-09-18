package main

import (
	"GW/cfg"
	"GW/inBoundHandle"
	"GW/log"
	"GW/net"
	"GW/outBoundHandle"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)
var (
	westIface string
	bpfOut string
	logLife time.Duration
	logSpliteTime time.Duration
	logFilePath string
	logger *logrus.Logger
	stdOut bool
	isDebug bool

	wgMain sync.WaitGroup
	configs cfg.Configs
)

func initConfigs(){
	//初始化配置
	westIface = configs.GWCfgs.LocalNic
	logLife = configs.GWCfgs.LogLife
	logSpliteTime = configs.GWCfgs.LogSpliteTime
	logFilePath = configs.GWCfgs.LogFilePath
	stdOut = configs.GWCfgs.STDOUT
	isDebug = configs.GWCfgs.IsDebug
	bpfOut = configs.NetCfgs.BPFOut
	net.EthInit(configs.GWCfgs.LocalNic)
}
func initLoggers(){
	//初始化loggers
	logger = log.SetLogger(logFilePath + "/GW", logLife, logSpliteTime, stdOut, isDebug)
	logger.Info("----------------------------------------------")
	logger.Info("[GW]: GW is starting")
	logger.Info("[GW]: configs has initialized successfully")
	logger.Info("[GW]: logger has initialized successfully")
}

func main() {
	//初始化配置文件，并保存在Configs变量中
	configs = cfg.InitConfigs()
	initConfigs()
	initLoggers()
	wgMain.Add(1)
	if configs.GWCfgs.PktCatch == "AF_Packet"{
		logger.Info("[Agent]: grab packet by AF_Packet", configs.GWCfgs.AF_PktVersion)
		go outBoundHandle.OutBoundAP(westIface, bpfOut, &configs, logger)
	}else if configs.GWCfgs.PktCatch == "Pcap"{
		logger.Info("[Agent]: grab packet by Pcap")
		go outBoundHandle.OutBoundPcap(westIface, bpfOut, &configs, logger)
	}else{
		logger.Fatal("[Agent]: the model of packet catch is unrecognized")
	}
	wgMain.Add(1)
	go inBoundHandle.InBound(&configs, logger)
	wgMain.Wait()
}
