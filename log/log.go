package log

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
)

type Log struct {
	fd io.Writer
	level int8
}

type logUtil int

func NewLog(w io.Writer, lvl int8) (*Log, error) {
	lu := logUtil(1)
	valid := lu.isValid(lvl)
	if !valid { return nil, ErrInvalidLogLevel }
	return &Log{ fd: w, level: lvl }, nil
}

func (l *Log) Fatal(params ...interface{}) {
	l.log(LOGFATAL, params...)
}

func (l *Log) Fatalf(fmtstr string, params ...interface{}) {
	l.logf(LOGFATAL, fmtstr, params...)
}

func (l *Log) log(lvl int8, params ...interface{}) {
	toLog := l.toLog(lvl)
	if !toLog { return }
	msg := fmt.Sprint(params...)
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
		"%s |%s| %s:%d: %s\n",
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

