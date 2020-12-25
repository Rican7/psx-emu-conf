// Copyright © Trevor N. Suarez (Rican7)

// Package gdocechoj2 provides a data source for the "PSX - General Game Info"
// spreadsheet that was created by Reddit user "Echoj2" and hosted on Google
// Docs.
//
// Sources:
//  - https://docs.google.com/spreadsheets/d/1D4FKPOWCi11zhVvUcS8Bv4-IzyxH9MZRldugigTc59E
//  - https://np.reddit.com/r/RetroPie/comments/9ala88/ps1_games_that_dont_require_l2_r2_and_analog/e4wa8p3/?context=100
//  - https://www.reddit.com/user/Echoj2
package gdocechoj2

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Rican7/psx-emu-conf/internal/data"
	"github.com/Rican7/psx-emu-conf/internal/data/source"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	googleDocID = "1D4FKPOWCi11zhVvUcS8Bv4-IzyxH9MZRldugigTc59E"

	sheetNTSCU = "NTSC-U"
	sheetNTSCJ = "NTSC-J"
	sheetPAL   = "PAL"

	dataCellRange = "A2:G99999"

	digitalOnlyLevelNo         = 1
	digitalOnlyLevelYes        = 2
	digitalOnlyLevelAnalogOnly = 3
	digitalOnlyLevelUnknown    = 4

	analogLevelYes     = 1
	analogLevelNo      = 2
	analogLevelUnknown = 3

	vibrationLevelYes     = 1
	vibrationLevelNo      = 2
	vibrationLevelUnknown = 3
)

// A map of regions to data cell ranges.
var regionRangeMap = map[data.Region]string{
	data.RegionNTSCU: fmt.Sprintf("%s!%s", sheetNTSCU, dataCellRange),
	data.RegionNTSCJ: fmt.Sprintf("%s!%s", sheetNTSCJ, dataCellRange),
	data.RegionPAL:   fmt.Sprintf("%s!%s", sheetPAL, dataCellRange),
}

var (
	// A regex for capturing the serial code from a psxdatacenter.com URL.
	regexLinkSerial = regexp.MustCompile(`\/([A-Za-z0-9-]+)\.html`)

	// A regex for capturing the feature level of a feature column.
	regexFeatureLevel = regexp.MustCompile(`(\d+) - (.*)`)
)

type src struct {
	apiKey string
}

// New returns a Source.
func New(googleAPIKey string) source.Source {
	return &src{
		apiKey: googleAPIKey,
	}
}

func (s *src) Fetch(ctx context.Context) ([]data.App, error) {
	var apps []data.App

	service, err := sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsReadonlyScope), option.WithAPIKey(s.apiKey))
	if err != nil {
		return nil, err
	}

	for region, cellRange := range regionRangeMap {
		dataResp, err := service.Spreadsheets.Get(googleDocID).
			Ranges(cellRange).
			Fields("sheets(data(rowData(values(hyperlink))))").
			IncludeGridData(true).
			Do()

		valueResp, err := service.Spreadsheets.Values.Get(googleDocID, cellRange).Do()
		if err != nil {
			return nil, err
		}

		apps = append(apps, processResponse(region, valueResp, dataResp)...)
	}

	return apps, nil
}

func processResponse(region data.Region, valueResp *sheets.ValueRange, dataResp *sheets.Spreadsheet) []data.App {
	var apps []data.App

	for i, row := range valueResp.Values {
		var link string
		var serialCode string

		if len(dataResp.Sheets[0].Data[0].RowData[i].Values) > 0 {
			link = dataResp.Sheets[0].Data[0].RowData[i].Values[0].Hyperlink
		}

		if serialCodeMatches := regexLinkSerial.FindStringSubmatch(link); len(serialCodeMatches) > 1 {
			serialCode = serialCodeMatches[1]
		}

		var digitalOnlyLevel int
		var analogLevel int
		var vibrationLevel int

		if digitalOnlyLevelMatches := regexFeatureLevel.FindStringSubmatch(row[2].(string)); len(digitalOnlyLevelMatches) > 1 {
			digitalOnlyLevel, _ = strconv.Atoi(digitalOnlyLevelMatches[1])
		}

		if analogLevelMatches := regexFeatureLevel.FindStringSubmatch(row[2].(string)); len(analogLevelMatches) > 1 {
			analogLevel, _ = strconv.Atoi(analogLevelMatches[1])
		}

		if vibrationLevelMatches := regexFeatureLevel.FindStringSubmatch(row[2].(string)); len(vibrationLevelMatches) > 1 {
			vibrationLevel, _ = strconv.Atoi(vibrationLevelMatches[1])
		}

		var analogSupport data.AnalogSupport
		var rumbleSupport data.RumbleSupport

		// Determine analog support
		switch {
		case digitalOnlyLevel == digitalOnlyLevelYes, analogLevel == analogLevelNo:
			analogSupport = data.AnalogSupportNo
		case digitalOnlyLevel == digitalOnlyLevelNo, analogLevel == analogLevelYes:
			analogSupport = data.AnalogSupportYes
		case digitalOnlyLevel == digitalOnlyLevelAnalogOnly:
			analogSupport = data.AnalogSupportRequired
		case digitalOnlyLevel == digitalOnlyLevelUnknown:
			fallthrough
		case analogLevel == analogLevelUnknown:
			fallthrough
		default:
			analogSupport = data.AnalogSupportUnknown
		}

		// Determine vibration support
		switch {
		case vibrationLevel == vibrationLevelNo:
			rumbleSupport = data.RumbleSupportNo
		case vibrationLevel == vibrationLevelYes:
			rumbleSupport = data.RumbleSupportYes
		case vibrationLevel == vibrationLevelUnknown:
			fallthrough
		default:
			rumbleSupport = data.RumbleSupportUnknown
		}

		app := data.App{
			Region:     region,
			SerialCode: serialCode,
			Title:      row[0].(string),

			FeatureSupport: data.FeatureSupport{
				AnalogSupport: analogSupport,
				RumbleSupport: rumbleSupport,
			},
		}

		apps = append(apps, app)
	}

	return apps
}
