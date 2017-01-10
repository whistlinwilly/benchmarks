package main

import (
	"runtime"
	"sync"
	"testing"
)

const ROUTINES = 10

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkOneToOneChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(b.N)
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- struct{}{}
		}
	}()
	for i := 0; i < b.N; i++ {
		<-ch
		wg.Done()
	}
	wg.Wait()
}

func BenchmarkFanOutChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			<-ch
			wg.Done()
		}()
	}
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
	wg.Wait()
}

func BenchmarkFanInChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			ch <- struct{}{}
		}()
	}
	for i := 0; i < b.N; i++ {
		<-ch
		wg.Done()
	}
	wg.Wait()
}

func BenchmarkOneToOneBufferedChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 100)
	wg.Add(b.N)
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- struct{}{}
		}
	}()
	for i := 0; i < b.N; i++ {
		<-ch
		wg.Done()
	}
	wg.Wait()
}

func BenchmarkFanOutBufferedChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 100)
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			<-ch
			wg.Done()
		}()
	}
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
	wg.Wait()
}

func BenchmarkFanInBufferedChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 100)
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			ch <- struct{}{}
		}()
	}
	for i := 0; i < b.N; i++ {
		<-ch
		wg.Done()
	}
	wg.Wait()
}

func BenchmarkFanOutBufferedBatchedChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 100)
	wg.Add(b.N)
	n := b.N / 10
	if n < 1 {
		n = 1
	}
	for i := 0; i < n; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				<-ch
				wg.Done()
			}
		}()
	}
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
	wg.Wait()
}

func BenchmarkFanInBufferedBatchChan(b *testing.B) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 100)
	wg.Add(b.N)
	n := b.N / 10
	if n < 1 {
		n = 1
	}
	for i := 0; i < n; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				ch <- struct{}{}
			}
		}()
	}
	for i := 0; i < b.N; i++ {
		<-ch
		wg.Done()
	}
	wg.Wait()
}
