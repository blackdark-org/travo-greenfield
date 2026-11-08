package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/BlackDark/travo-greenfield/internal/device"
)

// Timeout tiers: every shell-out gets one.
const (
	TierRead  = 5 * time.Second
	TierWrite = 15 * time.Second
	TierApply = 60 * time.Second
)

// Runner executes a command. Production uses ExecRunner; tests inject a fake.
type Runner func(ctx context.Context, name string, args ...string) (stdout string, exit int, err error)

// ExecRunner runs a real process.
func ExecRunner(ctx context.Context, name string, args ...string) (string, int, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return out.String(), ee.ExitCode(), fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
	}
	return out.String(), 0, err
}

// Device is the real-router adapter.
type Device struct{ Run Runner }

func New() *Device { return &Device{Run: ExecRunner} }

func (d *Device) run(ctx context.Context, tier time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, tier)
	defer cancel()
	out, code, err := d.Run(ctx, name, args...)
	if err != nil && name == "uci" && code == 1 {
		return out, fmt.Errorf("%w: %v", device.ErrNotFound, err)
	}
	return out, err
}

func (d *Device) Show(ctx context.Context, config string) ([]device.Section, error) {
	out, err := d.run(ctx, TierRead, "uci", "show", config)
	if err != nil {
		return nil, err
	}
	return ParseShow(config, out)
}

func (d *Device) Set(ctx context.Context, config, section, option, value string) error {
	// the section must exist: uci would silently create nothing useful otherwise
	if _, err := d.run(ctx, TierRead, "uci", "-q", "get", config+"."+section); err != nil {
		return fmt.Errorf("%s.%s: %w", config, section, device.ErrNotFound)
	}
	_, err := d.run(ctx, TierWrite, "uci", "set", fmt.Sprintf("%s.%s.%s=%s", config, section, option, value))
	return err
}

func (d *Device) SetType(ctx context.Context, config, section, typ string) error {
	_, err := d.run(ctx, TierWrite, "uci", "set", fmt.Sprintf("%s.%s=%s", config, section, typ))
	return err
}

func (d *Device) AddToList(ctx context.Context, config, section, option, value string) error {
	_, err := d.run(ctx, TierWrite, "uci", "add_list", fmt.Sprintf("%s.%s.%s=%s", config, section, option, value))
	return err
}

func (d *Device) DeleteOption(ctx context.Context, config, section, option string) error {
	_, err := d.run(ctx, TierWrite, "uci", "-q", "delete", fmt.Sprintf("%s.%s.%s", config, section, option))
	return err
}

func (d *Device) DeleteSection(ctx context.Context, config, section string) error {
	_, err := d.run(ctx, TierWrite, "uci", "delete", fmt.Sprintf("%s.%s", config, section))
	return err
}

func (d *Device) Commit(ctx context.Context, config string) error {
	_, err := d.run(ctx, TierWrite, "uci", "commit", config)
	return err
}

func (d *Device) Revert(ctx context.Context, config string) error {
	_, err := d.run(ctx, TierWrite, "uci", "revert", config)
	return err
}

func (d *Device) Call(ctx context.Context, object, method string, args any) (json.RawMessage, error) {
	payload := []byte("{}")
	if args != nil {
		var err error
		if payload, err = json.Marshal(args); err != nil {
			return nil, err
		}
	}
	out, err := d.run(ctx, TierRead, "ubus", "call", object, method, string(payload))
	if err != nil {
		return nil, err
	}
	if !json.Valid([]byte(out)) {
		return nil, fmt.Errorf("ubus %s.%s: reply is not JSON", object, method)
	}
	return json.RawMessage(out), nil
}

// HasDefaultRoute reads the kernel table (netifd status can lie, see device-facts §1).
func (d *Device) HasDefaultRoute(ctx context.Context) (bool, error) {
	out, err := d.run(ctx, TierRead, "ip", "-4", "route", "show", "default")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

var _ device.Device = (*Device)(nil)
