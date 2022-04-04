package rwc

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// All tests here are more or less copied over from rwmutex_test.go from the stdlib

// tlock is used as reference benchmark test case
type tlock struct{ lock sync.Mutex }

func (t *tlock) RLock()   { t.lock.Lock() }
func (t *tlock) RUnlock() { t.lock.Unlock() }
func (t *tlock) Lock()    { t.lock.Lock() }
func (t *tlock) Unlock()  { t.lock.Unlock() }

func parallelReader(m RWCMutexer, locked, unlocked, done chan bool) {
	m.RLock()
	locked <- true
	<-unlocked
	m.RUnlock()
	done <- true
}

func doTestParallelReaders(numReaders, gomaxprocs int) {
	runtime.GOMAXPROCS(gomaxprocs)
	var l = New()
	locked := make(chan bool)
	unlocked := make(chan bool)
	done := make(chan bool)
	for i := 0; i < numReaders; i++ {
		go parallelReader(l, locked, unlocked, done)
	}
	// Wait for all parallel RLock()s to succeed.
	for i := 0; i < numReaders; i++ {
		<-locked
	}
	for i := 0; i < numReaders; i++ {
		unlocked <- true
	}
	// Wait for the goroutines to finish.
	for i := 0; i < numReaders; i++ {
		<-done
	}
}

func TestParallelReaders(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(-1))
	doTestParallelReaders(1, 4)
	doTestParallelReaders(3, 4)
	doTestParallelReaders(4, 2)
}

func reader(l RWCMutexer, iterations int, activity *int32, done chan bool) {
	for i := 0; i < iterations; i++ {
		l.RLock()
		n := atomic.AddInt32(activity, 1)
		if n < 1 || n >= 10000 {
			l.RUnlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for i := 0; i < 100; i++ {
		}
		atomic.AddInt32(activity, -1)
		l.RUnlock()
	}
	done <- true
}

func writer(l RWCMutexer, iterations int, activity *int32, done chan bool) {
	for i := 0; i < iterations; i++ {
		l.Lock()
		n := atomic.AddInt32(activity, 10000)
		if n != 10000 {
			l.Unlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for i := 0; i < 100; i++ {
		}
		atomic.AddInt32(activity, -10000)
		l.Unlock()
	}
	done <- true
}

func HammerRWCMutex(gomaxprocs, numReaders, iterations int) {
	runtime.GOMAXPROCS(gomaxprocs)
	// Number of active readers + 10000 * number of active writers.
	var activity int32
	var l = New()
	done := make(chan bool)
	go writer(l, iterations, &activity, done)
	var i int
	for i = 0; i < numReaders/2; i++ {
		go reader(l, iterations, &activity, done)
	}
	go writer(l, iterations, &activity, done)
	for ; i < numReaders; i++ {
		go reader(l, iterations, &activity, done)
	}
	// Wait for the 2 writers and all readers to finish.
	for i := 0; i < 2+numReaders; i++ {
		<-done
	}
}

func TestLocker(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(-1))
	n := 1000
	if testing.Short() {
		n = 5
	}
	HammerRWCMutex(1, 1, n)
	HammerRWCMutex(1, 3, n)
	HammerRWCMutex(1, 10, n)
	HammerRWCMutex(4, 1, n)
	HammerRWCMutex(4, 3, n)
	HammerRWCMutex(4, 10, n)
	HammerRWCMutex(10, 1, n)
	HammerRWCMutex(10, 3, n)
	HammerRWCMutex(10, 10, n)
	HammerRWCMutex(10, 5, n)
}

func BenchmarkUncontended(b *testing.B) {
	const parallelism = 31

	b.Run("RWCMutex", func(b *testing.B) {
		b.SetParallelism(parallelism)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			l := New()

			for pb.Next() {
				l.RLock()
				l.RLock()

				l.RUnlock()
				l.RUnlock()

				l.Lock()
				l.Unlock()
			}
		})
	})

	b.Run("RWMutex", func(b *testing.B) {
		b.SetParallelism(parallelism)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			l := &sync.RWMutex{}

			for pb.Next() {
				l.RLock()
				l.RLock()

				l.RUnlock()
				l.RUnlock()

				l.Lock()
				l.Unlock()
			}
		})
	})
}

func BenchmarkWorkload(b *testing.B) {
	load := []struct {
		localWork  int // Work emulation
		writeRatio int // Higher the number, fewer the writes
	}{
		{localWork: 0, writeRatio: 100},
		{localWork: 0, writeRatio: 10},
		{localWork: 100, writeRatio: 100},
		{localWork: 100, writeRatio: 10},
		{localWork: 1000, writeRatio: 100},
		{localWork: 1000, writeRatio: 10},
	}

	locks := []struct {
		name string
		lock RWCMutexer
	}{
		{name: "RWCMutex", lock: New()},
		{name: "RWMutex", lock: &sync.RWMutex{}},
		{name: "tlock", lock: &tlock{}},
	}

	for _, tt := range load {
		for _, l := range locks {
			b.Run(fmt.Sprintf("%s/l%d/r%d", l.name, tt.localWork, tt.writeRatio), func(b *testing.B) {
				benchmarkRWLock(b, l.lock, tt.localWork, tt.writeRatio)
			})
		}
	}
}

func benchmarkRWLock(b *testing.B, l RWCMutexer, localWork, writeRatio int) {
	// From rwmutex_test.go
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			foo++
			if foo%writeRatio == 0 {
				l.Lock()
				l.Unlock()
			} else {
				l.RLock()
				for i := 0; i != localWork; i += 1 {
					foo *= 2
					foo /= 2
				}
				l.RUnlock()
			}
		}
		_ = foo
	})
}
