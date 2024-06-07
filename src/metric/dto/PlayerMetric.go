package dto

import (
	"fmt"
	metricUtils "mc_reverse_proxy/src/metric/utils"
	"strconv"
	"strings"
	"time"
)

type PlayerMetric struct {
	LoggedOut  bool
	PlayerName string
	IP         string
	Port       string
	LogginTime time.Time
	Playtime   time.Duration
	*NetworkMetric
	ErrorMetric
}

func NewPlayerMetric(addr string, name string) *PlayerMetric {
	s := strings.Split(addr, ":")
	return &PlayerMetric{
		LogginTime: time.Now(),
		PlayerName: name,
		IP:         s[0],
		Port:       s[1],
	}
}

func (m *PlayerMetric) GetPlayerMetric() string {
	if m.PlayerName == "" || m.IP == "" {
		return ""
	}
	filter := map[string]string{"player": m.PlayerName, "ip": m.IP}
	formatter := metricUtils.PrometheusFormatter{}
	formatter.Add("mcproxy_player_online", "1", filter)
	formatter.Add("mcproxy_player_playtime", fmt.Sprint(time.Since(m.LogginTime).Seconds()), filter)

	formatter.Add("mcproxy_player_error_accept_failed", fmt.Sprint(m.AcceptFailed), filter)
	formatter.Add("mcproxy_player_error_hanhshake_failed", fmt.Sprint(m.HandshakeFailed), filter)
	formatter.Add("mcproxy_player_error_deserialization_failed", fmt.Sprint(m.PacketDeserializeFailed), filter)
	formatter.Add("mcproxy_player_error_hostname_resolve_failed", fmt.Sprint(m.HostnameResolveFailed), filter)
	formatter.Add("mcproxy_player_error_server_connect_failed", fmt.Sprint(m.ServerConnectFailed), filter)
	formatter.Add("mcproxy_player_error_server_connect_failed", fmt.Sprint(m.ServerConnectFailed), filter)

	if m.NetworkMetric != nil {
		formatter.Add("mcproxy_player_network_client_packet_tx", strconv.FormatInt(int64(m.NetworkMetric.ClientPacketTx), 10), filter)
		formatter.Add("mcproxy_player_network_client_packet_rx", strconv.FormatInt(int64(m.NetworkMetric.ClientPacketRx), 10), filter)
		formatter.Add("mcproxy_player_network_server_packet_tx", strconv.FormatInt(int64(m.NetworkMetric.ServerPacketTx), 10), filter)
		formatter.Add("mcproxy_player_network_server_packet_rx", strconv.FormatInt(int64(m.NetworkMetric.ServerPacketRx), 10), filter)
		formatter.Add("mcproxy_player_network_client_data_tx", strconv.FormatInt(int64(m.NetworkMetric.ClientDataTx), 10), filter)
		formatter.Add("mcproxy_player_network_client_data_rx", strconv.FormatInt(int64(m.NetworkMetric.ClientPacketRx), 10), filter)
		formatter.Add("mcproxy_player_network_server_data_tx", strconv.FormatInt(int64(m.NetworkMetric.ServerDataTx), 10), filter)
		formatter.Add("mcproxy_player_network_server_data_rx", strconv.FormatInt(int64(m.NetworkMetric.ServerDataRx), 10), filter)
	}

	return formatter.Get()
}
