package main

// 　ネットワークインターフェース層

// import (
// 	"encoding/hex"
// 	"fmt"

// 	"github.com/kawa1214/tcp-ip-go/network"
// )

// func main() {
// 	network, _ := network.NewTun()
// 	network.Bind()

// 	for {
// 		pkt, _ := network.Read()
// 		fmt.Print(hex.Dump(pkt.Buf[:pkt.N]))
// 		network.Write(pkt)
// 	}
// }

//  インターネット層

import (
	"fmt"

	"github.com/kawa1214/tcp-ip-go/internet"
	"github.com/kawa1214/tcp-ip-go/network"
)

func main() {
	network, _ := network.NewTun()
	network.Bind()
	ip := internet.NewIpPacketQueue()
	ip.ManageQueues(network)

	for {
		pkt, _ := ip.Read()
		fmt.Printf("IP Header: %+v\n", pkt.IpHeader)
	}
}
