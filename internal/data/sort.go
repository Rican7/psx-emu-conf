// Copyright © Trevor N. Suarez (Rican7)

package data

import "sort"

type appsDefault []App

// AppsDefault returns a sort.Interface for []App, using a default sort order.
func AppsDefault(apps []App) sort.Interface {
	return appsDefault(apps)
}

func (a appsDefault) Len() int {
	return len(a)
}

func (a appsDefault) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a appsDefault) Less(i, j int) bool {
	if a[i].Region != a[j].Region {
		if a[i].Region == "" {
			return false
		}
		if a[j].Region == "" {
			return true
		}
		if a[i].Region < a[j].Region {
			return true
		}
		if a[i].Region > a[j].Region {
			return false
		}
	}

	if a[i].SerialCode != a[j].SerialCode {
		if a[i].SerialCode == "" {
			return false
		}
		if a[j].SerialCode == "" {
			return true
		}
		if a[i].SerialCode < a[j].SerialCode {
			return true
		}
		if a[i].SerialCode > a[j].SerialCode {
			return false
		}
	}

	if a[i].Title == "" {
		return false
	}
	if a[j].Title == "" {
		return true
	}
	return a[i].Title < a[j].Title
}
