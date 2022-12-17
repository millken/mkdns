package xdp

// Only support Little Endian
// 请务必不要在Big Endian的机器上使用，程序没有设计对其的支持

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -cc clang -cflags "-O2" xsk ./bpf/xsk_new.c
