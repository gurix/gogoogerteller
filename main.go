package main

import (
	"flag"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/gurix/gogoogerteller/util"

	"github.com/google/gopacket/pcap"
)

var iface string
var localIP string

func init() {
	flag.StringVar(&iface, "iface", "en0", "Network interface to listeon on")
}

func main() {
	InitCrackling()

	if addr, err := util.GetInterfaceIpv4Addr(iface); err != nil {
		panic(err)
	} else {
		localIP = addr
	}

	if handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter("tcp"); err != nil { // Only watch for tcp traffic
		panic(err)
	} else {

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			handlePacket(packet)
		}
	}
}

func handlePacket(p gopacket.Packet) {
	var eth layers.Ethernet
	var ip layers.IPv4
	var ipv6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	var payload gopacket.Payload

	var decoded []gopacket.LayerType
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip, &ipv6, &tcp, &udp, &payload)
	if err := parser.DecodeLayers(p.Data(), &decoded); err == nil {
		if ip.SrcIP.String() == localIP && ip.DstIP != nil {
			fmt.Printf("%s --> %s\n", ip.SrcIP, ip.DstIP)
			Crackle()
		}
	}
}
