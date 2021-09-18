package net

import (
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/bpf"
)

// SetBPFFilter translates a BPF filter string into BPF RawInstruction and applies them.
func (h *AfpacketHandle) SetBPFFilter(filter string, frame_size int) (err error) {
	pcapBPF, err := pcap.CompileBPFFilter(layers.LinkTypeEthernet, frame_size, filter)
	if err != nil {
		return err
	}
	bpfIns := []bpf.RawInstruction{}
	for _, ins := range pcapBPF {
		bpfIns2 := bpf.RawInstruction{
			Op: ins.Code,
			Jt: ins.Jt,
			Jf: ins.Jf,
			K:  ins.K,
		}
		bpfIns = append(bpfIns, bpfIns2)
	}
	if h.TPacket.SetBPF(bpfIns); err != nil {
		return err
	}
	return h.TPacket.SetBPF(bpfIns)
}

