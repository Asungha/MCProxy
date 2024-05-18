package main

import (
	"log"
	Logger "mc_reverse_proxy/src/logger"
	proxy "mc_reverse_proxy/src/proxy"
)

func main() {
	logger := Logger.NewLogger()
	err := logger.Init()
	if err != nil {
		log.Printf("[Proxy] Init logger error: %v", err)
	}
	go logger.Start()
	defer logger.Destroy()
	p, err := proxy.NewProxy("25565", logger)
	if err != nil {
		panic(err.Error())
	} else {
		for {
			log.Printf("[Proxy] Accept ready")
			p.Serve()
		}
	}
}
