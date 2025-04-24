package cache

import (
	"testing"
	"time"
)

func Test_Cache_SetWithTtl(t *testing.T) {
	type args struct {
		key   string
		value int
		ttl   int64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test cache with TTL",
			args: args{
				key:   "testKey",
				value: 42,
				ttl:   1,
			},
			want: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New[string, int]()
			c.SetWithTTL(tt.args.key, tt.args.value, tt.args.ttl)
			if got := c.Get(tt.args.key); got != tt.want {
				t.Errorf("Cache value = %v, want %v", got, tt.want)
			}
			// Wait for the TTL to expire
			time.Sleep(time.Duration(tt.args.ttl+1) * time.Second)
			if got := c.Get(tt.args.key); got != 0 {
				t.Errorf("Cache value = %v, want %v", got, 0)
			}
		})
	}
}

func Test_Cache_Set(t *testing.T) {
	type args struct {
		key   string
		value int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test cache set",
			args: args{
				key:   "testKey",
				value: 42,
			},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New[string, int]()
			c.Set(tt.args.key, tt.args.value)
			if got := c.Get(tt.args.key); got != tt.want {
				t.Errorf("Cache value = %v, want %v", got, tt.want)
			}
		})
	}
}
