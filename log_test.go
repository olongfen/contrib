package log

import (
	"testing"
)

func Test_File(t *testing.T) {
	var (
		l     = NewLogFile("./test.log")
		l1, _ = NewLog(l)
	)
	l1.Infof("ha ha")
	l1.Errorf("la la")
}
