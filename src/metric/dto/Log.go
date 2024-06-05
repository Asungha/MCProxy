package dto

import "time"

type Log struct {
	Timestamp time.Time
	NetworkMetric
	ErrorMetric
	ProxyMetric
	PlayerMetric
}

func NewLog() Log {
	return Log{
		Timestamp: time.Now(),
	}
}
