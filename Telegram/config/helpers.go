/*
 * Copyright © 2024 Abishnoi69 <github.com/Abishnoi69>
 */

package config

import (
	"log"
	"strconv"
)

func toInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func setDefaults() {
	if DbUrl == "" {
		log.Fatal("MongoDB DatabaseURL required")
	}

	if DbName == "" {
		DbName = "AutoApproveBot"
	}

	if OwnerId == 0 {
		OwnerId = 5938660179
	}

}
