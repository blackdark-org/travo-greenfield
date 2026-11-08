// Package app wires adapters into the daemon; it is the only place that picks real vs mock.
package app

import (
	"fmt"

	"github.com/BlackDark/travo-greenfield/internal/device"
	"github.com/BlackDark/travo-greenfield/internal/device/cli"
	"github.com/BlackDark/travo-greenfield/internal/device/mock"
	"github.com/BlackDark/travo-greenfield/internal/pkgmgr"
	"github.com/BlackDark/travo-greenfield/internal/runner"
)

// Deps are the ports every feature receives; features never import adapters.
type Deps struct {
	Device   device.Device
	Packages pkgmgr.PackageManager
	Mock     bool
}

// New builds real adapters, or mock adapters when useMock is set (develop without a device).
func New(useMock bool) (Deps, error) {
	if useMock {
		return Deps{Device: mock.New(), Packages: pkgmgr.NewMock(), Mock: true}, nil
	}
	pm, err := pkgmgr.DetectSystem(runner.Exec)
	if err != nil {
		return Deps{}, fmt.Errorf("package manager: %w", err)
	}
	return Deps{Device: cli.New(), Packages: pm}, nil
}
