FROM centos:7
MAINTAINER millken@gmail.com

ADD https://github.com/millken/mkdns/releases/download/v20161026-test/mkdns.tar.gz /usr/share/mkdns.tar.gz
RUN yum update -y \
 && yum install -y libpcap-devel geoip-devel net-tools iptables \
 && cd /usr/share/ && tar xvf mkdns.tar.gz && rm mkdns.tar.gz \
 && mv mkdns /sbin/mkdns \
 && mv config.toml /etc/mkdns.toml \
 && iptables -I OUTPUT -p tcp --sport 53 --tcp-flags ALL RST -j DROP \
 && iptables -I OUTPUT -p tcp --sport 53 --tcp-flags ALL RST,ACK -j DROP 

EXPOSE 53/udp 53/tcp 10000/tcp

CMD ["/sbin/mkdns", "-c", "/etc/mkdns.toml"]
