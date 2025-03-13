package driver

import (
	"context"
	"net"
)

// DNSRequest 表示一个 DNS 请求
type DNSRequest struct {
	RawData  []byte // 原始 DNS 数据
	ClientIP net.IP // 客户端 IP 地址
}

// DNSResponse 表示一个 DNS 响应
type DNSResponse struct {
	RawData []byte // 原始 DNS 数据
}

// Driver 是驱动接口
type Driver2 interface {
	// Start 启动驱动
	Start(ctx context.Context) error

	// Stop 停止驱动
	Stop() error

	// Receive 接收 DNS 请求
	Receive() <-chan DNSRequest

	// Send 发送 DNS 响应
	Send(response DNSResponse) error
}
