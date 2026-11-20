package txn

import (
	"context"

	"github.com/blackdark-org/travo-greenfield/internal/device"
)

// Settled is a config read that no transaction touched while it ran.
type Settled struct {
	Sections []device.Section
	// Quiet is false when a transaction began, or its guard was held, at any point between
	// the start and the end of the read: the sections may belong to a change that is not
	// confirmed or was rolled back meanwhile.
	Quiet bool
}

// ShowSettled reads config under its write lock, so it never sees another writer's staged
// delta: staging happens under that lock and ends in commit or revert before it is released.
// It waits for a transaction's apply phase (up to ApplyBudget); use plain Show for UI reads.
func (m *Manager) ShowSettled(ctx context.Context, config string) (Settled, error) {
	before := m.begun.Load()
	heldBefore, err := m.guardHeld()
	if err != nil {
		return Settled{}, err
	}
	if m.afterFirstCheck != nil {
		m.afterFirstCheck()
	}
	lctx, unlock, err := m.Locks.Lock(ctx, []string{config})
	if err != nil {
		return Settled{}, err
	}
	sections, err := m.Device.Show(lctx, config)
	unlock()
	if err != nil {
		return Settled{}, err
	}
	heldAfter, err := m.guardHeld()
	if err != nil {
		return Settled{}, err
	}
	quiet := !heldBefore && !heldAfter && m.begun.Load() == before
	return Settled{Sections: sections, Quiet: quiet}, nil
}

func (m *Manager) guardHeld() (bool, error) {
	held, err := m.Guards.Stale()
	for _, g := range held {
		if g.Name == GuardName {
			return true, err
		}
	}
	return false, err
}
