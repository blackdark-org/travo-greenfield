package app

import "testing"

func TestMockModeNeedsNoDevice(t *testing.T) {
	d, err := New(true)
	if err != nil || !d.Mock || d.Device == nil || d.Packages == nil {
		t.Fatalf("%+v %v", d, err)
	}
}
