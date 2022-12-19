https://github.com/facebookincubator/katran/blob/main/katran/decap/bpf/decap_kern.c
https://gitlab.nic.cz/knot/knot-dns/-/blob/master/src/libknot/xdp/bpf-kernel.c
https://github.com/OSH-2019/x-xdp-on-android/blob/master/docs/research.md
https://github.com/cloudflare/cbpfc/blob/master/c_example_test.go

Debian下libbpf编译失败提示<asm/types.h>文件不存在解决方法
apt install gcc-multilib
clang -I include/ -O2 -target bpf -c xsk.c -o xsk.o
# compile to .o
clang -O2 -target bpf -c xsk.c -o xsk.o
# compile to .s
clang -O2 -target bpf -c xsk.c -S -o xsk.s

# detach
sudo ip link set dev ens19 xdp off

# attach
sudo ip link set dev ens19 xdp obj xsk.o sec xsk_program
# also you can use cilium ebpf loader