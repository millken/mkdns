mkdns
=====
test dns
```
sudo apt-get install ethtool
sudo ethtool -K eth0 gso off
sudo ethtool -K eth0 tso off
sudo ethtool -K eth0 gro off
```

using tcp
```
iptables -I OUTPUT -p tcp --sport 53 --tcp-flags ALL RST -j DROP
iptables -I OUTPUT -p tcp --sport 53 --tcp-flags ALL RST,ACK -j DROP
```

```
CC=$(which musl-gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' server.go

```
==== Link ====

https://github.com/jerome-laforge/ClientInjector/blob/master/src/cmd/ClientInjector/network/network.go
https://github.com/nightcoffee/http-hijack/blob/3bc6b1ec68bab21c94ce75c868928eb7756605b4/http-hijack.go
https://github.com/david415/HoneyBadger
http://www.devdungeon.com/content/packet-capture-injection-and-analysis-gopacket#decoding-packet-layers
https://github.com/grahamking/latency/blob/master/tcp.go
https://github.com/portal-rs/portal-go/blob/main/pkg/server/server.go
https://github.com/miekg/dnsv2/blob/main/rr.go