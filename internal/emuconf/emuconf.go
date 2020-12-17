// Copyright © Trevor N. Suarez (Rican7)

// Package emuconf provides mechanisms to configure PlayStation emulators.
package emuconf

import (
	"io"

	"github.com/Rican7/psx-emu-conf/internal/data"
)

// Locator defines an interface for emulation configurators that can determine a
// path for a given app.
type Locator interface {
	Path(app data.App) string
}

// Configurator defines a common interface for emulation configurators that are
// capable of configuring an emulator for a given app.
type Configurator interface {
	EmulatorName() string
	Configure(writer io.Writer, app data.App) error
}
