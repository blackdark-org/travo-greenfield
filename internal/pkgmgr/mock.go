package pkgmgr

import (
	"context"
	"sync"
)

// Mock is an in-memory PackageManager for mock mode and tests.
type Mock struct {
	mu   sync.Mutex
	pkgs map[string]string
}

func NewMock() *Mock {
	return &Mock{pkgs: map[string]string{"dnsmasq": "2.93-r1", "dropbear": "2025.89-r1"}}
}

func (*Mock) Name() string                 { return "mock" }
func (*Mock) Update(context.Context) error { return nil }

func (m *Mock) Install(_ context.Context, pkg string) error {
	if err := validPkg(pkg); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pkgs[pkg] = "1.0-r1"
	return nil
}

func (m *Mock) Remove(_ context.Context, pkg string) error {
	if err := validPkg(pkg); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.pkgs, pkg)
	return nil
}

func (m *Mock) Installed(context.Context) ([]Package, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Package, 0, len(m.pkgs))
	for n, v := range m.pkgs {
		out = append(out, Package{Name: n, Version: v})
	}
	return out, nil
}

func (m *Mock) Version(ctx context.Context, pkg string) (string, error) {
	pkgs, _ := m.Installed(ctx)
	return find(pkgs, pkg)
}

var _ PackageManager = (*Mock)(nil)
var _ PackageManager = (*APK)(nil)
var _ PackageManager = (*OPKG)(nil)
