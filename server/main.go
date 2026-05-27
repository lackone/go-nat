package main

import (
	"flag"
	"go-nat/common"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "c", "", "config file")
	flag.Parse()

	config, err := ParseConfig(confFile)
	if err != nil {
		panic(err)
	}

	listenerCfg, err := ParseListenerConfig(config.ListenerFile)
	if err != nil {
		panic(err)
	}

	sessionMgr := NewSessionManager()

	for _, cfg := range listenerCfg {
		listener := NewListener(&common.ProxyProtocol{
			ClientId:         cfg.ClientId,
			PublicIp:         cfg.PublicIp,
			PublicPort:       cfg.PublicPort,
			PublicProtocol:   cfg.PublicProtocol,
			InternalIp:       cfg.InternalIp,
			InternalPort:     cfg.InternalPort,
			InternalProtocol: cfg.InternalProtocol,
		}, sessionMgr)

		go func() {
			defer listener.Close()
			err := listener.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
	}

	server := NewServer(":9999", sessionMgr)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
