package logrus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	mu         sync.Mutex
	loggers    = make(map[string]*logHandle)
	syslogHook logrus.Hook
)

type logHandle struct {
	logrus.Logger
	name  string
	level *logrus.Level
}

// GetLogger returns a logger mapped to `name`
func GetLogger(name string) *logHandle {
	mu.Lock()
	defer mu.Unlock()

	if logger, ok := loggers[name]; ok {
		return logger
	}
	loggers[name] = newLogger(name)
	return loggers[name]
}

func newLogger(name string) *logHandle {
	l := &logHandle{name: name}
	l.Out = os.Stderr
	l.Formatter = l
	l.Level = logrus.InfoLevel
	l.Hooks = make(logrus.LevelHooks)
	if syslogHook != nil {
		l.Hooks.Add(syslogHook)
	}
	return l
}

// SetLevel set DEBUG/INFO/WARNING/ERROR log level
func SetLevel(logger *logHandle, lvl string) (err error) {
	logger.Level, err = logrus.ParseLevel(lvl)
	return
}

// Format format log entry
func (l logHandle) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006/01/02 15:04:05.000000")
	str := fmt.Sprintf("%v [%v,%v] <%v>: %v",
		timestamp,
		filepath.Base(entry.Caller.File),
		entry.Caller.Line,
		strings.ToUpper(entry.Level.String()), entry.Message)

	if len(entry.Data) != 0 {
		str += " " + fmt.Sprint(entry.Data)
	}
	str += "\n"
	return []byte(str), nil
}
