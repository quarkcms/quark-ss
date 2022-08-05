package log

import (
	"time"

	"github.com/cihub/seelog"
)

func init() {
	currentTime := time.Now()
	config := `<seelog>
					<outputs formatid="main">
						<file path='./storage/logs/` + currentTime.Format("2006-01-02") + `.log'/>
					</outputs>
					<formats>
        				<format id="main" format="%Date %Time - [%LEVEL] - %Msg%n"/>
    				</formats>
				</seelog>`

	logger, _ := seelog.LoggerFromConfigAsBytes([]byte(config))
	seelog.ReplaceLogger(logger)
}

// Trace formats message using the default formats for its operands and writes to default logger with log level = Trace
func Trace(v ...interface{}) {

	seelog.Trace(v)
}

// Debug formats message using the default formats for its operands and writes to default logger with log level = Debug
func Debug(v ...interface{}) {

	seelog.Debug(v)
}

// Info formats message using the default formats for its operands and writes to default logger with log level = Info
func Info(v ...interface{}) {

	seelog.Info(v)
}

// Warn formats message using the default formats for its operands and writes to default logger with log level = Warn
func Warn(v ...interface{}) {

	seelog.Warn(v)
}

// Error formats message using the default formats for its operands and writes to default logger with log level = Error
func Error(v ...interface{}) {

	seelog.Error(v)
}

// Critical formats message using the default formats for its operands and writes to default logger with log level = Critical
func Critical(v ...interface{}) {

	seelog.Critical(v)
}
