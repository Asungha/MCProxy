package main

import (
	"log"
	proxy "mc_reverse_proxy/src/proxy"
)

func main() {
	p, err := proxy.NewProxy("25565")
	if err != nil {
		panic(err.Error())
	} else {
		for {
			log.Printf("[Proxy] Accept ready")
			p.Serve()
		}
	}
}
