package main

import (
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
	// f, err := os.Create("goroutine.prof")
	// if err != nil {
	// 	log.Fatal("could not create goroutine profile: ", err)
	// }
	// defer f.Close()
	// go func() {
	// 	<-time.After(1 * time.Minute)
	// 	log.Printf("[Profiler] writting")
	// 	if gprof := pprof.Lookup("goroutine"); gprof != nil {
	// 		if err := gprof.WriteTo(f, 0); err != nil {
	// 			log.Fatal("could not write goroutine profile: ", err)
	// 		}
	// 	}
	// }()
	if err != nil {
		panic(err.Error())
	} else {
		p.Serve()
	}

}
