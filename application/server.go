package application

import (
	"fmt"

	"github.com/kawa1214/tcp-ip-go/datalink"
	"github.com/kawa1214/tcp-ip-go/network"
	"github.com/kawa1214/tcp-ip-go/transport"
)

type Server struct {
	device         *datalink.NetDevice
	ipPacketQueue  *network.IpPacketQueue
	tcpPacketQueue *transport.TcpPacketQueue
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ListenAndServe() error {
	device, err := datalink.NewTun()
	device.Bind()
	s.device = device
	if err != nil {
		return err
	}
	s.serve()
	return nil
}

func (s *Server) serve() {
	ipPacketQueue := network.NewIpPacketQueue()
	ipPacketQueue.ManageQueues(s.device)
	s.ipPacketQueue = ipPacketQueue

	tcpPacketQueue := transport.NewTcpPacketQueue()
	tcpPacketQueue.ManageQueues(ipPacketQueue)
	s.tcpPacketQueue = tcpPacketQueue
}

func (s *Server) Close() {
	s.device.Close()
	s.ipPacketQueue.Close()
	s.tcpPacketQueue.Close()
}

func (s *Server) Accept() (transport.Connection, error) {
	conn, err := s.tcpPacketQueue.ReadAcceptConnection()
	if err != nil {
		return transport.Connection{}, fmt.Errorf("accept error: %s", err)
	}

	return conn, nil
}

func (s *Server) Write(conn transport.Connection, resp *HTTPResponse) {
	pkt := conn.Pkt
	tcpDataLen := int(pkt.Packet.N) - (int(pkt.IpHeader.IHL) * 4) - (int(pkt.TcpHeader.DataOff) * 4)

	payload := resp.String()
	respNewIPHeader := network.NewIp(pkt.IpHeader.DstIP, pkt.IpHeader.SrcIP, transport.LENGTH+len(payload))
	respNewTcpHeader := transport.New(
		pkt.TcpHeader.DstPort,
		pkt.TcpHeader.SrcPort,
		pkt.TcpHeader.AckNum,
		pkt.TcpHeader.SeqNum+uint32(tcpDataLen),
		transport.HeaderFlags{
			PSH: true,
			ACK: true,
		},
	)
	sendPkt := transport.TcpPacket{
		IpHeader:  respNewIPHeader,
		TcpHeader: respNewTcpHeader,
	}

	s.tcpPacketQueue.Write(pkt, sendPkt, []byte(payload))
}
