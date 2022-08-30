package main

import (
	"flag"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
)

var iface string

func init() {
	flag.StringVar(&iface, "iface", "en0", "Network interface to listeon on")
}

func main() {
	log.Info("Starting Application")

	InitCrackling()

	if handle, err := pcap.OpenLive("en0", 1600, true, pcap.BlockForever); err != nil {
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
	if err := parser.DecodeLayers(p.Data(), &decoded); err != nil {
		log.Printf("processor %v", err)
	} else {
		fmt.Println(ip.DstIP)
		Crackle()
	}
}
