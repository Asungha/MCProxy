package utils

import (
	"fmt"
	"log"
	. "mc_reverse_proxy/src/common"
)

func Color(text string, color TextColor) string {
	return fmt.Sprintf("%s%s%s", color, text, COLOR_Reset)
}

func Log(serviceName string, serviceNameColor TextColor, textColor TextColor, text string) {
	log.Printf("%s[%s]%s %s%s%s", serviceNameColor, serviceName, COLOR_Reset, textColor, text, COLOR_Reset)
}

func fLogFunc(serviceName string, serviceNameColor, textColor TextColor, textFormat string, text ...any) {
	log.Printf("%s[%s]%s\t%s%s%s", serviceNameColor, serviceName, COLOR_Reset, textColor, fmt.Sprintf(textFormat, text...), COLOR_Reset)
}

func FFatal(serviceName string, serviceNameColor, textColor TextColor, textFormat string, text ...any) {
	log.Fatalf("%s[%s]%s %s%s%s", serviceNameColor, serviceName, COLOR_Reset, textColor, fmt.Sprintf(textFormat, text...), COLOR_Reset)
}

type flog struct {
	Proxy           func(string, ...any)
	Prometheus      func(string, ...any)
	PacketLogger    func(string, ...any)
	HTTPBackend     func(string, ...any)
	HTTPFrontend    func(string, ...any)
	GRPCControl     func(string, ...any)
	Connection      func(string, ...any)
	ApplicationConn func(string, ...any)
}

var FLog = &flog{
	Proxy: func(textFormat string, text ...any) {
		fLogFunc("Game Proxy", PROXY_COLOR, PROXY_COLOR, textFormat, text...)
	},
	Connection: func(textFormat string, text ...any) {
		fLogFunc("Connection", Connection_COLOR_info, Connection_COLOR_info, textFormat, text...)
	},
	ApplicationConn: func(textFormat string, text ...any) {
		fLogFunc("Connection", Connection_COLOR_info, Connection_COLOR_info, textFormat, text...)
	},
	Prometheus: func(textFormat string, text ...any) {
		fLogFunc("Prometheus", METRIC_COLOR, METRIC_COLOR, textFormat, text...)
	},
	PacketLogger: func(textFormat string, text ...any) {
		fLogFunc("Packet Logger", METRIC_COLOR, METRIC_COLOR, textFormat, text...)
	},
	HTTPBackend: func(textFormat string, text ...any) {
		fLogFunc("HTTP Backend", HTTP_COLOR, HTTP_COLOR, textFormat, text...)
	},
	HTTPFrontend: func(textFormat string, text ...any) {
		fLogFunc("HTTP Frontend", HTTP_COLOR, HTTP_COLOR, textFormat, text...)
	},
	GRPCControl: func(textFormat string, text ...any) {
		fLogFunc("GRPC Control", GRPC_COLOR, GRPC_COLOR, textFormat, text...)
	},
}

var FLogDebug = &flog{
	Proxy: func(textFormat string, text ...any) {
		fLogFunc("Game Proxy", PROXY_COLOR_info, PROXY_COLOR_info, "\t"+textFormat, text...)
	},
	Connection: func(textFormat string, text ...any) {
		fLogFunc("Connection", Connection_COLOR_info, Connection_COLOR_info, "\t"+textFormat, text...)
	},
	ApplicationConn: func(textFormat string, text ...any) {
		fLogFunc("Connection", App_conn_color, App_conn_color, "\t"+textFormat, text...)
	},
	Prometheus: func(textFormat string, text ...any) {
		fLogFunc("Prometheus", METRIC_COLOR_info, METRIC_COLOR_info, "\t"+textFormat, text...)
	},
	PacketLogger: func(textFormat string, text ...any) {
		fLogFunc("Packet Logger", METRIC_COLOR_info, METRIC_COLOR_info, "\t"+textFormat, text...)
	},
	HTTPBackend: func(textFormat string, text ...any) {
		fLogFunc("HTTP Backend", HTTP_COLOR_info, HTTP_COLOR_info, "\t"+textFormat, text...)
	},
	HTTPFrontend: func(textFormat string, text ...any) {
		fLogFunc("HTTP Frontend", HTTP_COLOR_info, HTTP_COLOR_info, "\t"+textFormat, text...)
	},
	GRPCControl: func(textFormat string, text ...any) {
		fLogFunc("GRPC Control", GRPC_COLOR_info, GRPC_COLOR_info, "\t"+textFormat, text...)
	},
}

var FLogErr = &flog{
	Proxy: func(textFormat string, text ...any) {
		fLogFunc("Game Proxy", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
	Connection: func(textFormat string, text ...any) {
		fLogFunc("Connection", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
	ApplicationConn: func(textFormat string, text ...any) {
		fLogFunc("Connection", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
	Prometheus: func(textFormat string, text ...any) {
		fLogFunc("Prometheus", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
	PacketLogger: func(textFormat string, text ...any) {
		fLogFunc("Packet Logger", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
	HTTPBackend: func(textFormat string, text ...any) {
		fLogFunc("HTTP Backend", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
	HTTPFrontend: func(textFormat string, text ...any) {
		fLogFunc("HTTP Frontend", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
	GRPCControl: func(textFormat string, text ...any) {
		fLogFunc("GRPC Control", COLOR_Red, COLOR_Red, "\t"+textFormat, text...)
	},
}
