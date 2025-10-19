package main

import (
	"testing"
	"time"
)

func Test_or(t *testing.T) {
	t.Run("no channels", func(t *testing.T) {
		c := or()
		select {
		case <-c:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("should have been closed immediately")
		}
	})

	t.Run("one channel", func(t *testing.T) {
		ch1 := make(chan interface{})
		go func() {
			time.Sleep(50 * time.Millisecond)
			close(ch1)
		}()

		start := time.Now()
		<-or(ch1)

		if time.Since(start) < 50*time.Millisecond {
			t.Fatal("closed too early")
		}
	})

	t.Run("two channels", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})

		go func() {
			time.Sleep(50 * time.Millisecond)
			close(ch1)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			close(ch2)
		}()

		start := time.Now()
		<-or(ch1, ch2)

		if time.Since(start) > 75*time.Millisecond {
			t.Fatal("closed too late")
		}
	})

	t.Run("multiple channels, first one is fastest", func(t *testing.T) {
		sig := func(after time.Duration) <-chan interface{} {
			c := make(chan interface{})
			go func() {
				defer close(c)
				time.Sleep(after)
			}()
			return c
		}

		start := time.Now()

		<-or(
			sig(50*time.Millisecond),
			sig(100*time.Millisecond),
			sig(150*time.Millisecond),
		)

		if time.Since(start) > 75*time.Millisecond {
			t.Fatal("closed too late")
		}
	})

	t.Run("multiple channels, middle one is fastest", func(t *testing.T) {
		sig := func(after time.Duration) <-chan interface{} {
			c := make(chan interface{})
			go func() {
				defer close(c)
				time.Sleep(after)
			}()
			return c
		}

		start := time.Now()

		<-or(
			sig(100*time.Millisecond),
			sig(50*time.Millisecond),
			sig(150*time.Millisecond),
		)

		if time.Since(start) > 75*time.Millisecond {
			t.Fatal("closed too late")
		}
	})
}
