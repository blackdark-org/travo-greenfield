package mock

import (
	"bytes"
	"context"
	"os"
	"testing"
)

// The embed cannot reach outside the package, so the copies must not drift from the captures.
func TestSeedCopiesMatchCaptures(t *testing.T) {
	for _, name := range []string{"system-board.json", "system-info.json"} {
		orig, err := os.ReadFile("../../../testdata/device/axt1800-25.12.6/" + name)
		if err != nil {
			t.Fatal(err)
		}
		cp, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(orig, cp) {
			t.Errorf("%s: mock copy differs from testdata/device/axt1800-25.12.6", name)
		}
	}
}

func TestSeedSystemAnswersUbus(t *testing.T) {
	d := New()
	SeedSystem(d)
	for _, m := range []string{"board", "info"} {
		if _, err := d.Call(context.Background(), "system", m, nil); err != nil {
			t.Errorf("system %s: %v", m, err)
		}
	}
}
