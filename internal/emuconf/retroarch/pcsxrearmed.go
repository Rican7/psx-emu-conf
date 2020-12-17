// Copyright © Trevor N. Suarez (Rican7)

package retroarch

import (
	"fmt"
	"io"

	"github.com/Rican7/psx-emu-conf/internal/data"
	"github.com/Rican7/psx-emu-conf/internal/emuconf"
)

const (
	// CoreNamePCSXReARMed defines the proper name of the PCSX-ReARMed emulator.
	CoreNamePCSXReARMed = "PCSX-ReARMed"
)

const (
	fmtConfigControllerTypeKey = "pcsx_rearmed_pad%dtype"

	configControllerTypeValueStandard  = `"standard"`
	configControllerTypeValueAnalog    = `"analog"`
	configControllerTypeValueDualShock = `"dualshock"`
)

type pcsxReARMed struct {
	*core
}

// NewPCSXReARMed returns a Configurator for the PCSX-ReARMed core in RetroArch.
func NewPCSXReARMed() emuconf.Configurator {
	return &pcsxReARMed{
		core: &core{
			name: CoreNamePCSXReARMed,
		},
	}
}

func (e *pcsxReARMed) Configure(writer io.Writer, app data.App) error {
	var err error

	var controllerTypeValue string

	switch {
	case app.FeatureSupport.RumbleSupport == data.RumbleSupportYes:
		controllerTypeValue = configControllerTypeValueDualShock
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportYes,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportRequired:
		controllerTypeValue = configControllerTypeValueAnalog
	case app.FeatureSupport.AnalogSupport == data.AnalogSupportNo,
		app.FeatureSupport.AnalogSupport == data.AnalogSupportUnknown:
		fallthrough
	default:
		controllerTypeValue = configControllerTypeValueStandard
	}

	// Write a value for each controller.
	//
	// TODO: This currently writes configs for two controllers. Support more
	// than that? How many? A configurable amount?
	for i := 1; i <= 2; i++ {
		key := fmt.Sprintf(fmtConfigControllerTypeKey, i)

		_, err = fmt.Fprintf(writer, "%s = %s\n", key, controllerTypeValue)
	}

	return err
}
