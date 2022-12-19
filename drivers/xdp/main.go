package xdp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/millken/goenv"
	"github.com/millken/mkdns/drivers"
	"github.com/millken/mkdns/internal/ebpf"
	"github.com/millken/mkdns/internal/scheduler"
	"github.com/millken/mkdns/internal/xdp"
	"github.com/pkg/errors"
)

type xdpDriver struct {
	scheduler *scheduler.Scheduler
	linkName  string
	queueID   int
	txChan    chan []byte
	rxChan    chan []byte
}

func New(scheduler *scheduler.Scheduler) drivers.Driver {
	return &xdpDriver{
		scheduler: scheduler,
		linkName:  goenv.Get("XDP_LINKNAME", "enp3s0"),
		queueID:   goenv.Int("XDP_QUEUEID", 0),
		txChan:    make(chan []byte),
		rxChan:    make(chan []byte),
	}
}

func (d *xdpDriver) Name() string {
	return "xdp"
}

func (d *xdpDriver) Version() string {
	return "0.1"
}

func (d *xdpDriver) Start(ctx context.Context) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("error: failed to fetch the list of network interfaces on the system: %w", err)
	}
	ifIndex := -1
	for _, iface := range interfaces {
		if iface.Name == d.linkName {
			ifIndex = iface.Index
			break
		}
	}
	if ifIndex == -1 {
		return fmt.Errorf("error: couldn't find a suitable network interface to attach to")
	}
	program, err := ebpf.NewDNSProtoProgram(nil)
	if err != nil {
		return fmt.Errorf("error: failed to create xdp program: %w", err)
	}
	defer program.Close()
	if err := program.Attach(ifIndex); err != nil {
		return fmt.Errorf("error: failed to attach xdp program: %w", err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	go func() {
		<-sc
		err := program.Detach(ifIndex)
		if err != nil {
			panic(errors.Wrap(err, "detach failed"))
		}
		os.Exit(0)
	}()

	xsk, err := xdp.NewSocket(ifIndex, d.queueID, nil)
	if err != nil {
		return fmt.Errorf("error: failed to create XDP socket: %w", err)
	}

	if err := program.SetOption(d.queueID, &xdp.KnotXdpOpts{
		Flags:   uint16(xdp.KNOT_XDP_FILTER_TCP) | uint16(xdp.KNOT_XDP_FILTER_UDP) | uint16(xdp.KNOT_XDP_FILTER_ON),
		UdpPort: 53,
	}); err != nil {
		return fmt.Errorf("error: failed to set XDP options: %w", err)
	}

	// Register our XDP socket file descriptor with the eBPF program so it can be redirected packets
	if err := program.Register(d.queueID, xsk.FD()); err != nil {
		return fmt.Errorf("error: failed to register XDP socket with eBPF program: %w", err)
	}
	defer program.Unregister(d.queueID)

	for {
		// If there are any free slots on the Fill queue...
		if n := xsk.NumFreeFillSlots(); n > 0 {
			// ...then fetch up to that number of not-in-use
			// descriptors and push them onto the Fill ring queue
			// for the kernel to fill them with the received
			// frames.
			xsk.Fill(xsk.GetDescs(n, false))
		}

		// Poll the Rx ring queue for any received frames.
		numRx, _, err := xsk.Poll(-1)
		if err != nil {
			return fmt.Errorf("error: failed to poll the XDP socket: %w", err)
		}

		// Consume the descriptors filled with received frames
		// from the Rx ring queue.
		rxDescs := xsk.Receive(numRx)

		// Print the received frames and also modify them
		// in-place replacing the destination MAC address with
		// broadcast address.
		for i := 0; i < len(rxDescs); i++ {
			d.scheduler.TxChan() <- xsk.GetFrame(rxDescs[i])

		}

	}
	return nil
}

func (d *xdpDriver) Stop(ctx context.Context) error {

	return nil
}

func (d *xdpDriver) TxChan(ch chan []byte) error {
	return nil
}

func (d *xdpDriver) RxChan() chan []byte {

	return nil
}
