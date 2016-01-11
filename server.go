// udpserver runs a simple UDP server. It is invoked with the PacketHandler
// interface which will Receive packets as byte slices along with the address
// they were received from.

package udpserver

import (
	"github.com/adamcolton/err"
	"net"
	"time"
)

const MaxUdpPacketLength = 65507

type UDPServer struct {
	conn          *net.UDPConn
	packetHandler PacketHandler
	localIP       string
	stop, running bool
}

type PacketHandler interface {
	Receive([]byte, *net.UDPAddr)
}

type UDPAddr net.UDPAddr

// New creates a UDPserver
// passing in ":0" for port will select any open port
func New(port string, packetHandler PacketHandler) (*UDPServer, error) {
	laddr, e := net.ResolveUDPAddr("udp", port)
	if !err.Check(e) {
		return nil, e
	}
	conn, e := net.ListenUDP("udp", laddr)
	if !err.Check(e) {
		return nil, e
	}

	server := &UDPServer{
		conn:          conn,
		packetHandler: packetHandler,
		localIP:       getLocalIP(),
		stop:          false,
		running:       false,
	}
	return server, nil
}

func RunNew(port string, packetHandler PacketHandler) (*UDPServer, error) {
	s, e := New(port, packetHandler)
	if err.Check(e) {
		go s.Run()
	}
	return s, e
}

// run is the servers listen loop
func (s *UDPServer) Run() {
	if s.running {
		return
	}
	s.running = true
	buf := make([]byte, MaxUdpPacketLength)
	for {
		l, addr, e := s.conn.ReadFromUDP(buf)
		if s.stop {
			break
		}
		if err.Log(e) {
			packet := make([]byte, l)
			copy(packet, buf[:l])
			go s.packetHandler.Receive(packet, addr)
		}
	}
	s.running = false
}

// Stop will stop the server
func (s *UDPServer) Stop() error {
	s.stop = true
	return s.conn.SetReadDeadline(time.Now()) // kill all reads
}

func (s *UDPServer) Close() error {
	e := s.Stop()
	if !err.Check(e) {
		return e
	}
	conn := s.conn
	s.conn = nil
	return conn.Close()
}

// Send will send a single packe (byte slice) to an address
// just a wrapper around WriteToUDP
func (s *UDPServer) Send(packet []byte, addr *net.UDPAddr) (int, error) {
	return s.conn.WriteToUDP(packet, addr)
}

// SendAll sends a slice of packets (byte slices) to an address
// this will return the last error it encoutered, if it encountered any
func (s *UDPServer) SendAll(packets [][]byte, addr *net.UDPAddr) error {
	var last error
	for _, p := range packets {
		if _, e := s.Send(p, addr); e != nil {
			last = e
		}
		time.Sleep(time.Millisecond)
	}
	return last
}

// LocalIP is a getter for the localIP, which is set when the server starts
func (s *UDPServer) LocalIP() string { return s.localIP }

// getLocalIp loops over the interface addresses to find one that is not a loopback
// address and uses that as it's local IP. It may not be fool proof and requires
// further investigation, but it does seem to work.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		var ip *net.IP
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ip = &ipnet.IP
		} else if ipaddr, ok := a.(*net.IPAddr); ok && !ipaddr.IP.IsLoopback() {
			ip = &ipaddr.IP
		}
		if ip != nil {
			if ip.To4() != nil {
				addr := ip.String()
				if addr != "0.0.0.0" {
					return addr
				}
			}
		}
	}
	return ""
}
