// Copyright © Trevor N. Suarez (Rican7)

// Package source defines mechanisms for interacting with data sources.
package source

import (
	"context"

	"github.com/Rican7/psx-emu-conf/internal/data"
)

// Source defines a common interface for data sources.
type Source interface {
	// Fetch retrieves a list of apps, or returns an error.
	Fetch(ctx context.Context) ([]data.App, error)
}
