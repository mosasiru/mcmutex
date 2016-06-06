package mcmutex

import (
	"github.com/douban/libmc/golibmc"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	mc := golibmc.SimpleNew([]string{"localhost:11211"})
	mutex := NewMCMutex(mc)
	mutex.Expiration = 1
	key := "hoge"
	defer func() {
		if err := mutex.Unlock(key); err != nil {
			t.Error(err)
		}
	}()
	if err := mutex.Lock(key); err != nil {
		t.Error(err)
	}
	if err := mutex.Lock(key); err == nil {
		t.Errorf("lock through")
	}
	time.Sleep(1 * time.Second)
	if err := mutex.Lock(key); err != nil {
		t.Error(err)
	}
}

func TestRetry(t *testing.T) {
	mc := golibmc.SimpleNew([]string{"localhost:11211"})
	key := "fuga"
	mutex := NewMCMutex(mc)
	defer mutex.Unlock(key)
	mutex.Expiration = 1
	mutex.Retry = 10
	mutex.Interval = 1 * time.Millisecond

	if err := mutex.Lock(key); err != nil {
		t.Error(err)
	}

	if err := mutex.Lock(key); err == nil {
		t.Errorf("lock through")
	}

	mutex.Interval = 100 * time.Millisecond
	if err := mutex.Lock(key); err != nil {
		t.Error(err)
	}
}
