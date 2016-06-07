// Package mcmutex provides a mutex using memcached(golibmc)
package mcmutex

import (
	"errors"
	"github.com/douban/libmc/golibmc"
	"time"
)

const (
	defaultRetry      = 0
	defaultInterval   = 10 * time.Millisecond
	defaultExpiration = 30
)

// ErrLockFailed means failure to acquire lock after all retrys.
var ErrLockFailed = errors.New("failed to acquire lock")

// MCMutex can create mutex using golibmc
type MCMutex struct {
	client *golibmc.Client

	// retry interval
	Interval time.Duration

	// retry count before acquisition lock (default: 0)
	Retry int

	// lock will be expired after Expiration time (default: 30s)
	Expiration int64
}

// NewMCMutex create *MCMutex using default configure.
func NewMCMutex(mc *golibmc.Client) *MCMutex {
	return &MCMutex{
		client:     mc,
		Interval:   defaultInterval,
		Retry:      defaultRetry,
		Expiration: defaultExpiration,
	}
}

// Lock the key, or sleep and retry to lock according to configuration.
// it returns err when fail to acquire lock or got memcached error.
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

// Unlock the key.
// it returns err when lock is already free or got memcached error.
func (m *MCMutex) Unlock(key string) error {
	return m.client.Delete(key)
}
