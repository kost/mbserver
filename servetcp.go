package mbserver

import (
	"io"
	"log"
	"net"
	"strings"
)

func (s *Server) ServeConn(conn net.Conn) (err error) {
	defer conn.Close()

	for {
		packet := make([]byte, 512)
		bytesRead, err := conn.Read(packet)
		if err != nil {
			return err
		}
		// Set the length of the packet to the number of read bytes.
		packet = packet[:bytesRead]

		frame, err := NewTCPFrame(packet)
		if err != nil {
			return err
		}

		request := &Request{conn, frame}

		s.requestChan <- request
	}
	return nil
}

func (s *Server) accept(listen net.Listener) error {
	for {
		conn, err := listen.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			log.Printf("Unable to accept connections: %#v\n", err)
			return err
		}

		go s.ServeConn(conn)
	}
}

// ListenTCP starts the Modbus server listening on "address:port".
func (s *Server) ListenTCP(addressPort string) (err error) {
	listen, err := net.Listen("tcp", addressPort)
	if err != nil {
		log.Printf("Failed to Listen: %v\n", err)
		return err
	}
	s.listeners = append(s.listeners, listen)
	go s.accept(listen)
	return err
}
