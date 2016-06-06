package mcmutex

import (
	"errors"
	"github.com/douban/libmc/golibmc"
	"time"
)

const (
	DefaultInterval   = 10 * time.Millisecond
	DefaultRetry      = 0
	DefaultExpiration = 30
)

var ErrLockFailed = errors.New("failed to acquire lock")

type MCMutex struct {
	Interval   time.Duration
	Retry      int
	Expiration int64
	client     *golibmc.Client
}

func NewMCMutex(mc *golibmc.Client) *MCMutex {
	return &MCMutex{
		client:     mc,
		Interval:   DefaultInterval,
		Retry:      DefaultRetry,
		Expiration: DefaultExpiration,
	}
}

func (m *MCMutex) Lock(key string) error {
	for i := 0; i <= m.Retry; i++ {
		if err := m.client.Add(&golibmc.Item{Key: key, Value: []byte{1}, Expiration: m.Expiration}); err != nil {
			if err != golibmc.ErrNotStored {
				return err
			}
			time.Sleep(m.Interval)
			continue
		}
		return nil
	}
	return ErrLockFailed
}

func (m *MCMutex) Unlock(key string) error {
	return m.client.Delete(key)
}
