// Copyright © Trevor N. Suarez (Rican7)

// Package data defines structures and mechanisms for PlayStation data.
package data

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
