package metric

import (

	// "database/sql"

	"fmt"
	"time"
	// _ "github.com/mattn/go-sqlite3"
)

type Action uint

func (a *Action) String() string {
	switch *a {
	case 1:
		return "connect"
	case 2:
		return "handshake"
	case 3:
		return "status"
	case 4:
		return "login"
	case 5:
		return "disconnect"
	default:
		return "unknown"
	}
}

const (
	UNKNOWN    Action = 0
	CONNECT    Action = 1
	HANDSHAKE  Action = 2
	STATUS     Action = 3
	LOGIN      Action = 4
	DISCONNECT Action = 5
)

type Log struct {
	Timestamp time.Time
	// IP         string
	// Action     Action
	// State      string
	// Host       string
	// Playername string
	// Message    string
	NetworkMetric
	ErrorMetric
	ProxyMetric
	PlayerMetric
}

func NewLog(ip string, action Action, state string, host *string, playername *string, message *string) Log {
	return Log{
		Timestamp: time.Now(),
	}
}

type LogPusher struct {
	Collector MetricCollecter
}

func (p *LogPusher) PushErrorMetric(log ErrorMetric) error {
	fmt.Println("ErrorMetric")
	return p.Collector.PushLog(Log{ErrorMetric: log})
}

func (p *LogPusher) PushProxyMetric(log ProxyMetric) error {
	fmt.Println("ProxyMetric")
	return p.Collector.PushLog(Log{ProxyMetric: log})
}

func (p *LogPusher) PushPlayerMetric(log PlayerMetric) error {
	fmt.Println("PlayerMetric")
	return p.Collector.PushLog(Log{PlayerMetric: log})
}

func (p *LogPusher) PushNetworkMetric(log NetworkMetric) error {
	fmt.Println("NetworkMetric")
	return p.Collector.PushLog(Log{NetworkMetric: log})
}
