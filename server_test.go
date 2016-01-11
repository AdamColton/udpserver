package udpserver

import (
	"github.com/adamcolton/err"
	"net"
	"testing"
	"time"
)

type PH struct {
	packet []byte
	addr   *net.UDPAddr
}

func (ph *PH) Receive(pck []byte, addr *net.UDPAddr) {
	ph.packet = pck
	ph.addr = addr
}

func TestServer(t *testing.T) {
	p := &PH{}
	a := New(":5555", p)
	New(":5556", p)

	ip, e := net.ResolveUDPAddr("udp", ":5556")
	err.Test(e, t)
	a.Send([]byte{1, 2, 3}, ip)

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
		if len(p.packet) > 0 {
			break
		}
	}

	if len(p.packet) != 3 || p.packet[0] != 1 || p.packet[1] != 2 || p.packet[2] != 3 {
		t.Error("Incorrect Packet")
	}

	s := p.addr.String()
	l := len(s)
	if s[l-4:] != "5555" {
		t.Error("Incorrect Address")
	}
}
