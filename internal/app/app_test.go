package app

import (
	"context"
	"testing"
)

func TestMockModeNeedsNoDevice(t *testing.T) {
	d, err := New(Options{Mock: true})
	if err != nil || !d.Mock || d.Device == nil || d.Packages == nil || d.Accounts == nil || d.Txn == nil {
		t.Fatalf("%+v %v", d, err)
	}
}

func TestMockIsNotOpenUnlessFirstRun(t *testing.T) {
	closed, _ := New(Options{Mock: true})
	if set, _ := closed.Accounts.RootPasswordSet(context.Background()); !set {
		t.Fatal("plain mock mode must have a root password")
	}
	open, _ := New(Options{Mock: true, FirstRun: true})
	if set, _ := open.Accounts.RootPasswordSet(context.Background()); set {
		t.Fatal("--mock-firstrun must start without a password")
	}
}
