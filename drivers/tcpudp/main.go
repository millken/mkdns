package tcpudp

import (
	"context"
	"net"
	"sync"

	"github.com/yourproject/driver"
)

type TCPUDPDriver struct {
	addr        string
	udpConn     *net.UDPConn
	tcpListener *net.TCPListener
	requests    chan driver.DNSRequest
	wg          sync.WaitGroup
}

func NewTCPUDPDriver(addr string) *TCPUDPDriver {
	return &TCPUDPDriver{
		addr:     addr,
		requests: make(chan driver.DNSRequest, 100),
	}
}

func (d *TCPUDPDriver) Start(ctx context.Context) error {
	udpAddr, err := net.ResolveUDPAddr("udp", d.addr)
	if err != nil {
		return err
	}
	d.udpConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", d.addr)
	if err != nil {
		return err
	}
	d.tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	d.wg.Add(2)
	go d.handleUDP(ctx)
	go d.handleTCP(ctx)

	return nil
}

func (d *TCPUDPDriver) Stop() error {
	d.udpConn.Close()
	d.tcpListener.Close()
	d.wg.Wait()
	close(d.requests)
	return nil
}

func (d *TCPUDPDriver) Receive() <-chan driver.DNSRequest {
	return d.requests
}

func (d *TCPUDPDriver) Send(response driver.DNSResponse) error {
	// 实现发送逻辑
	return nil
}

func (d *TCPUDPDriver) handleUDP(ctx context.Context) {
	defer d.wg.Done()
	buf := make([]byte, 1024)
	for {
		n, addr, err := d.udpConn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		d.requests <- driver.DNSRequest{
			RawData:  buf[:n],
			ClientIP: addr.IP,
		}
	}
}

func (d *TCPUDPDriver) handleTCP(ctx context.Context) {
	defer d.wg.Done()
	for {
		conn, err := d.tcpListener.AcceptTCP()
		if err != nil {
			return
		}
		go func(conn *net.TCPConn) {
			defer conn.Close()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			d.requests <- driver.DNSRequest{
				RawData:  buf[:n],
				ClientIP: conn.RemoteAddr().(*net.TCPAddr).IP,
			}
		}(conn)
	}
}
