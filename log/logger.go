package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

type Logger struct {
	*logrus.Logger
	mu        logrus.MutexWrap
	entryPool sync.Pool
}

func (l *Logger) level() logrus.Level {
	level := (*uint32)(&l.Level)
	return logrus.Level(atomic.LoadUint32(level))
}

func (l *Logger) newEntry() *logrus.Entry {
	oldEntry, ok := l.entryPool.Get().(*logrus.Entry)
	var entry *logrus.Entry
	if ok {
		entry = oldEntry
	} else {
		entry = logrus.NewEntry(l.Logger)
	}
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	entry = entry.WithField("file", fmt.Sprintf("%s:%d", file, line))
	//fmt.Println(fmt.Sprintf("aaaa %s:%d", file, line))
	return entry
}

func (logger *Logger) releaseEntry(entry *logrus.Entry) {
	logger.entryPool.Put(entry)
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *logrus.Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *logrus.Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.level() >= logrus.DebugLevel {
		entry := logger.newEntry()
		entry.Debugf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.level() >= logrus.InfoLevel {
		entry := logger.newEntry()
		entry.Infof(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.level() >= logrus.WarnLevel {
		entry := logger.newEntry()
		entry.Warnf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	if logger.level() >= logrus.WarnLevel {
		entry := logger.newEntry()
		entry.Warnf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.level() >= logrus.ErrorLevel {
		entry := logger.newEntry()
		entry.Errorf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.level() >= logrus.FatalLevel {
		entry := logger.newEntry()
		entry.Fatalf(format, args...)
		logger.releaseEntry(entry)
	}
	logrus.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.level() >= logrus.PanicLevel {
		entry := logger.newEntry()
		entry.Panicf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	if logger.level() >= logrus.DebugLevel {
		entry := logger.newEntry()
		entry.Debug(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Info(args ...interface{}) {
	if logger.level() >= logrus.InfoLevel {
		entry := logger.newEntry()
		entry.Info(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Print(args ...interface{}) {
	entry := logger.newEntry()
	entry.Info(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warn(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warn(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Warning(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warn(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if logger.level() >= ErrorLevel {
		entry := logger.newEntry()
		entry.Error(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.level() >= FatalLevel {
		entry := logger.newEntry()
		entry.Fatal(args...)
		logger.releaseEntry(entry)
	}
	logrus.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	if logger.level() >= PanicLevel {
		entry := logger.newEntry()
		entry.Panic(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Debugln(args ...interface{}) {
	if logger.level() >= DebugLevel {
		entry := logger.newEntry()
		entry.Debugln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Infoln(args ...interface{}) {
	if logger.level() >= InfoLevel {
		entry := logger.newEntry()
		entry.Infoln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Println(args ...interface{}) {
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warnln(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warnln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Warningln(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warnln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Errorln(args ...interface{}) {
	if logger.level() >= ErrorLevel {
		entry := logger.newEntry()
		entry.Errorln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Fatalln(args ...interface{}) {
	if logger.level() >= FatalLevel {
		entry := logger.newEntry()
		entry.Fatalln(args...)
		logger.releaseEntry(entry)
	}
	logrus.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	if logger.level() >= PanicLevel {
		entry := logger.newEntry()
		entry.Panicln(args...)
		logger.releaseEntry(entry)
	}
}

func New() *Logger {
	return &Logger{&logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}, logrus.MutexWrap{}, sync.Pool{}}
}
