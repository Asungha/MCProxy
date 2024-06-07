package dto

import (
	metricUtils "mc_reverse_proxy/src/metric/utils"
	"strconv"
)

type ErrorMetric struct {
	AcceptFailed            uint
	HandshakeFailed         uint
	PacketDeserializeFailed uint
	HostnameResolveFailed   uint
	ServerConnectFailed     uint

	LogOverflow uint
}

func (m *ErrorMetric) Sum(a ErrorMetric) {
	m.AcceptFailed += a.AcceptFailed
	m.HandshakeFailed += a.HandshakeFailed
	m.PacketDeserializeFailed += m.PacketDeserializeFailed
	m.HostnameResolveFailed += a.HostnameResolveFailed
	m.ServerConnectFailed += a.ServerConnectFailed
	m.LogOverflow += a.LogOverflow
}

func (m *ErrorMetric) GetErrorMetric() string {
	filter := map[string]string{}
	formatter := metricUtils.PrometheusFormatter{}
	formatter.Add("mcproxy_error_accept_failed", strconv.FormatInt(int64(m.AcceptFailed), 10), filter)
	formatter.Add("mcproxy_error_hanhshake_failed", strconv.FormatInt(int64(m.HandshakeFailed), 10), filter)
	formatter.Add("mcproxy_error_deserialization_failed", strconv.FormatInt(int64(m.PacketDeserializeFailed), 10), filter)
	formatter.Add("mcproxy_error_hostname_resolve_failed", strconv.FormatInt(int64(m.HostnameResolveFailed), 10), filter)
	formatter.Add("mcproxy_error_server_connect_failed", strconv.FormatInt(int64(m.ServerConnectFailed), 10), filter)
	formatter.Add("mcproxy_error_log_overflow", strconv.FormatInt(int64(m.LogOverflow), 10), filter)
	return formatter.Get()
}
