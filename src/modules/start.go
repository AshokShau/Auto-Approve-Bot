/*
 * Copyright (c) 2026. AshokShau <github.com/AshokShau>
 */

package modules

import (
	"fmt"
	"time"

	"github.com/AshokShau/Auto-Approve-Bot/src/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var (
	startMessage = "Welcome to the Auto-Approve Bot! This bot will automatically approve all chat join requests. Please ensure that you have the necessary permissions to manage chat members."
	StartTime    = time.Now()
)

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	button := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{
				Text: "Add Me To Your Group",
				Url:  fmt.Sprintf("https://t.me/%s?startgroup=true", b.Username),
			}},
			{
				{
					Text: "Update Channel",
					Url:  "https://t.me/FallenProjects",
				},
				{
					Text: "Help & Commands",
					Url:  "https://AshokShau.github.io/Auto-Approve-Bot/#usage",
				},
			},
		},
	}

	_, err := b.SendPhoto(msg.Chat.Id, gotgbot.InputFileByURL("https://graph.org/file/32166f1230f0a4368438f.jpg"), &gotgbot.SendPhotoOpts{Caption: startMessage, ReplyMarkup: button, HasSpoiler: true, ProtectContent: true})
	if err != nil {
		return fmt.Errorf("[start]failed to send message: %w", err)
	}

	if msg.Chat.Type == gotgbot.ChatTypePrivate {
		_ = db.AddServedUser(msg.Chat.Id)
	} else {
		_ = db.AddServedChat(msg.Chat.Id)
	}

	return ext.EndGroups
}

func ping(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	startTime := time.Now()

	rest, err := msg.Reply(b, "<code>Pinging</code>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return fmt.Errorf("[Ping] failed to send message: %w", err)
	}

	// Calculate latency
	elapsedTime := time.Since(startTime)

	// Calculate uptime
	uptime := time.Since(StartTime)
	formattedUptime := getFormattedDuration(uptime)

	location, _ := time.LoadLocation("Asia/Kolkata")
	responseText := fmt.Sprintf("Pinged in %vms (Latency: %.2fs) at %s\n\nUptime: %s", elapsedTime.Milliseconds(), elapsedTime.Seconds(), time.Now().In(location).Format(time.RFC1123), formattedUptime)

	_, _, err = rest.EditText(b, responseText, nil)
	if err != nil {
		return fmt.Errorf("[Ping] failed to edit message: %w", err)
	}

	return ext.EndGroups
}
