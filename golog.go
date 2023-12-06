package golog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

const (
	LOGFATAL int8 = iota + 1
	LOGERR
	LOGWARN
	LOGINFO
	LOGDEBUG
	LOGTRACE
)

var (
	ErrInvalidLogLevel = errors.New("log: invalid log level")
	ErrLogFatal = errors.New("log: error log FATAL")
)

// Log is the struct of the logger
type Log struct {
	fd io.Writer
	level int8
}

type logUtil int

// CfgFunc is a function that changes logger configuration
type CfgFunc func(*Log)

// WithWriter sets the logging writer
func WithWriter(writer io.Writer) CfgFunc {
	return func(l *Log) {
		l.fd = writer
	}
}

// WithLevel sets the logging level
func WithLevel(level int8) CfgFunc {
	return func(l *Log) {
		l.level = level
	}
}

// NewLog creates a new instance of Log
// Returns error on invalid log level
func NewLog(cfgFuncs ...CfgFunc) (*Log, error) {
	lu := logUtil(0)
	log := lu.defaultLog();
	for _, f := range cfgFuncs {
		f(log)
	}
	valid := lu.isValid(log.level)
	if !valid { return nil, ErrInvalidLogLevel }
	return log, nil
}

// Fatal prints the fatal log and panics
func (l *Log) Fatal(params ...interface{}) {
	l.log(LOGFATAL, params...)
	panic(ErrLogFatal)
}

// Fatalf prints the fatal log and panics
func (l *Log) Fatalf(fmtstr string, params ...interface{}) {
	l.logf(LOGFATAL, fmtstr, params...)
	panic(ErrLogFatal)
}

// Err prints the err log
func (l *Log) Err(params ...interface{}) {
	l.log(LOGERR, params...)
}

// Errf prints the err log
func (l *Log) Errf(fmtstr string, params ...interface{}) {
	l.logf(LOGERR, fmtstr, params...)
}

// Warn prints the warn log
func (l *Log) Warn(params ...interface{}) {
	l.log(LOGWARN, params...)
}

// Warnf prints the warn log
func (l *Log) Warnf(fmtstr string, params ...interface{}) {
	l.logf(LOGWARN, fmtstr, params...)
}

// Info prints the info log
func (l *Log) Info(params ...interface{}) {
	l.log(LOGINFO, params...)
}

// Infof prints the info log
func (l *Log) Infof(fmtstr string, params ...interface{}) {
	l.logf(LOGINFO, fmtstr, params...)
}

// Debug prints the debug log
func (l *Log) Debug(params ...interface{}) {
	l.log(LOGDEBUG, params...)
}

// Debugf prints the debug log
func (l *Log) Debugf(fmtstr string, params ...interface{}) {
	l.logf(LOGDEBUG, fmtstr, params...)
}

// Trace prints the trace log
func (l *Log) Trace(params ...interface{}) {
	l.log(LOGTRACE, params...)
}

// Tracef prints the trace log
func (l *Log) Tracef(fmtstr string, params ...interface{}) {
	l.logf(LOGTRACE, fmtstr, params...)
}

func (l *Log) log(lvl int8, params ...interface{}) {
	toLog := l.toLog(lvl)
	if !toLog { return }
	msg := fmt.Sprintln(params...)
	l.logWrite(lvl, msg)
}

func (l *Log) logf(lvl int8, fmtstr string, params ...interface{}) {
	toLog := l.toLog(lvl)
	if !toLog { return }
	msg := fmt.Sprintf(fmtstr, params...)
	l.logWrite(lvl, msg)
}

func (l *Log) logWrite(lvl int8, msg string) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		fmt.Fprintln(os.Stderr, "log: err getting caller")
	}
	log := fmt.Sprintf(
		"%s |%s| %s:%d: %s",
		time.Now().Format(time.DateTime),
		logUtil(0).lvlStr(lvl),
		file,
		line,
		msg,
	)
	n, err := l.fd.Write([]byte(log))
	if err != nil {
		fmt.Fprintln(os.Stderr, "log: err writing log")
	}
	if n != len(log) {
		fmt.Fprintln(os.Stderr, "log: err unexpected bytes written:", len(log), n)
	}
}

func (l *Log) toLog(lvl int8) bool {
	return l.level >= lvl
}

func (l logUtil) isValid(lvl int8) bool {
	return lvl >= LOGFATAL && lvl <= LOGTRACE
}

func (l logUtil) defaultLog() *Log {
	return &Log{
		fd: os.Stderr,
		level: LOGINFO,
	}
}

func (l logUtil) lvlStr(lvl int8) string {
	switch lvl {
	case LOGFATAL:
		return "FATAL"
	case LOGERR:
		return "ERR"
	case LOGWARN:
		return "WARN"
	case LOGINFO:
		return "INFO"
	case LOGDEBUG:
		return "DEBUG"
	case LOGTRACE:
		return "TRACE"
	}
	return ""
}
