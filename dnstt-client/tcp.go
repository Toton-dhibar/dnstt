package main

import (
    "bufio"
    "context"
    "encoding/binary"
    "io"
    "log"
    "net"
    "sync"
    "time"

    "www.bamsoftware.com/git/dnstt.git/turbotunnel"
)

const tcpDialTimeout = 30 * time.Second

// TCPPacketConn is a TCP-based transport for DNS messages. It maintains a TCP
// connection to the resolver, reconnecting as necessary. It closes the connection
// if any reconnection attempt fails.
//
// TCPPacketConn deals only with already formatted DNS messages. It does not
// handle encoding information into the messages. That is rather the
// responsibility of DNSPacketConn.
type TCPPacketConn struct {
    // QueuePacketConn is the direct receiver of ReadFrom and WriteTo calls.
    // recvLoop and sendLoop take the messages out of the receive and send
    // queues and actually put them on the network.
    *turbotunnel.QueuePacketConn
}

// NewTCPPacketConn creates a new TCPPacketConn configured to use the TCP
// server at addr as a DNS over TCP resolver. It maintains a TCP connection to
// the resolver, reconnecting as necessary. It closes the connection if any
// reconnection attempt fails.
func NewTCPPacketConn(addr string) (*TCPPacketConn, error) {
	dial := func() (net.Conn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), tcpDialTimeout)
		defer cancel()
		return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	}
	// We maintain one TCP connection at a time, redialing it whenever it
	// becomes disconnected. We do the first dial here, outside the
	// goroutine, so that any immediate and permanent connection errors are
	// reported directly to the caller of NewTCPPacketConn.
	conn, err := dial()
	if err != nil {
		return nil, err
	}
	c := &TCPPacketConn{
		QueuePacketConn: turbotunnel.NewQueuePacketConn(turbotunnel.DummyAddr{}, 0),
	}
	go func() {
		defer c.Close()
		for {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				err := c.recvLoop(conn)
				if err != nil {
					log.Printf("recvLoop: %v", err)
				}
				wg.Done()
			}()
			go func() {
				err := c.sendLoop(conn)
				if err != nil {
					log.Printf("sendLoop: %v", err)
				}
				wg.Done()
			}()
			wg.Wait()
			conn.Close()

			// Whenever the TCP connection dies, redial a new one.
			// Keep trying until we succeed.
			for {
				conn, err = dial()
				if err != nil {
					log.Printf("dial tcp: %v", err)
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}
				break // Successfully reconnected
			}
		}
	}()
	return c, nil
}

// recvLoop repeatedly reads DNS messages from conn and queues them in the
// PacketConn.
func (c *TCPPacketConn) recvLoop(conn net.Conn) error {
	br := bufio.NewReader(conn)
	for {
		var length uint16
		err := binary.Read(br, binary.BigEndian, &length)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		p := make([]byte, int(length))
		_, err = io.ReadFull(br, p)
		if err != nil {
			return err
		}
		c.QueuePacketConn.QueueIncoming(p, turbotunnel.DummyAddr{})
	}
}

// sendLoop repeatedly reads DNS messages from the PacketConn and writes them to conn.
func (c *TCPPacketConn) sendLoop(conn net.Conn) error {
	bw := bufio.NewWriter(conn)
	for p := range c.QueuePacketConn.OutgoingQueue(turbotunnel.DummyAddr{}) {
		length := uint16(len(p))
		if int(length) != len(p) {
			panic(len(p))
		}
		err := binary.Write(bw, binary.BigEndian, &length)
		if err != nil {
			return err
		}
		_, err = bw.Write(p)
		if err != nil {
			return err
		}
		err = bw.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}