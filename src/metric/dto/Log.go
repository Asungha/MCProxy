package dto

import "time"

type Log struct {
	Timestamp time.Time
	NetworkMetric
	ErrorMetric
	ProxyMetric
	PlayerMetric
	SystemMetric
	*GameServerMetric
}

func NewLog() Log {
	return Log{
		Timestamp: time.Now(),
	}
}
