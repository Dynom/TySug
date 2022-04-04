package rwc

import (
	"sync"
	"testing"
)

func TestReadLocksDontBlock(t *testing.T) {
	locks := []struct {
		name string
		lock RWCMutexer
	}{
		{name: "RWCMutex", lock: New()},
		{name: "RWMutex", lock: &sync.RWMutex{}},
	}

	for _, tt := range locks {
		t.Run(tt.name, func(t *testing.T) {
			const iterations = 1000

			wg := sync.WaitGroup{}
			wg.Add(iterations)

			for i := 0; i < iterations; i++ {
				go func() {
					// 1st
					tt.lock.RLock()

					// 2nd open
					tt.lock.RLock()
					tt.lock.RUnlock() // closing 2nd

					// 2nd open, again
					tt.lock.RLock()
					tt.lock.RUnlock() // closing 2nd

					tt.lock.RUnlock() // closing 1

					wg.Done()
				}()
			}

			wg.Wait()
		})
	}
}
