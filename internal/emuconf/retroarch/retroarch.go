// Copyright © Trevor N. Suarez (Rican7)

// Package retroarch provides emulator configurators for the RetroArch frontend.
//
// RetroArch is a popular frontend for numerous emulator "cores", that each
// allow for their own specific configuration, with the ability to provide
// global configurations, per-core configurations, per-directory configurations,
// and even per-game configurations.
//
// See:
//  - https://www.retroarch.com/
//  - https://docs.libretro.com/guides/overrides/
package retroarch

import (
	"fmt"
	"path"

	"github.com/Rican7/psx-emu-conf/internal/data"
)

const (
	// Name defines the proper name of the RetroArch frontend.
	Name = "RetroArch"

	// PathConfigDirectory defines a path to the directory where configurations
	// are stored within RetroArch.
	//
	// NOTE: This is the DEFAULT value, and is relative to the path where
	// RetroArch stores its default directories. As RetroArch allows for the
	// specification of directory paths, this could differ per user,
	// installation, or configuration.
	PathConfigDirectory = "config"

	// ExtensionPerGameCoreOption defines the file extension used for per-game
	// core option files.
	ExtensionPerGameCoreOption = ".opt"
)

type core struct {
	name string
}

func (c *core) EmulatorName() string {
	return fmt.Sprintf("%s - %s", Name, c.name)
}

func (c *core) Path(app data.App) string {
	return pathForGameCoreOptionFile(c.name, app)
}

func pathForGameCoreOptionFile(coreName string, app data.App) string {
	return path.Join(coreName, app.Title+ExtensionPerGameCoreOption)
}
