package main

import (
	"fmt"
	"github.com/olongfen/contrib/log"
)

func main() {
	l1, _ := log.NewLog(nil)
	defer log.PanicRecover(l1)
	a := []int{}
	fmt.Println(a[1])

}
