// Copyright © Trevor N. Suarez (Rican7)

package retroarch

import (
	"fmt"
	"io"

	"github.com/Rican7/psx-emu-conf/internal/data"
	"github.com/Rican7/psx-emu-conf/internal/emuconf"
)

const (
	// CoreNameBeetlePSX defines the proper name of the Beetle PSX emulator.
	CoreNameBeetlePSX = "Beetle PSX"

	beetlePSXInternalName = CoreNameBeetlePSX
)

const (
	beetlePSXConfigAnalogToggleKey = "beetle_psx_analog_toggle"

	beetlePSXConfigAnalogToggleValueDisabled = `"disabled"`
	beetlePSXConfigAnalogToggleValueEnabled  = `"enabled"`
)

// beetlePSX represents the Beetle PSX emulator core in RetroArch.
//
// See:
//  - https://mednafen.github.io/documentation/psx.html
//  - https://github.com/libretro/beetle-psx-libretro
//  - https://docs.libretro.com/library/beetle_psx/
type beetlePSX struct {
	*core
}

// NewBeetlePSX returns a Configurator for the Beetle PSX core in RetroArch.
func NewBeetlePSX() emuconf.Configurator {
	return &beetlePSX{
		core: &core{
			internalName: beetlePSXInternalName,
			displayName:  CoreNameBeetlePSX,
		},
	}
}

func (e *beetlePSX) Configure(writer io.Writer, app data.App) error {
	var err error

	var analogToggleValue string

	switch {
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportYes,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportRequired:
		analogToggleValue = beetlePSXConfigAnalogToggleValueEnabled
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportNo,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportUnknown:
		fallthrough
	default:
		analogToggleValue = beetlePSXConfigAnalogToggleValueDisabled
	}

	_, err = fmt.Fprintf(writer, "%s = %s\n", beetlePSXConfigAnalogToggleKey, analogToggleValue)

	return err
}
