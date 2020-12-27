// Copyright © Trevor N. Suarez (Rican7)

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Rican7/psx-emu-conf/internal/data"
	"github.com/Rican7/psx-emu-conf/internal/emuconf"
	"github.com/Rican7/psx-emu-conf/internal/emuconf/retroarch"
)

const (
	defaultPathToData        = "_data/data.json"
	defaultPathToConfigFiles = "_configs"
)

// TODO:
//
//  - Abstract and organize a bit
//  - Handle errors far better... os.Stderr?
//  - Add flags/options
//  - Get data path from flags
//  - Allow for piping data via stdin
//  - Get output path from flags
func main() {
	dataStream, err := os.Open(defaultPathToData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer dataStream.Close()

	var apps []data.App

	decoder := json.NewDecoder(dataStream)
	if err := decoder.Decode(&apps); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Close the open stream now that its been decoded, to save on file handles
	// (we're going to need ALL the file handles that we can get... 😅)
	dataStream.Close()

	for i := range apps {
		apps[i].Normalize()

		if err := apps[i].Validate(); err != nil {
			fmt.Println(err)
			continue
		}
	}

	configurators := []emuconf.Configurator{
		retroarch.NewPCSXReARMed(),
		retroarch.NewBeetlePSX(),
		retroarch.NewBeetlePSXHW(),
	}

	for _, app := range apps {
		for _, configurator := range configurators {
			var configFilePaths []string

			mainConfigFilePath, err := buildConfigPath(app, configurator)
			if err != nil {
				fmt.Println(err)
				continue
			}

			configFilePaths = append(configFilePaths, mainConfigFilePath)

			if altLocator, ok := configurator.(emuconf.AlternativesLocator); ok {
				altConfigFilePaths := buildAltConfigPaths(app, altLocator)

				configFilePaths = append(configFilePaths, altConfigFilePaths...)
			}

			for _, configFilePath := range configFilePaths {
				configFilePath = filepath.Join(defaultPathToConfigFiles, configFilePath)
				configFileDir := filepath.Dir(configFilePath)

				err = os.MkdirAll(configFileDir, 0777)
				if err != nil {
					fmt.Println(err)
					continue
				}

				file, err := os.Create(configFilePath)
				if err != nil {
					fmt.Println(err)
					continue
				}

				defer file.Close()

				err = configurator.Configure(file, app)
				if err != nil {
					fmt.Println(err)
				}

				file.Close()
			}
		}
	}
}

func buildConfigPath(app data.App, configurator emuconf.Configurator) (string, error) {
	var confPath string

	if locator, ok := configurator.(emuconf.Locator); ok {
		confPath = locator.Path(app)
	}

	if confPath == "" {
		switch {
		case app.SerialCode != "":
			confPath = app.SerialCode
		case app.Title != "":
			confPath = app.Title
		}
	}

	confPath = filepath.Clean(confPath)

	base := filepath.Base(confPath)
	if base == "" || base == "." || base == string(filepath.Separator) {
		return "", errors.New("incomplete file path")
	}

	return confPath, nil
}

func buildAltConfigPaths(app data.App, altLocator emuconf.AlternativesLocator) []string {
	var altConfPaths []string

	for _, altConfPath := range altLocator.AlternativePaths(app) {
		altConfPath = filepath.Clean(altConfPath)

		base := filepath.Base(altConfPath)
		if base == "" || base == "." || base == string(filepath.Separator) {
			// Output the error, but ignore this path, as its invalid
			fmt.Println("incomplete alt file path")
		}

		altConfPaths = append(altConfPaths, altConfPath)
	}

	return altConfPaths
}
