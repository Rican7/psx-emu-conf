// Copyright © Trevor N. Suarez (Rican7)

package data

import (
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/Rican7/psx-emu-conf/internal/data/normalize"
)

var (
	// regexSerialCode defines a regular expression for matching against valid
	// serial codes.
	//
	// A valid PlayStation software serial code is in the following format:
	//  - 4 letters, a dash (`-`), and then 5 digits
	//  - The 1st letters should always be an `S` for PlayStation 1 titles
	//  - The 2nd letter should be either a `C` or `L`
	//  - The 3rd letter should be one of A, C, E, K, P, or U
	//  - The 4th letter should be one of D, M, S, T, or X
	//
	// See: https://serialstation.com/serials/guide/
	//  (though, SerialStation seems to mix up "J" and "P"...)
	regexSerialCode = regexp.MustCompile(`^(S(C|L)(A|C|E|K|P|U)(D|M|S|T|X))-(\d{5})$`)
)

// Normalize modifies an app in-place by performing some normalizations on the
// data contained within the App.
func (a *App) Normalize() {
	title, titleVariations := normalize.Title(a.Title)

	titleVariationsSet := make(map[string]struct{})
	var normalizedTitleVariations []string
	for _, titleVariation := range append(a.TitleVariations, titleVariations...) {
		if _, ok := titleVariationsSet[titleVariation]; !ok {
			titleVariationsSet[titleVariation] = struct{}{}
			normalizedTitleVariations = append(normalizedTitleVariations, titleVariation)
		}
	}
	sort.Strings(normalizedTitleVariations)

	discNamesSet := make(map[string]struct{})
	var normalizedDiscNames []string
	for _, discName := range a.DiscNames {
		if _, ok := discNamesSet[discName]; !ok {
			discNamesSet[discName] = struct{}{}
			normalizedDiscNames = append(normalizedDiscNames, discName)
		}
	}
	sort.Strings(normalizedDiscNames)

	a.Region = Region(normalize.Region(string(a.Region)))
	a.SerialCode = normalize.SerialCode(a.SerialCode)
	a.Title = title
	a.TitleVariations = normalizedTitleVariations
	a.DiscNames = normalizedDiscNames
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

	if a.SerialCode != "" && !regexSerialCode.MatchString(a.SerialCode) {
		fmt.Println(a.SerialCode)
		return errors.New("invalid SerialCode")
	}

	if a.Title == "" {
		return errors.New("missing Title")
	}

	if a.FeatureSupport.AnalogSupport < AnalogSupportUnknown ||
		a.FeatureSupport.AnalogSupport > AnalogSupportRequired {
		return errors.New("invalid FeatureSupport.AnalogSupport level")
	}

	if a.FeatureSupport.RumbleSupport < RumbleSupportUnknown ||
		a.FeatureSupport.RumbleSupport > RumbleSupportYes {
		return errors.New("invalid FeatureSupport.RumbleSupport level")
	}

	return nil
}
