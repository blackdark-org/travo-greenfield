package uciwrite

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestReentrantLockIsRefused(t *testing.T) {
	var l Locks
	ctx, unlock, err := l.Lock(context.Background(), []string{"network"})
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	if _, _, err := l.Lock(ctx, []string{"dhcp"}); !errors.Is(err, ErrReentrant) {
		t.Fatalf("got %v", err)
	}
}

func TestNormalizeRejectsInjectionAndSorts(t *testing.T) {
	for _, bad := range []string{"", "a.b", "../x", "a b", "x;rm"} {
		if _, err := Normalize([]string{bad}); !errors.Is(err, ErrBadConfig) {
			t.Errorf("accepted %q", bad)
		}
	}
	got, _ := Normalize([]string{"wireless", "dhcp", "wireless", "firewall"})
	if len(got) != 3 || got[0] != "dhcp" || got[2] != "wireless" {
		t.Fatal(got)
	}
}

func TestOverlappingSetsDoNotDeadlock(t *testing.T) {
	var l Locks
	var wg sync.WaitGroup
	sets := [][]string{{"network", "firewall"}, {"firewall", "network"}, {"dhcp", "network"}}
	for i := 0; i < 200; i++ {
		for _, s := range sets {
			wg.Add(1)
			go func(s []string) {
				defer wg.Done()
				_, unlock, err := l.Lock(context.Background(), s)
				if err != nil {
					t.Error(err)
					return
				}
				unlock()
			}(s)
		}
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("deadlock")
	}
}

func TestSameConfigIsExclusive(t *testing.T) {
	var l Locks
	_, unlock, _ := l.Lock(context.Background(), []string{"dhcp"})
	got := make(chan struct{})
	go func() {
		_, u, _ := l.Lock(context.Background(), []string{"dhcp"})
		close(got)
		u()
	}()
	select {
	case <-got:
		t.Fatal("second writer got the lock")
	case <-time.After(50 * time.Millisecond):
	}
	unlock()
	<-got
}
