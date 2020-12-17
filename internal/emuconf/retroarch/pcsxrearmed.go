// Copyright © Trevor N. Suarez (Rican7)

package retroarch

import (
	"fmt"
	"io"

	"github.com/Rican7/psx-emu-conf/internal/data"
	"github.com/Rican7/psx-emu-conf/internal/emuconf"
)

const (
	// CoreNamePCSXReARMed defines the proper name of the PCSX ReARMed emulator.
	CoreNamePCSXReARMed = "PCSX ReARMed"

	pcsxReARMedInternalName = "PCSX-ReARMed"
)

const (
	pcsxReARMedConfigControllerTypeKeyFormat = "pcsx_rearmed_pad%dtype"

	pcsxReARMedConfigControllerTypeValueStandard  = `"standard"`
	pcsxReARMedConfigControllerTypeValueAnalog    = `"analog"`
	pcsxReARMedConfigControllerTypeValueDualShock = `"dualshock"`
)

// pcsxReARMed emulator core in RetroArch.
//
// See:
//  - https://github.com/notaz/pcsx_rearmed
//  - https://github.com/libretro/pcsx_rearmed
//  - https://docs.libretro.com/library/pcsx_rearmed/
type pcsxReARMed struct {
	*core
}

// NewPCSXReARMed returns a Configurator for the PCSX ReARMed core in RetroArch.
func NewPCSXReARMed() emuconf.Configurator {
	return &pcsxReARMed{
		core: &core{
			internalName: pcsxReARMedInternalName,
			displayName:  CoreNamePCSXReARMed,
		},
	}
}

func (e *pcsxReARMed) Configure(writer io.Writer, app data.App) error {
	var err error

	var controllerTypeValue string

	switch {
	case app.FeatureSupport.RumbleSupport == data.RumbleSupportYes:
		controllerTypeValue = pcsxReARMedConfigControllerTypeValueDualShock
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportYes,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportRequired:
		controllerTypeValue = pcsxReARMedConfigControllerTypeValueAnalog
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportNo,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportUnknown:
		fallthrough
	default:
		controllerTypeValue = pcsxReARMedConfigControllerTypeValueStandard
	}

	// Write a value for each controller.
	//
	// TODO: This currently writes configs for two controllers. Support more
	// than that? How many? A configurable amount?
	for i := 1; i <= 2; i++ {
		key := fmt.Sprintf(pcsxReARMedConfigControllerTypeKeyFormat, i)

		_, err = fmt.Fprintf(writer, "%s = %s\n", key, controllerTypeValue)
	}

	return err
}
