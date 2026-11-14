package runner

import (
	"context"
	"errors"
	"os/exec"
	"testing"
)

func TestRunInputSeparatesStderr(t *testing.T) {
	_, code, err := RunInput(context.Background(), "", "sh", "-c", "printf 'first line\\nsecond line\\n' >&2; exit 7")
	if code != 7 || err == nil {
		t.Fatalf("exit code %d, error %v", code, err)
	}
	stderr, ok := Stderr(err)
	if !ok || stderr != "first line\nsecond line" {
		t.Fatalf("stderr %q, %v", stderr, ok)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("command cause lost: %v", err)
	}
}
