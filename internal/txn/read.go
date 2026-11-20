package txn

import (
	"context"

	"github.com/blackdark-org/travo-greenfield/internal/device"
)

// ShowSettled reads config under its write lock, so it never sees another writer's staged
// delta: staging happens under that lock and ends in commit or revert before it is released.
// It waits for a transaction's apply phase (up to ApplyBudget); use plain Show for UI reads.
func (m *Manager) ShowSettled(ctx context.Context, config string) ([]device.Section, error) {
	lctx, unlock, err := m.Locks.Lock(ctx, []string{config})
	if err != nil {
		return nil, err
	}
	defer unlock()
	return m.Device.Show(lctx, config)
}
