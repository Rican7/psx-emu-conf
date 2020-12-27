// Copyright © Trevor N. Suarez (Rican7)

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/Rican7/psx-emu-conf/internal/data"
	"github.com/Rican7/psx-emu-conf/internal/data/source/gdocechoj2"
)

// TODO:
//
//  - Abstract and organize a bit
//  - Handle errors far better... os.Stderr?
//  - Fetch from multiple sources concurrently
//  - Merge the results into a combined result-set
//  - Potentially report any differences between sources?
func main() {
	ctx := context.Background()

	src := gdocechoj2.New(os.Getenv("GOOGLE_API_KEY"))

	apps, err := src.Fetch(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := range apps {
		apps[i].Normalize()

		if err := apps[i].Validate(); err != nil {
			fmt.Println(err)
			continue
		}
	}

	apps = mergeAppCollections(apps)

	sort.Sort(data.AppsDefault(apps))

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	encoder.Encode(apps)
}

// mergeAppCollections merges multiple collections of apps into one large
// collection of apps, merging matching app references into each other, merging
// left-to-right (first-to-last) and taking the left-most (first) collection as
// the source of truth if differences are found, and only merging other data if
// the data doesn't already exist in the app data being merged into.
func mergeAppCollections(appCollections ...[]data.App) []data.App {
	var mergedApps []data.App

	appsBySerialCode := make(map[string]data.App)
	var appsWithoutIndex []data.App

	for _, appCollection := range appCollections {
		for _, app := range appCollection {
			if app.SerialCode != "" {
				appBySerialCode, ok := appsBySerialCode[app.SerialCode]

				if ok {
					app = mergeApps(appBySerialCode, app)
				}

				appsBySerialCode[app.SerialCode] = app
			} else {
				appsWithoutIndex = append(appsWithoutIndex, app)
			}
		}
	}

	for _, app := range appsBySerialCode {
		mergedApps = append(mergedApps, app)
	}
	for _, app := range appsWithoutIndex {
		mergedApps = append(mergedApps, app)
	}

	return mergedApps
}

// mergeApps merges two apps together, treating the left (first) passed app as
// the "primary" and the right (second) passed app as the secondary. The data is
// merged into the primary app only if the data doesn't already exist in the
// primary app.
//
// When disagreements are present between merging apps, the primary app's data
// always wins.
//
// TODO: Handle both having data but disagreeing?
func mergeApps(appPrimary data.App, appSecondary data.App) data.App {
	app := appPrimary

	if app.Region == "" && appSecondary.Region != "" {
		app.Region = appSecondary.Region
	}

	if app.SerialCode == "" && appSecondary.SerialCode != "" {
		app.SerialCode = appSecondary.SerialCode
	}

	if app.Title == "" && appSecondary.Title != "" {
		app.Title = appSecondary.Title
	}

	if app.Title != appSecondary.Title {
		// Add the other title as a variation
		app.TitleVariations = append(app.TitleVariations, appSecondary.Title)
	}

	if app.NumberOfDiscs == 0 && appSecondary.NumberOfDiscs != 0 {
		app.NumberOfDiscs = appSecondary.NumberOfDiscs
	}

	if app.FeatureSupport.AnalogSupport == data.AnalogSupportUnknown &&
		appSecondary.FeatureSupport.AnalogSupport != data.AnalogSupportUnknown {
		app.FeatureSupport.AnalogSupport = appSecondary.FeatureSupport.AnalogSupport
	}

	if app.FeatureSupport.RumbleSupport == data.RumbleSupportUnknown &&
		appSecondary.FeatureSupport.RumbleSupport != data.RumbleSupportUnknown {
		app.FeatureSupport.RumbleSupport = appSecondary.FeatureSupport.RumbleSupport
	}

	app.TitleVariations = append(app.TitleVariations, appSecondary.TitleVariations...)
	app.DiscNames = append(app.DiscNames, appSecondary.DiscNames...)

	// Normalize the app data before returning
	app.Normalize()

	return app
}
