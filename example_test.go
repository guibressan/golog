package golog_test

import (
	"os"
	"testing"

	"github.com/guibressan/golog"
)

func TestExampleLogDefault(t *testing.T) {
	log, _ := golog.NewLog();

	log.Info(1, "==", 1)
}

func TestExampleLogWithConfig(t *testing.T) {
	log, _ := golog.NewLog(
		golog.WithLevel(golog.LOGINFO),
		golog.WithWriter(os.Stderr),
	);

	log.Infof("1 == %d\n", 1)
}
