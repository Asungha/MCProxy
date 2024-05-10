package main

import (
	"log"
	proxy "mc_reverse_proxy/src/proxy"
	"runtime"
	"time"
)

func getThreadData() {
	for {
		log.Println("Thread count: ", runtime.NumGoroutine())
		time.Sleep(5 * time.Second)
	}
}

func main() {
	p, err := proxy.NewProxy("25565")
	if err != nil {
		log.Printf(err.Error())
	} else {
		for {
			log.Printf("Accept ready")
			p.Serve()
			if err != nil {
				log.Printf(err.Error())
				break
			}
		}
	}
}
