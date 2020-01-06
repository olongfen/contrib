package cache

import (
	"testing"
	"time"
)

//
func Test_Map(t *testing.T) {
	m := NewMap(time.Millisecond * 300)
	m.Store("aaa", "111")
	m.Store("bbb", "222")
	m.StoreWithTimeout("ccc", "333", time.Millisecond*400)

	for _i, _t := range []int{0, 200, 100, 200} {
		time.Sleep(time.Millisecond * time.Duration(_t))
		t.Logf("#%d length:%d", _i, m.Length())
	}
}
