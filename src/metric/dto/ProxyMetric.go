package dto

import (
	metricUtils "mc_reverse_proxy/src/metric/utils"
	"strconv"
)

type ProxyMetric struct {
	PlayerGetStatus int `json:"get_status"`
	PlayerLogin     int `json:"login"`
	PlayerPlaying   int `json:"playing"`
	Ping            int `json:"ping"`
	Connected       int `json:"connected"`
}

func (m *ProxyMetric) GetProxyMetric() string {
	filter := map[string]string{}
	formatter := metricUtils.PrometheusFormatter{}
	formatter.Add("mcproxy_proxy_get_status", strconv.FormatInt(int64(m.PlayerGetStatus), 10), filter)
	formatter.Add("mcproxy_proxy_login", strconv.FormatInt(int64(m.PlayerLogin), 10), filter)
	formatter.Add("mcproxy_proxy_ping", strconv.FormatInt(int64(m.Ping), 10), filter)
	formatter.Add("mcproxy_proxy_playing", strconv.FormatInt(int64(m.PlayerPlaying), 10), filter)
	formatter.Add("mcproxy_proxy_connected", strconv.FormatInt(int64(m.Connected), 10), filter)
	return formatter.Get()
}

func (m *ProxyMetric) Sum(a ProxyMetric) {
	m.PlayerGetStatus += a.PlayerGetStatus
	m.PlayerLogin += a.PlayerLogin
	m.PlayerPlaying += a.PlayerPlaying
	m.Connected += a.Connected
}
