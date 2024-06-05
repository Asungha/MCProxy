package adaptor

import (
	"log"
	"net/http"
	"time"

	service "mc_reverse_proxy/src/metric/service"
)

type PrometheusAdaptor struct {
	MetricCollecter *service.MetricService
	ListenAddress   string
	MetricAdaptor
}

func (e *PrometheusAdaptor) handler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to plain text
	w.Header().Set("Content-Type", "text/plain")

	// Write the text response
	l, err := e.MetricCollecter.Collect()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Printf("[Prometheus Exporter] Error: %v", err)
		return
	}
	// fmt.Println(l.GetMetric())
	w.Write([]byte(l.GetMetric()))
}

func (e *PrometheusAdaptor) Serve() {
	<-time.After(1 * time.Second)
	http.HandleFunc("/metrics", e.handler)
	log.Printf("[Prometheus Exporter] Starting server at %s", e.ListenAddress)
	if err := http.ListenAndServe(e.ListenAddress, nil); err != nil {
		log.Fatalf("[Prometheus Exporter] Error starting server: %v", err)
	}
}
