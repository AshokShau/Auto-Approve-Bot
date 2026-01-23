/*
 * Copyright (c) 2026. AshokShau <github.com/AshokShau>
 */

package config

import (
	"strconv"
)

func toInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}
