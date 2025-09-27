package dnsttclient

import (
    "io"
    "net"
    "time"
)

// TCPPacketConn represents a TCP connection for DNS over TCP.
type TCPPacketConn struct {
    conn net.Conn
}

// NewTCPPacketConn establishes a new TCP connection.
func NewTCPPacketConn(conn net.Conn) *TCPPacketConn {
    return &TCPPacketConn{conn: conn}
}

// ReadFrom reads a DNS message from the connection.
func (c *TCPPacketConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
    // Read the 2-byte length prefix
    var length [2]byte
    if _, err := io.ReadFull(c.conn, length[:]); err != nil {
        return 0, nil, err
    }
    msgLen := int(length[0])<<8 | int(length[1])
    
    // Read the DNS message
    return io.ReadFull(c.conn, b[:msgLen])
}

// WriteTo writes a DNS message to the connection.
func (c *TCPPacketConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
    // Write the 2-byte length prefix
    length := []byte{byte(len(b) >> 8), byte(len(b))}
    if _, err := c.conn.Write(length); err != nil {
        return 0, err
    }
    
    // Write the DNS message
    return c.conn.Write(b)
}

// Close closes the connection.
func (c *TCPPacketConn) Close() error {
    return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *TCPPacketConn) LocalAddr() net.Addr {
    return c.conn.LocalAddr()
}

// SetDeadline sets the read and write deadlines.
func (c *TCPPacketConn) SetDeadline(t time.Time) error {
    return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
func (c *TCPPacketConn) SetReadDeadline(t time.Time) error {
    return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
func (c *TCPPacketConn) SetWriteDeadline(t time.Time) error {
    return c.conn.SetWriteDeadline(t)
}