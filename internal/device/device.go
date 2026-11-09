// Package device is the only way the daemon touches OpenWrt: UCI, ubus and the kernel route table.
// Adapters: cli (real router) and mock (mock mode and tests).
package device

import (
	"context"
	"encoding/json"
	"errors"
)

// Section is one UCI section. Options hold strings, lists hold list options.
type Section struct {
	Name    string              `json:"name"`
	Type    string              `json:"type"`
	Options map[string]string   `json:"options"`
	Lists   map[string][]string `json:"lists"`
}

// ErrNotFound is returned when a config, section or ubus object does not exist.
var ErrNotFound = errors.New("not found")

// UCI mirrors `uci` semantics. Set is create-or-update (upsert) like the real tool.
// Staged changes are process-global per config until Commit or Revert.
type UCI interface {
	Show(ctx context.Context, config string) ([]Section, error)
	// Set writes one option; the section must exist (use AddSection first).
	Set(ctx context.Context, config, section, option, value string) error
	// SetType creates the section if missing, or updates its type.
	SetType(ctx context.Context, config, section, typ string) error
	AddToList(ctx context.Context, config, section, option, value string) error
	DeleteOption(ctx context.Context, config, section, option string) error
	DeleteSection(ctx context.Context, config, section string) error
	Commit(ctx context.Context, config string) error
	Revert(ctx context.Context, config string) error
	// Pending reports whether config has staged, uncommitted changes (`uci changes`).
	Pending(ctx context.Context, config string) (bool, error)
}

// Ubus calls an object method and returns the raw JSON reply.
type Ubus interface {
	Call(ctx context.Context, object, method string, args any) (json.RawMessage, error)
}

// Routes reads the KERNEL table. netifd status is not a source for route truth.
type Routes interface {
	HasDefaultRoute(ctx context.Context) (bool, error)
}

// Reloader applies committed config to running services: `reload_config`
// (procd ucitrack). Never `wifi` (device-facts §3).
type Reloader interface {
	Reload(ctx context.Context) error
}

// Device is the full port.
type Device interface {
	UCI
	Ubus
	Routes
	Reloader
}
