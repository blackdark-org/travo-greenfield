package app

import (
	"context"
	"os"
	"testing"
)

func TestMockModeNeedsNoDevice(t *testing.T) {
	d, err := New(Options{Mock: true})
	t.Cleanup(func() { _ = d.Close() })
	if err != nil || !d.Mock || d.Device == nil || d.Packages == nil || d.Accounts == nil || d.Txn == nil {
		t.Fatalf("%+v %v", d, err)
	}
}

func TestMockIsNotOpenUnlessFirstRun(t *testing.T) {
	closed, _ := New(Options{Mock: true})
	t.Cleanup(func() { _ = closed.Close() })
	if set, _ := closed.Accounts.RootPasswordSet(context.Background()); !set {
		t.Fatal("plain mock mode must have a root password")
	}
	open, _ := New(Options{Mock: true, FirstRun: true})
	t.Cleanup(func() { _ = open.Close() })
	if set, _ := open.Accounts.RootPasswordSet(context.Background()); set {
		t.Fatal("--mock-firstrun must start without a password")
	}
}

func TestCloseRemovesMockStateDir(t *testing.T) {
	d, err := New(Options{Mock: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(d.Paths.State); err != nil {
		t.Fatal(err)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(d.Paths.State); !os.IsNotExist(err) {
		t.Fatalf("mock dir left behind: %v", err)
	}
}
