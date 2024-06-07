package dto

import "strings"

type Metric struct {
	SystemMetric
	NetworkMetric
	ErrorMetric
	ProxyMetric
	playerMetric map[string]PlayerMetric
}

func (m *Metric) GetMetric() string {
	m.ProxyMetric.PlayerPlaying = len(m.playerMetric) - 1
	d := (m.GetNetworkMetric() + m.GetErrorMetric() + m.GetSystemMetric() + m.GetProxyMetric())
	for _, playerMetric := range m.playerMetric {
		d += playerMetric.GetPlayerMetric()
	}
	d = strings.ReplaceAll(d, "\t", "")
	d = strings.ReplaceAll(d, "\n\n", "\n")
	return d
}
