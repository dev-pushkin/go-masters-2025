package cache

import "time"

type val[V any] struct {
	v   V
	ttl int64
	// функция, которая будет вызвана при истечении ttl
	deleteFunc func()
}

func newVal[V any](value V, ttl int64) *val[V] {
	return &val[V]{
		v:   value,
		ttl: ttl,
	}
}

func (v *val[V]) run() {
	if v.ttl <= 0 {
		return
	}
	<-time.NewTimer(time.Duration(v.ttl) * time.Second).C
	v.deleteFunc()
}
