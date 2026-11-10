// Package guard implements crash-guard files for live-state operations.
// A present guard means "an operation started and did not finish": features skip.
package guard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/blackdark-org/travo-greenfield/internal/atomicfile"
)

// ErrActive means a guard of that name already exists.
var ErrActive = errors.New("guard already active")

var validName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*-in-progress$`)

// Dir manages guards inside one directory that holds nothing else.
type Dir struct {
	Path string
	Now  func() time.Time
}

// Info describes a guard found on disk.
type Info struct {
	Name    string    `json:"name"`
	Started time.Time `json:"started"`
}

func (d Dir) file(name string) (string, error) {
	if !validName.MatchString(name) {
		return "", fmt.Errorf("invalid guard name %q", name)
	}
	return filepath.Join(d.Path, name), nil
}

// Acquire writes the guard first. A failed write aborts the caller's operation.
func (d Dir) Acquire(name string) error {
	p, err := d.file(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(d.Path, 0o700); err != nil {
		return fmt.Errorf("guard dir: %w", err)
	}
	if _, err := os.Lstat(p); err == nil {
		return ErrActive
	}
	return atomicfile.WriteJSON(p, Info{Name: name, Started: d.Now().UTC()}, 0o600)
}

// Release removes the guard after the operation fully succeeded.
func (d Dir) Release(name string) error {
	p, err := d.file(name)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// Clear is the explicit retry path; it removes one guard by exact name, never by glob.
func (d Dir) Clear(name string) error { return d.Release(name) }

// Stale lists guards present on disk, for the startup report.
func (d Dir) Stale() ([]Info, error) {
	entries, err := os.ReadDir(d.Path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := []Info{}
	for _, e := range entries {
		if !validName.MatchString(e.Name()) {
			continue
		}
		var in Info
		if err := atomicfile.ReadJSON(filepath.Join(d.Path, e.Name()), &in); err != nil {
			in = Info{Name: e.Name()} // present but unreadable still counts as active
		}
		out = append(out, in)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
