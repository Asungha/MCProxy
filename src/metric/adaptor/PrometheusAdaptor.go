package adaptor

import (
	"log"
	"net/http"

	. "mc_reverse_proxy/src/common"
	service "mc_reverse_proxy/src/metric/service"
	"mc_reverse_proxy/src/utils"
)

const serviceName = "Prometheus"

type PrometheusAdaptor struct {
	MetricCollecter *service.MetricService

	ListenAddress string
	serviceName   string
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
	w.Write([]byte(l.GetMetric()))
}

func (e *PrometheusAdaptor) Serve() {
	// <-time.After(1 * time.Second)
	http.HandleFunc("/metrics", e.handler)
	// log.Printf("%s Starting server at %s", utils.Color("[Prometheus Exporter]", utils.COLOR_Cyan), e.ListenAddress)
	utils.FLog.Prometheus("Starting server at %s", e.ListenAddress)
	if err := http.ListenAndServe(e.ListenAddress, nil); err != nil {
		// log.Fatalf("%s Error starting server: %v", utils.Color("[Prometheus Exporter]", utils.COLOR_Red), err)
		utils.FFatal(serviceName, COLOR_Red, COLOR_Red, "Error starting server: %v", err)
	}
}
