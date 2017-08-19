package log_test

import (
	"github.com/Desserts/goutil/log"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	log.Info("test")
}

func TestLog2(t *testing.T) {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
	log.Error("bbb")
}

func TestFile(t *testing.T) {
	file, err := os.OpenFile("./a.log", os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		t.Error("create file error", err)
	}
	log.SetOutput(file)
	log.Error("HAHA")
}
