package xdp

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
	"github.com/vishvananda/netlink"
)

type KnotXdpFilter uint16

const (
	KNOT_XDP_FILTER_ON    KnotXdpFilter = 1 << 0 /*!< Filter enabled. */
	KNOT_XDP_FILTER_UDP   KnotXdpFilter = 1 << 1 /*!< Apply filter to UDP. */
	KNOT_XDP_FILTER_TCP   KnotXdpFilter = 1 << 2 /*!< Apply filter to TCP. */
	KNOT_XDP_FILTER_QUIC  KnotXdpFilter = 1 << 3 /*!< Apply filter to QUIC/UDP. */
	KNOT_XDP_FILTER_PASS  KnotXdpFilter = 1 << 4 /*!< Pass incoming messages to ports >= port value. */
	KNOT_XDP_FILTER_DROP  KnotXdpFilter = 1 << 5 /*!< Drop incoming messages to ports >= port value. */
	KNOT_XDP_FILTER_ROUTE KnotXdpFilter = 1 << 6
)

type KnotXdpOpts struct {
	Flags    uint16 /*!< XDP filter flags \a knot_xdp_filter_flag_t. */
	UdpPort  uint16 /*!< UDP/TCP port to listen on. */
	QuicPort uint16 /*!< QUIC/UDP port to listen on. */
}

// Program represents the necessary data structures for a simple XDP program that can filter traffic
// based on the attached rx queue.
type Program struct {
	Program *ebpf.Program
	Options *ebpf.Map
	Sockets *ebpf.Map
}

// Attach the XDP Program to an interface.
func (p *Program) Attach(Ifindex int) error {
	if err := removeProgram(Ifindex); err != nil {
		return err
	}
	return attachProgram(Ifindex, p.Program)
}

// Detach the XDP Program from an interface.
func (p *Program) Detach(Ifindex int) error {
	return removeProgram(Ifindex)
}

// Register adds the socket file descriptor as the recipient for packets from the given queueID.
func (p *Program) Register(queueID int, fd int) error {
	if err := p.Sockets.Put(uint32(queueID), uint32(fd)); err != nil {
		return fmt.Errorf("failed to update xsksMap: %v", err)
	}

	return nil
}

func (p *Program) SetOption(queueID int, opts *KnotXdpOpts) error {
	if err := p.Options.Put(uint32(queueID), opts); err != nil {
		return fmt.Errorf("failed to update optsMap: %v", err)
	}
	return nil
}

// Unregister removes any associated mapping to sockets for the given queueID.
func (p *Program) Unregister(queueID int) error {
	if err := p.Options.Put(uint32(queueID), &KnotXdpOpts{}); err != nil {
		return err
	}
	if err := p.Sockets.Delete(uint32(queueID)); err != nil {
		return err
	}
	return nil
}

// Close closes and frees the resources allocated for the Program.
func (p *Program) Close() error {
	allErrors := []error{}
	if p.Sockets != nil {
		if err := p.Sockets.Close(); err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to close xsksMap: %v", err))
		}
		p.Sockets = nil
	}

	if p.Options != nil {
		if err := p.Options.Close(); err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to close optsMap: %v", err))
		}
		p.Options = nil
	}

	if p.Program != nil {
		if err := p.Program.Close(); err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to close XDP program: %v", err))
		}
		p.Program = nil
	}

	if len(allErrors) > 0 {
		return allErrors[0]
	}
	return nil
}

// NewProgram returns a translation of the default eBPF XDP program found in the
// xsk_load_xdp_prog() function in <linux>/tools/lib/bpf/xsk.c:
// https://github.com/torvalds/linux/blob/master/tools/lib/bpf/xsk.c#L259
func NewProgram(maxQueueEntries int) (*Program, error) {
	optsMap, err := ebpf.NewMap(&ebpf.MapSpec{
		Name:       "opts_map",
		Type:       ebpf.Array,
		KeySize:    uint32(unsafe.Sizeof(int32(0))),
		ValueSize:  uint32(unsafe.Sizeof(int32(0))),
		MaxEntries: uint32(maxQueueEntries),
		Flags:      0,
		InnerMap:   nil,
	})
	if err != nil {
		return nil, fmt.Errorf("ebpf.NewMap opts_map failed (try increasing RLIMIT_MEMLOCK): %v", err)
	}

	xsksMap, err := ebpf.NewMap(&ebpf.MapSpec{
		Name:       "xsks_map",
		Type:       ebpf.XSKMap,
		KeySize:    uint32(unsafe.Sizeof(int32(0))),
		ValueSize:  uint32(unsafe.Sizeof(int32(0))),
		MaxEntries: uint32(maxQueueEntries),
		Flags:      0,
		InnerMap:   nil,
	})
	if err != nil {
		return nil, fmt.Errorf("ebpf.NewMap xsks_map failed (try increasing RLIMIT_MEMLOCK): %v", err)
	}

	program, err := ebpf.NewProgram(&ebpf.ProgramSpec{
		Name: "xsk_ebpf",
		Type: ebpf.XDP,
		Instructions: asm.Instructions{
			{OpCode: 97, Dst: 1, Src: 1, Offset: 16},                               // 0: code: 97 dst_reg: 1 src_reg: 1 off: 16 imm: 0   // 0
			{OpCode: 99, Dst: 10, Src: 1, Offset: -4},                              // 1: code: 99 dst_reg: 10 src_reg: 1 off: -4 imm: 0  // 1
			{OpCode: 191, Dst: 2, Src: 10},                                         // 2: code: 191 dst_reg: 2 src_reg: 10 off: 0 imm: 0  // 2
			{OpCode: 7, Dst: 2, Src: 0, Offset: 0, Constant: -4},                   // 3: code: 7 dst_reg: 2 src_reg: 0 off: 0 imm: -4    // 3
			{OpCode: 24, Dst: 1, Src: 1, Offset: 0, Constant: int64(optsMap.FD())}, // 4: code: 24 dst_reg: 1 src_reg: 1 off: 0 imm: 4    // 4 XXX use optsMap.FD as IMM
			//{ OpCode: 0 },                                                                 // 5: code: 0 dst_reg: 0 src_reg: 0 off: 0 imm: 0     //   part of the same instruction
			{OpCode: 133, Dst: 0, Src: 0, Constant: 1},                  // 6: code: 133 dst_reg: 0 src_reg: 0 off: 0 imm: 1   // 5
			{OpCode: 191, Dst: 1, Src: 0},                               // 7: code: 191 dst_reg: 1 src_reg: 0 off: 0 imm: 0   // 6
			{OpCode: 180, Dst: 0, Src: 0},                               // 8: code: 180 dst_reg: 0 src_reg: 0 off: 0 imm: 0   // 7
			{OpCode: 21, Dst: 1, Src: 0, Offset: 8},                     // 9: code: 21 dst_reg: 1 src_reg: 0 off: 8 imm: 0    // 8
			{OpCode: 180, Dst: 0, Src: 0, Constant: 2},                  // 10: code: 180 dst_reg: 0 src_reg: 0 off: 0 imm: 2  // 9
			{OpCode: 97, Dst: 1, Src: 1},                                // 11: code: 97 dst_reg: 1 src_reg: 1 off: 0 imm: 0   // 10
			{OpCode: 21, Dst: 1, Offset: 5},                             // 12: code: 21 dst_reg: 1 src_reg: 0 off: 5 imm: 0   // 11
			{OpCode: 24, Dst: 1, Src: 1, Constant: int64(xsksMap.FD())}, // 13: code: 24 dst_reg: 1 src_reg: 1 off: 0 imm: 5   // 12 XXX use xsksMap.FD as IMM
			//{ OpCode: 0 },                                                                 // 14: code: 0 dst_reg: 0 src_reg: 0 off: 0 imm: 0    //    part of the same instruction
			{OpCode: 97, Dst: 2, Src: 10, Offset: -4}, // 15: code: 97 dst_reg: 2 src_reg: 10 off: -4 imm: 0 // 13
			{OpCode: 180, Dst: 3},                     // 16: code: 180 dst_reg: 3 src_reg: 0 off: 0 imm: 0  // 14
			{OpCode: 133, Constant: 51},               // 17: code: 133 dst_reg: 0 src_reg: 0 off: 0 imm: 51 // 15
			{OpCode: 149},                             // 18: code: 149 dst_reg: 0 src_reg: 0 off: 0 imm: 0  // 16
		},
		License:       "LGPL-2.1 or BSD-2-Clause",
		KernelVersion: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("error: ebpf.NewProgram failed: %v", err)
	}

	return &Program{Program: program, Options: optsMap, Sockets: xsksMap}, nil
}

// removeProgram removes an existing XDP program from the given network interface.
func removeProgram(Ifindex int) error {
	var link netlink.Link
	var err error
	link, err = netlink.LinkByIndex(Ifindex)
	if err != nil {
		return err
	}
	if !isXdpAttached(link) {
		return nil
	}
	if err = netlink.LinkSetXdpFd(link, -1); err != nil {
		return fmt.Errorf("netlink.LinkSetXdpFd(link, -1) failed: %v", err)
	}
	for {
		link, err = netlink.LinkByIndex(Ifindex)
		if err != nil {
			return err
		}
		if !isXdpAttached(link) {
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}

func isXdpAttached(link netlink.Link) bool {
	if link.Attrs() != nil && link.Attrs().Xdp != nil && link.Attrs().Xdp.Attached {
		return true
	}
	return false
}

// attachProgram attaches the given XDP program to the network interface.
func attachProgram(Ifindex int, program *ebpf.Program) error {
	link, err := netlink.LinkByIndex(Ifindex)
	if err != nil {
		return err
	}
	return netlink.LinkSetXdpFdWithFlags(link, program.FD(), int(DefaultXdpFlags))
}
