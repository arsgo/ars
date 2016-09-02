package base

import (
	"time"

	"github.com/arsgo/lib4go/logger"
)

var runtimeLogger logger.ILogger
func RunTime(msg string, start time.Time) {
	if runtimeLogger == nil {
		runtimeLogger, _ = logger.Get("run.time." + logger.MainLoggerName)
	}
	now := time.Now()
	tk := now.Sub(start)
	if tk.Nanoseconds()/1000/1000 > 1 {
		runtimeLogger.Fatalf("%s:[%v]", msg, tk)
	}
}
