/*
 * Copyright © 2024 AshokShau <github.com/AshokShau>
 */

package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Token   string
	OwnerId int64
	DbUrl   string
	DbName  string
)

func init() {
	_ = godotenv.Load()

	DbUrl = os.Getenv("DB_URL")
	DbName = os.Getenv("DB_NAME")
	Token = os.Getenv("TOKEN")
	OwnerId = toInt64(os.Getenv("OWNER_ID"))

	setDefaults()
}
