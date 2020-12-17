// Copyright © Trevor N. Suarez (Rican7)

package retroarch

import (
	"fmt"
	"io"

	"github.com/Rican7/psx-emu-conf/internal/data"
	"github.com/Rican7/psx-emu-conf/internal/emuconf"
)

const (
	// CoreNameBeetlePSXHW defines the proper name of the Beetle PSX HW emulator.
	CoreNameBeetlePSXHW = "Beetle PSX HW"

	beetlePSXHWInternalName = CoreNameBeetlePSXHW
)

const (
	beetlePSXHWConfigAnalogToggleKey = "beetle_psx_hw_analog_toggle"

	beetlePSXHWConfigAnalogToggleValueDisabled = `"disabled"`
	beetlePSXHWConfigAnalogToggleValueEnabled  = `"enabled"`
)

// beetlePSXHW represents the Beetle PSX HW emulator core in RetroArch.
//
// See:
//  - https://mednafen.github.io/documentation/psx.html
//  - https://github.com/libretro/beetle-psx-libretro
//  - https://docs.libretro.com/library/beetle_psx_hw/
type beetlePSXHW struct {
	*core
}

// NewBeetlePSXHW returns a Configurator for the Beetle PSX HW core in RetroArch.
func NewBeetlePSXHW() emuconf.Configurator {
	return &beetlePSXHW{
		core: &core{
			internalName: beetlePSXHWInternalName,
			displayName:  CoreNameBeetlePSXHW,
		},
	}
}

func (e *beetlePSXHW) Configure(writer io.Writer, app data.App) error {
	var err error

	var analogToggleValue string

	switch {
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportYes,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportRequired:
		analogToggleValue = beetlePSXHWConfigAnalogToggleValueEnabled
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportNo,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportUnknown:
		fallthrough
	default:
		analogToggleValue = beetlePSXHWConfigAnalogToggleValueDisabled
	}

	_, err = fmt.Fprintf(writer, "%s = %s\n", beetlePSXHWConfigAnalogToggleKey, analogToggleValue)

	return err
}
