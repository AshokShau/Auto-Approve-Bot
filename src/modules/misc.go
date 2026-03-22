/*
 * Copyright (c) 2026. AshokShau <github.com/AshokShau>
 */

package modules

import (
	"fmt"
	"strings"
	"time"
)

// plural returns the unit with correct singular/plural form.
func plural(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}

// getFormattedDuration returns a human-readable string for the given duration.
func getFormattedDuration(diff time.Duration) string {
	totalSeconds := int(diff.Seconds())

	months := totalSeconds / (30 * 24 * 3600)
	remaining := totalSeconds % (30 * 24 * 3600)

	weeks := remaining / (7 * 24 * 3600)
	remaining = remaining % (7 * 24 * 3600)

	days := remaining / (24 * 3600)
	remaining = remaining % (24 * 3600)

	hours := remaining / 3600
	remaining = remaining % 3600

	minutes := remaining / 60
	seconds := remaining % 60

	var parts []string

	if months > 0 {
		parts = append(parts, plural(months, "month"))
	}
	if weeks > 0 {
		parts = append(parts, plural(weeks, "week"))
	}
	if days > 0 {
		parts = append(parts, plural(days, "day"))
	}
	if hours > 0 {
		parts = append(parts, plural(hours, "hour"))
	}
	if minutes > 0 {
		parts = append(parts, plural(minutes, "minute"))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, plural(seconds, "second"))
	}

	return strings.Join(parts, " ")
}
