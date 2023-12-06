package log

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestCreate(t *testing.T){
	w := bytes.NewBuffer([]byte{})

	_, err := NewLog(w, LOGERR)

	if err != nil { t.Fatal(err) }
}

func TestLogLevel(t *testing.T) {
	w := bytes.NewBuffer([]byte{})
	log, err := NewLog(w, LOGERR)
	if err != nil { t.Fatal(err) }

	res := log.toLog(LOGWARN)
	res1 := log.toLog(LOGERR)
	res2 := log.toLog(LOGFATAL)

	if res || !res1 || !res2 { 
		t.Fatal("unexpected result:", res, res1, res2) 
	}
}

func TestLog(t *testing.T){
	w := bytes.NewBuffer([]byte{})
	log, err := NewLog(w, LOGERR)
	if err != nil { t.Fatal(err) }
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expecting panic")
		} 
		perror, ok := err.(error)
		if !ok {
			t.Fatal("unexpected panic:", err)
		}
		if !errors.Is(ErrLogFatal, perror) {
			t.Fatal("unexpected error type:", err)
		}
	}()

	// expected panic
	log.Fatal("testing log fatal")
}
