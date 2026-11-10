package mock

import (
	_ "embed"
)

//go:embed testdata/system-board.json
var systemBoard []byte

//go:embed testdata/system-info.json
var systemInfo []byte

// SeedSystem scripts the system ubus replies with the verbatim AXT1800 captures.
func SeedSystem(d *Device) {
	d.SetUbusReply("system", "board", systemBoard)
	d.SetUbusReply("system", "info", systemInfo)
}
