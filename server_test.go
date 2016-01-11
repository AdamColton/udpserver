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
	s, e := New(":5555", p)
	err.Test(e, t)
	s2, e := RunNew(":5556", p)
	err.Test(e, t)

	ip, e := net.ResolveUDPAddr("udp", ":5556")
	err.Test(e, t)
	s.Send([]byte{1, 2, 3}, ip)

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
		if len(p.packet) > 0 {
			break
		}
	}

	if len(p.packet) != 3 || p.packet[0] != 1 || p.packet[1] != 2 || p.packet[2] != 3 {
		t.Error("Incorrect Packet")
	}

	addr := p.addr.String()
	l := len(addr)
	if l < 4 || addr[l-4:] != "5555" {
		t.Error("Incorrect Address")
	}

	s.Close()
	s2.Close()
}

func TestStop(t *testing.T) {
	p := &PH{}
	s, e := RunNew(":5557", p)
	err.Test(e, t)
	time.Sleep(time.Millisecond)
	if !s.running {
		t.Error("Server is not running")
	}
	s.Stop()
	time.Sleep(time.Millisecond)
	if s.running {
		t.Error("Server has not stopped")
	}
	s.Close()
}

func TestClose(t *testing.T) {
	p := &PH{}
	s, e := RunNew(":5558", p)
	err.Test(e, t)
	time.Sleep(time.Millisecond)
	s.Close()
	time.Sleep(time.Millisecond)
	s2, e := RunNew(":5558", p)
	err.Test(e, t) // if 5558 is still in use, this will fail
	time.Sleep(time.Millisecond)
	s2.Close()
}
