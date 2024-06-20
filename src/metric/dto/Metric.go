package dto

import "strings"

type Metric struct {
	SystemMetric     `json:"system"`
	NetworkMetric    `json:"newtork"`
	ErrorMetric      `json:"exception"`
	ProxyMetric      `json:"proxy"`
	playerMetric     map[string]PlayerMetric      `json:"player"`
	GameServerMetric map[string]*GameServerMetric `json:"server"`
}

func (m *Metric) GetMetric() string {
	m.ProxyMetric.PlayerPlaying = len(m.playerMetric) - 1
	d := (m.GetNetworkMetric() + m.GetErrorMetric() + m.GetSystemMetric() + m.GetProxyMetric())
	for _, playerMetric := range m.playerMetric {
		d += playerMetric.GetPlayerMetric()
	}
	for _, playerMetric := range m.GameServerMetric {
		d += playerMetric.GetGameServerMetric()
	}
	d = strings.ReplaceAll(d, "\t", "")
	d = strings.ReplaceAll(d, "\n\n", "\n")
	return d
}
