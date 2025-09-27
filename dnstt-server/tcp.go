package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

// TCPServer represents the TCP server.
type TCPServer struct {
	address string
}

// NewTCPServer creates a new instance of TCPServer.
func NewTCPServer(address string) *TCPServer {
	return &TCPServer{address: address}
}

// Start starts the TCP server.
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue // handle error appropriately
		}
		go s.handleConnection(conn)
	}
}

// handleConnection manages a single connection to the server.
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// Read the 2-byte length prefix
		lengthBytes := make([]byte, 2)
		_, err := io.ReadFull(conn, lengthBytes)
		if err != nil {
			break // handle error appropriately
		}

		// Convert lengthBytes to an integer
		messageLength := binary.BigEndian.Uint16(lengthBytes)

		// Read the DNS message
		message := make([]byte, messageLength)
		_, err = io.ReadFull(conn, message)
		if err != nil {
			break // handle error appropriately
		}

		// Here you would handle the DNS message (existing server logic)
		// For example: processDNSMessage(message)
	}
}

// main function to start the TCP server.
func main() {
	server := NewTCPServer(":53") // Change port as needed
	if err := server.Start(); err != nil {
		panic(err)
	}
}