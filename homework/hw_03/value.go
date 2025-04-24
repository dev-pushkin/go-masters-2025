package cache

import "time"

type val[K comparable, V any] struct {
	k   K
	v   V
	ttl int64
	// функция, которая будет вызвана при истечении ttl
	deleteFunc func(K)
}

func newVal[K comparable, V any](key K, value V, ttl int64) *val[K, V] {
	return &val[K, V]{
		k:   key,
		v:   value,
		ttl: ttl,
	}
}

func (v *val[K, V]) run() {
	if v.ttl <= 0 {
		return
	}
	<-time.NewTimer(time.Duration(v.ttl) * time.Second).C
	v.deleteFunc(v.k)
}
