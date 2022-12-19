package ebpf

import (
	"github.com/cilium/ebpf"
	"github.com/millken/mkdns/internal/xdp"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go xsk  ./bpf_kernel.c -- -I/usr/include/  -I./include  -nostdinc -O3
////go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -cc clang -cflags "-O2" xsk ./bpf_kernel.c -- -I./include  -nostdinc -O3
// NewDNSProtoProgram returns an new eBPF that directs packets of the given ip protocol to to XDP sockets
func NewDNSProtoProgram(options *ebpf.CollectionOptions) (*xdp.Program, error) {
	spec, err := loadXsk()
	if err != nil {
		return nil, err
	}

	var program xskObjects
	if err := spec.LoadAndAssign(&program, options); err != nil {
		return nil, err
	}

	p := &xdp.Program{Program: program.XdpRedirectDnsFunc, Options: program.OptsMap, Sockets: program.XsksMap}
	return p, nil
}
