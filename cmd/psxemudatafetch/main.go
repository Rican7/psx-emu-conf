// Copyright © Trevor N. Suarez (Rican7)

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	encoder.Encode(apps)
}
