package common

type TextColor string

const (
	COLOR_Reset   TextColor = "\033[0m"
	COLOR_Red     TextColor = "\033[31m"
	COLOR_Green   TextColor = "\033[32m"
	COLOR_Yellow  TextColor = "\033[33m"
	COLOR_Blue    TextColor = "\033[34m"
	COLOR_Magenta TextColor = "\033[35m"
	COLOR_Cyan    TextColor = "\033[36m"
	COLOR_Gray    TextColor = "\033[37m"
	COLOR_White   TextColor = "\033[97m"

	COLOR_Bright_Green   TextColor = "\033[42m"
	COLOR_Bright_Yellow  TextColor = "\033[43m"
	COLOR_Bright_Magenta TextColor = "\033[45m"
	COLOR_Bright_Cyan    TextColor = "\033[46m"

	PROXY_COLOR  TextColor = COLOR_Bright_Magenta
	METRIC_COLOR TextColor = COLOR_Bright_Cyan
	HTTP_COLOR   TextColor = COLOR_Bright_Green
	GRPC_COLOR   TextColor = COLOR_Bright_Yellow

	PROXY_COLOR_info  TextColor = COLOR_Magenta
	METRIC_COLOR_info TextColor = COLOR_Cyan
	HTTP_COLOR_info   TextColor = COLOR_Green
	GRPC_COLOR_info   TextColor = COLOR_Yellow

	App_conn_color        TextColor = "\033[90m"
	Connection_COLOR_info TextColor = COLOR_White
)
