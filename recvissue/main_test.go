package main

import (
	"fmt"
	"net"
	"testing"

	"golang.org/x/net/ipv4"
)

func TestRecv(t *testing.T) {
	const lport = 12345

	c1, err := net.ListenUDP("udp4", &net.UDPAddr{Port: lport})
	if err != nil {
		t.Fatal(err)
	}
	pc := newPacketConn4(c1, lport)
	t.Log("Listening on", pc.LocalAddr())

	c2, err := net.Dial("udp4", fmt.Sprintf("127.0.0.1:%d", lport))
	if err != nil {
		t.Fatal(err)
	}

	c2.Write([]byte("hello"))

	buf := make([]byte, 100)
	_, src, dst, err := pc.ReadFromEx(buf)
	if err != nil {
		t.Fatal(err)
	}

	if src == nil {
		t.Error("src should not be nil")
	}
	if dst == nil {
		t.Error("dst should not be nil")
	}
}

// packetConn4 is a concrete packetConn for IPv4
type packetConn4 struct {
	net.PacketConn
	pc    *ipv4.PacketConn
	lport int
}

func newPacketConn4(conn net.PacketConn, lport int) *packetConn4 {
	pc := ipv4.NewPacketConn(conn)
	pc.SetControlMessage(ipv4.FlagDst|ipv4.FlagSrc, true)
	return &packetConn4{
		PacketConn: conn,
		pc:         pc,
		lport:      lport,
	}
}

func (c *packetConn4) ReadFromEx(buf []byte) (n int, src, dst net.Addr, err error) {
	n, cm, src, err := c.pc.ReadFrom(buf)
	if cm != nil {
		dst = &net.UDPAddr{IP: cm.Dst, Port: c.lport}
	}
	return n, src, dst, err
}
