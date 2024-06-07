package dto

import (
	"fmt"
)

type PrometheusConfig struct {
	ListenIP string
	Port     int
}

func (c *PrometheusConfig) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.ListenIP, c.Port)
}
