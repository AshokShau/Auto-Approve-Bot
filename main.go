/*
 * Copyright © 2024 AshokShau <github.com/AshokShau>
 */

package main

import (
	"github.com/AshokShau/Auto-Approve-Bot/Telegram/config"
	"github.com/AshokShau/Auto-Approve-Bot/Telegram/modules"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func main() {
	if config.Token == "" {
		log.Fatal("Bot Token required")
	}

	bot, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		log.Fatal("failed to create bot:", err)
	}

	updater := ext.NewUpdater(modules.Dispatcher, nil)
	allowedUpdates := []string{gotgbot.UpdateTypeMessage, gotgbot.UpdateTypeCallbackQuery, gotgbot.UpdateTypeChannelPost, gotgbot.UpdateTypeChatJoinRequest}

	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout:        9,
			AllowedUpdates: allowedUpdates,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	if err != nil {
		log.Fatal("failed to start polling:", err)
	}

	log.Printf("Bot started as %s", bot.Username)
	_, _ = bot.SendMessage(config.OwnerId, "Bot started;", nil)

	updater.Idle()
}
