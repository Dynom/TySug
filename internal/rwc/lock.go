package rwc

import (
	"sync"
	"sync/atomic"
)

const maxReaders int32 = 1 << 30

// RWCMutexer is mostly just a copy & paste from sync.RWMutex but with channels. It outperforms sync.RWMutex
// on workloads with many short-read locks scenarios versus small wlock. See benchmarks for context
type RWCMutexer interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
}

func New() *RWCMutex {
	return &RWCMutex{
		lock:      sync.Mutex{},
		writerSem: make(chan struct{}),
		readerSem: make(chan struct{}),
	}
}

type RWCMutex struct {
	lock        sync.Mutex
	readerCount int32 // Readers locked
	readerWait  int32 // Readers locked, waiting to depart

	writerSem chan struct{}
	readerSem chan struct{}
}

// Lock locks write workloads
func (l *RWCMutex) Lock() {
	// First, resolve competition with other writers.
	l.lock.Lock()

	// Announce to readers there is a pending writer.
	r := atomic.AddInt32(&l.readerCount, -maxReaders) + maxReaders

	// RLock for active readers.
	if r != 0 && atomic.AddInt32(&l.readerWait, r) != 0 {
		<-l.writerSem
	}
}

// Unlock releases locked write workloads
func (l *RWCMutex) Unlock() {
	// Announce to readers there is no active writer.
	r := atomic.AddInt32(&l.readerCount, maxReaders)
	if r >= maxReaders {
		panic("Unlock called during an unreleased reader")
	}

	// Unblock blocked readers, if [there are] any.
	for i := 0; i < int(r); i++ {
		l.readerSem <- struct{}{}
	}

	// Allow other writers to proceed.
	l.lock.Unlock()
}

// RLock allows read workloads to wait for writers to complete.
func (l *RWCMutex) RLock() {
	if atomic.AddInt32(&l.readerCount, 1) < 0 {
		// A writer is pending, wait for it.
		<-l.readerSem
	}
}

// RUnlock signals that a read workload is completed, so that a writer can start
func (l *RWCMutex) RUnlock() {
	r := atomic.AddInt32(&l.readerCount, -1)
	if r >= 0 {
		return
	}

	// A writer is pending.
	if atomic.AddInt32(&l.readerWait, -1) == 0 {
		// The last reader unblocks the writer.
		l.writerSem <- struct{}{}
	} else if r+1 == 0 || r+1 == -maxReaders {
		panic("Done called on un-Waiting")
	}
}
