// Copyright © Trevor N. Suarez (Rican7)

// Package data defines structures and mechanisms for PlayStation data.
package data

import (
	"errors"

	"github.com/Rican7/psx-emu-conf/internal/data/normalize"
)

// App defines the structure of a PlayStation software title.
//
// These are called "App" rather than "Game", to support non-game releases
// without discrimination.
type App struct {
	Region          Region   `json:",omitempty"`
	SerialCode      string   `json:",omitempty"`
	Title           string   `json:",omitempty"`
	TitleVariations []string `json:",omitempty"`

	NumberOfDiscs uint     `json:",omitempty"`
	DiscNames     []string `json:",omitempty"`

	FeatureSupport FeatureSupport `json:",omitempty"`
}

// Region defines a "Region" of a PlayStation software title.
type Region string

// Available regions.
const (
	RegionNTSCU Region = "NTSC-U"
	RegionNTSCJ Region = "NTSC-J"
	RegionPAL   Region = "PAL"
)

// FeatureSupport defines a structure that represents the PlayStation features
// and peripherals support matrix.
type FeatureSupport struct {
	AnalogSupport AnalogSupport
	RumbleSupport RumbleSupport
}

// AnalogSupport defines the level of support of an "Analog" controller.
type AnalogSupport uint

// Analog support levels.
const (
	AnalogSupportUnknown  AnalogSupport = iota // Unknown level of support.
	AnalogSupportNo                            // No support. Digital only.
	AnalogSupportYes                           // Supports analog controllers.
	AnalogSupportRequired                      // Analog controller is required.
)

// RumbleSupport defines the level of support of the "rumble" feature.
type RumbleSupport uint

// Rumble support levels.
const (
	RumbleSupportUnknown RumbleSupport = iota // Unknown level of support.
	RumbleSupportNo                           // No support.
	RumbleSupportYes                          // Supports rumble.
)

// Normalize modifies an app in-place by performing some normalizations on the
// data contained within the App.
func (a *App) Normalize() {
	title, titleVariations := normalize.Title(a.Title)

	var titleVariationsSet map[string]struct{}
	var normalizedTitleVariations []string
	for _, titleVariation := range append(a.TitleVariations, titleVariations...) {
		if _, ok := titleVariationsSet[titleVariation]; !ok {
			titleVariationsSet[titleVariation] = struct{}{}
			normalizedTitleVariations = append(normalizedTitleVariations, titleVariation)
		}
	}

	a.Region = Region(normalize.Region(string(a.Region)))
	a.SerialCode = normalize.SerialCode(a.SerialCode)
	a.Title = title
	a.TitleVariations = normalizedTitleVariations
}

// Validate performs validations and returns an error if the data isn't valid.
// The returned error should describe what is invalid about the data.
func (a *App) Validate() error {
	switch a.Region {
	case RegionNTSCU, RegionNTSCJ, RegionPAL:
		// Valid
	default:
		return errors.New("invalid Region")
	}

	if a.Title == "" {
		return errors.New("missing Title")
	}

	return nil
}
