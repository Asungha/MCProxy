package metric

import (

	// "database/sql"

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
		// IP:        ip,
		// Action:    action,
		// State:     state,
		// Host: func() string {
		// 	if host != nil {
		// 		return *host
		// 	} else {
		// 		return "null"
		// 	}

		// }(),
		// Playername: func() string {
		// 	if playername != nil {
		// 		return *playername
		// 	} else {
		// 		return "null"
		// 	}

		// }(),
		// Message: func() string {
		// 	if message != nil {
		// 		return *message
		// 	} else {
		// 		return "null"
		// }

	}
}

// type Logger struct {
// 	Ready  bool
// 	Ctx    context.Context
// 	Cancle context.CancelFunc
// 	// DB      *sql.DB
// 	repo []Log
// 	// LogChan chan Log
// 	LogEntities []Loggable
// }

// func (l *Logger) Init() error {
// 	l.repo = make([]Log, 256)
// 	l.Ready = true
// 	return nil
// }

// func (l *Logger) Collect()

// func (l *Logger) Start() {

// }

// func (l *Logger) Destroy() {
// 	l.Cancle()
// }

// func NewLogger() *Logger {
// 	ctx, cancle := context.WithCancel(context.Background())
// 	return &Logger{
// 		Ctx:    ctx,
// 		Cancle: cancle,
// 	}
// }
