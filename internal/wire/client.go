package wire

import (
	"errors"
	"net"
	"net/netip"
	"os"
	"sync"
	"time"
)

var (
	// ErrMaxConns is returned when dns client reaches the max connections limitation.
	ErrMaxConns = errors.New("dns client reaches the max connections limitation")
)

// Client is an UDP client that supports DNS protocol.
type Client struct {
	AddrPort netip.AddrPort

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections. Zero means no limit.
	MaxIdleConns int

	// MaxConns optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, ErrMaxConns will be return.
	//
	// Zero means no limit.
	MaxConns int

	// ReadTimeout is the maximum duration for reading the dns server response.
	ReadTimeout time.Duration

	mu    sync.Mutex
	conns []*net.UDPConn
}

// Exchange executes a single DNS transaction, returning
// a Response for the provided Request.
func (c *Client) Exchange(req, resp *Message) (err error) {
	err = c.exchange(req, resp)
	if err != nil && os.IsTimeout(err) {
		err = c.exchange(req, resp)
	}
	return err
}

func (c *Client) exchange(req, resp *Message) error {
	var fresh bool
	conn, err := c.get()
	if conn == nil && err == nil {
		conn, err = c.dial()
		fresh = true
	}
	if err != nil {
		return err
	}

	_, err = conn.Write(req.Raw)
	if err != nil && !fresh {
		// if error is a pooled conn, let's close it & retry again
		conn.Close()
		if conn, err = c.dial(); err != nil {
			return err
		}
		if _, err = conn.Write(req.Raw); err != nil {
			return err
		}
	}

	if c.ReadTimeout > 0 {
		err = conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
		if err != nil {
			return err
		}
	}

	resp.Raw = resp.Raw[:cap(resp.Raw)]
	n, err := conn.Read(resp.Raw)
	if err == nil {
		resp.Raw = resp.Raw[:n]
		err = ParseMessage(resp, resp.Raw, false)
	}

	c.put(conn)

	return err
}

func (c *Client) dial() (conn *net.UDPConn, err error) {
	conn, err = net.DialUDP("udp", nil, net.UDPAddrFromAddrPort(c.AddrPort))
	return
}

func (c *Client) get() (conn *net.UDPConn, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := len(c.conns)
	if c.MaxConns != 0 && count > c.MaxConns {
		err = ErrMaxConns

		return
	}
	if count > 0 {
		conn = c.conns[len(c.conns)-1]
		c.conns = c.conns[:len(c.conns)-1]
	}

	return
}

func (c *Client) put(conn *net.UDPConn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if (c.MaxIdleConns != 0 && len(c.conns) > c.MaxIdleConns) ||
		(c.MaxConns != 0 && len(c.conns) > c.MaxConns) {
		conn.Close()

		return
	}

	c.conns = append(c.conns, conn)
}
