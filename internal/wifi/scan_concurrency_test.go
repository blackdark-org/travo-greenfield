package wifi

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

type scanCall struct {
	result ScanResult
	err    error
}

func TestConcurrentScansJoinOneFlightAndNextCallRescans(t *testing.T) {
	svc, scanner := connectedScanService(t)
	release := make(chan struct{})
	scanner.blocks = map[string]chan struct{}{"phy0-sta0": release}
	scanner.started = make(chan string, 4)
	calls := make(chan scanCall, 3)
	go func() {
		result, err := svc.Scan(context.Background())
		calls <- scanCall{result: result, err: err}
	}()
	if name := <-scanner.started; name != "phy0-sta0" {
		t.Fatalf("first target %s", name)
	}
	for range 2 {
		go func() {
			result, err := svc.Scan(context.Background())
			calls <- scanCall{result: result, err: err}
		}()
	}
	// Keep the first radio blocked while both additional callers reach the flight.
	time.Sleep(20 * time.Millisecond)
	close(release)
	var first ScanResult
	for range 3 {
		call := <-calls
		if call.err != nil {
			t.Fatal(call.err)
		}
		if first.Networks == nil {
			first = call.result
		} else if !reflect.DeepEqual(call.result, first) {
			t.Fatalf("joined result differs:\n%+v\n%+v", call.result, first)
		}
	}
	if want := []string{"phy0-sta0", "phy1-ap0"}; !reflect.DeepEqual(scanner.called(), want) {
		t.Fatalf("calls after joined scan: %v", scanner.called())
	}
	if _, err := svc.Scan(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"phy0-sta0", "phy1-ap0", "phy0-sta0", "phy1-ap0"}; !reflect.DeepEqual(scanner.called(), want) {
		t.Fatalf("next call did not rescan: %v", scanner.called())
	}
}

func TestScanCancelledWaiterDoesNotCancelFlight(t *testing.T) {
	svc, scanner := connectedScanService(t)
	release := make(chan struct{})
	scanner.blocks = map[string]chan struct{}{"phy0-sta0": release}
	scanner.started = make(chan string, 2)
	first := make(chan scanCall, 1)
	go func() {
		result, err := svc.Scan(context.Background())
		first <- scanCall{result: result, err: err}
	}()
	<-scanner.started
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Scan(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waiter: %v", err)
	}
	close(release)
	call := <-first
	if call.err != nil || len(call.result.Networks) == 0 {
		t.Fatalf("flight was cancelled: %+v", call)
	}
}
