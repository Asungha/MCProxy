package main

import (
	"log"
	proxy "mc_reverse_proxy/src/proxy"
)

func main() {
	p, err := proxy.NewProxy("25565")
	// ticker := time.NewTicker(5 * time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			m, err := p.GetMC().Collect()
	// 			if err != nil {
	// 				log.Printf("[MC] error: %v", err)
	// 			}
	// 			log.Printf("%s", m.GetMetric())
	// 		}
	// 	}
	// }()
	if err != nil {
		panic(err.Error())
	} else {
		for {
			log.Printf("[Proxy] Accept ready")
			p.Serve()
		}
	}
}
