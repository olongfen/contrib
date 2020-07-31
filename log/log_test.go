package log

import (
	"fmt"
	"testing"
)

var (
	l = NewLogFile("./test", true)
)

func Test_File(t *testing.T) {
	l.Infof("ha ha")
	l.Errorf("la la")
}

func TestPanicRecover(t *testing.T) {
	var (
		a = []int{}
	)
	defer l.Recover()
	fmt.Println(a[1])
}
