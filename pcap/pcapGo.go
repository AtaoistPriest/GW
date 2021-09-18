package pcap

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"time"
)

/*
*	iface  		网卡名称
*	snaplen		单个数据包的最大长度
*	timeOut		handle读取数据包的超时时间
*	promisc		是否启用混杂模式
*	bpfStr		bpf过滤规则
*/

func BuildHandle(iface string, snaplen int, timeOut int, promisc bool, bpfStr string) (*gopacket.PacketSource, error){
	var handle *pcap.Handle
	var err error
	inactive, err := pcap.NewInactiveHandle(iface)
	if err != nil {
		return nil, err
	}
	defer inactive.CleanUp()
	if err = inactive.SetSnapLen(snaplen); err != nil {
		return nil, err
	}
	if err = inactive.SetPromisc(promisc); err != nil {
		return nil, err
	}
	//SetTimeout sets the read timeout for the handle.
	if err = inactive.SetTimeout(time.Millisecond * time.Duration(timeOut)); err != nil {
		return nil, err
	}
	if handle, err = inactive.Activate(); err != nil {
		return nil, err
	}
	if err = handle.SetBPFFilter(bpfStr); err != nil {
		return nil, err
	}
	source := gopacket.NewPacketSource(handle, nil)
	source.Lazy = false
	source.NoCopy = true
	source.DecodeStreamsAsDatagrams = false
	source.SkipDecodeRecovery = false
	//自定义修改:此处已修改gopacket的源码，已获取source的赋值
	source.ResetSource()
	return source, nil
}

func ReadNextPkt(source *gopacket.PacketSource) ([]byte, error) {
	data, _, err := source.SourceTmp.ReadPacketData()
	return data, err
}
