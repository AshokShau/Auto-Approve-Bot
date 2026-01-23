/*
 * Copyright (c) 2026. AshokShau <github.com/AshokShau>
 */

package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	Token      string
	OwnerId    int64
	DbUrl      string
	DbName     string
	WebhookUrl string
	Port       string
)

func LoadEnv() error {
	_ = godotenv.Load()

	DbUrl = os.Getenv("DB_URL")
	DbName = os.Getenv("DB_NAME")
	Token = os.Getenv("TOKEN")
	OwnerId = toInt64(os.Getenv("OWNER_ID"))
	WebhookUrl = os.Getenv("WEBHOOK_URL")
	Port = os.Getenv("PORT")

	if Token == "" {
		return errors.New("bot token required; set TOKEN env variable")

	}
	if DbUrl == "" {
		return errors.New("DB_URL required")
	}
	if DbName == "" {
		DbName = "AutoApproveBot"
	}

	if OwnerId == 0 {
		OwnerId = 5938660179
	}

	return nil
}
