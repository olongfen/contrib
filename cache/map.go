package cache

import (
	"sync"
	"time"
)

// CacheMap
type Map struct {
	// 缓存数据
	data sync.Map
	// 时间记录
	time sync.Map
	// 默认有效期超时
	timeout time.Duration
	// 定时清理
	ticker chan bool
}

// Delete 删除数据
func (m *Map) Delete(key interface{}) {
	m.data.Delete(key)
	m.time.Delete(key)
}

// Store 缓存数据
func (m *Map) Store(key, value interface{}) {
	m.data.Store(key, value)
	if m.timeout > 0 {
		m.time.Store(key, time.Now().Add(m.timeout))
	} else {
		m.time.Delete(key)
	}
}

// Load
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	if value, ok = m.data.Load(key); ok {
		// 如果超时
		if _timeout, _ok := m.time.Load(key); _ok {
			if time.Now().After(_timeout.(time.Time)) == true {
				m.Delete(key)
				value = nil
				ok = false
			}
		}
	}
	return
}

// LoadOrStore
func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	if actual, loaded = m.data.LoadOrStore(key, value); loaded == true {
		// 读取:超时验证
		if _timeout, _ok := m.time.Load(key); _ok {
			if time.Now().After(_timeout.(time.Time)) {
				// 新建
				loaded = false
				if actual = value; actual != nil {
					m.Store(key, value)
				} else {
					m.Delete(key)
				}
			}
		}
	} else {
		// 新建
		actual = value
		m.Store(key, actual)
	}
	return
}

// Range
func (m *Map) Range(f func(key, value interface{}) bool) {
	m.data.Range(func(key, value interface{}) bool {
		// 超时检测
		if _timeout, _ok := m.time.Load(key); _ok {
			if time.Now().After(_timeout.(time.Time)) {
				m.Delete(key)
				return true
			}
		}
		return f(key, value)
	})
}

// StoreWithTime
func (m *Map) StoreWithTimeout(key, value interface{}, timeout time.Duration) {
	if m.timeout > 0 {
		// 有默认时长的map
		if timeout > 0 {
			m.data.Store(key, value)
			m.time.Store(key, time.Now().Add(timeout))
		} else {
			m.data.Delete(key)
			m.time.Delete(key)
		}
	} else {
		// 无默认时长的map
		if timeout > 0 {
			m.data.Store(key, value)
			m.time.Store(key, time.Now().Add(m.timeout))
		} else {
			m.data.Store(key, value)
			m.time.Delete(key)
		}
	}
}

// Length 新增方法:map长度
func (m *Map) Length() (ret int) {
	m.Range(func(key, value interface{}) bool {
		ret += 1
		return true
	})
	return
}

// AutoClean 定期清理map
func (m *Map) AutoClean(t time.Duration) *Map {
	//
	if m.ticker != nil {
		m.ticker <- true
		m.ticker = nil
	}

	//
	if t == 0 {
		t = m.timeout
	}

	//
	if t > 0 {
		m.ticker = make(chan bool)
		// 开启定时清理模式
		go func() {
			ticker := time.NewTicker(t)
			for {
				select {
				case <-ticker.C:
					m.Length()
				case <-m.ticker:
					ticker.Stop()
					return
				}
			}
		}()
	}

	return m
}

// new
func NewMap(timeout time.Duration) (m *Map) {
	m = new(Map)
	if timeout > 0 {
		m.timeout = timeout
	}
	return
}
