/*
 * Copyright (c) 2026. AshokShau <github.com/AshokShau>
 */

package modules

import (
	"fmt"
	"log"
	"runtime"

	"github.com/AshokShau/Auto-Approve-Bot/src/config"
	"github.com/AshokShau/Auto-Approve-Bot/src/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func getGoroutineCount() int {
	return runtime.NumGoroutine()
}

func getGoVersion() string {
	return runtime.Version()
}

func stats(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user.Id != config.OwnerId {
		return ext.EndGroups
	}

	chatCount, err := db.GetChatCount()
	if err != nil {
		return fmt.Errorf("GetChatCount: %w", err)
	}

	userCount, err := db.GetUserCount()
	if err != nil {
		return fmt.Errorf("GetUserCount: %w", err)
	}

	statsMessage := fmt.Sprintf(
		"📊 <b>Statistics</b>\n\n"+
			"👥 Users: %d\n"+
			"👥 Chats: %d\n"+
			"🌀 Goroutines: %d\n"+
			"🔢 Go Version: %s",
		userCount, chatCount, getGoroutineCount(), getGoVersion())

	_, err = ctx.EffectiveMessage.Reply(b, statsMessage, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return fmt.Errorf("ReplyMessage: %w", err)
	}

	return ext.EndGroups
}

func broadCast(b *gotgbot.Bot, ctx *ext.Context) error {
	useR := ctx.EffectiveUser
	if useR.Id != config.OwnerId {
		return ext.EndGroups
	}

	reply := ctx.EffectiveMessage.ReplyToMessage
	if reply == nil {
		_, err := ctx.EffectiveMessage.Reply(b, "❌ <b>Reply to a message to broadcast</b>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
		if err != nil {
			return fmt.Errorf("[broadcast] failed to send reply message ")
		}
		return ext.EndGroups
	}

	button := &gotgbot.InlineKeyboardMarkup{}
	if reply.ReplyMarkup != nil {
		button.InlineKeyboard = reply.ReplyMarkup.InlineKeyboard
	}

	servedUsers, err := db.GetServedUsers()
	if err != nil {
		return fmt.Errorf("[broadcast] failed to get servedUsers %w", err)
	}

	successfulBroadcasts := 0
	for _, user := range servedUsers {
		userId := user["user_id"].(int64)
		_, err = b.CopyMessage(userId, ctx.EffectiveMessage.Chat.Id, reply.MessageId, &gotgbot.CopyMessageOpts{ReplyMarkup: button})
		if err != nil {
			log.Printf("[broadcast] failed to copy message to %d: %v", userId, err)
		} else {
			successfulBroadcasts++
		}
	}

	servedChats, err := db.GetServedChats()
	if err != nil {
		return fmt.Errorf("[broadcast] failed to get servedChats %w", err)
	}

	successfulBroadcastsChats := 0
	for _, chat := range servedChats {
		chatId := chat["chat_id"].(int64)
		_, err = b.CopyMessage(chatId, ctx.EffectiveMessage.Chat.Id, reply.MessageId, &gotgbot.CopyMessageOpts{ReplyMarkup: button})
		if err != nil {
			log.Printf("[broadcast] failed to copy message to %d: %v", chatId, err)
		} else {
			successfulBroadcastsChats++
		}
	}

	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("✅ <b>Broadcast successfully to %d users and %d chats</b>", successfulBroadcasts, successfulBroadcastsChats), &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return fmt.Errorf("[broadcast] failed to send reply message %w", err)
	}

	return ext.EndGroups
}
