package main

type ClientInfo struct {
	ClientId         string // 客户端ID
	PublicIp         string // 公共IP
	PublicPort       int    // 公共端口
	PublicProtocol   string // 公共协议
	InternalIp       string // 内网IP
	InternalPort     int    // 内网端口
	InternalProtocol string // 内网协议
}
