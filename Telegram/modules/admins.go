/*
 * Copyright © 2024 AshokShau <github.com/AshokShau>
 */

package modules

import (
	"fmt"
	"github.com/AshokShau/Auto-Approve-Bot/Telegram/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/go-faster/errors"
)

func isAdmin(b *gotgbot.Bot, chat *gotgbot.Chat, userId int64) bool {
	userMember, err := chat.GetMember(b, userId, nil)
	if err != nil {
		return false
	}

	mem := userMember.MergeChatMember()
	if mem.Status == "member" {
		return false
	}

	return true
}

func autoApprove(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	if chat.Type == "private" {
		_, _ = msg.Reply(b, "This command can only be used in groups or channels", nil)
		return ext.EndGroups
	}

	// Check if the user is an admin in the chat
	if chat.Type != gotgbot.ChatTypeChannel && !isAdmin(b, chat, msg.From.Id) {
		_, _ = msg.Reply(b, "You must be an admin to use this command", nil)
		return ext.EndGroups
	}

	_, _ = b.DeleteMessage(chat.Id, msg.MessageId, nil)
	button := &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Enable", CallbackData: "app_enable"},
				{Text: "Disable", CallbackData: "app_disable"},
			},
		},
	}

	approveStats, _ := db.IsDisabledChat(msg.Chat.Id)
	text := fmt.Sprintf("Auto approve is currently %v\nDo you want to change it?", approveStats)

	_, err := b.SendMessage(ctx.EffectiveChat.Id, text, &gotgbot.SendMessageOpts{ReplyMarkup: button})
	if err != nil {
		return errors.Wrap(err, "[autoApprove] failed to send auto approve message")
	}

	return ext.EndGroups
}

func autoApproveCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	data := ctx.CallbackQuery.Data
	query := ctx.Update.CallbackQuery
	chatId := ctx.EffectiveChat.Id
	user := query.From

	var text string

	if !isAdmin(b, ctx.EffectiveChat, user.Id) {
		_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You must be an admin to use this command", ShowAlert: true})
		return nil
	}

	if data == "app_enable" {
		err := db.EnableApprove(chatId)
		if err != nil {
			return errors.Wrap(err, "[autoApproveCallback] failed to enable auto approve")
		}
		text = "Auto approve enabled"

	} else if data == "app_disable" {
		err := db.DisableApprove(chatId)
		if err != nil {
			return errors.Wrap(err, "[autoApproveCallback] failed to disable auto approve")
		}
		text = "Auto approve disabled"
	}

	_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: text, ShowAlert: true})
	_, _ = b.DeleteMessage(chatId, query.Message.GetMessageId(), nil)
	return ext.EndGroups
}
