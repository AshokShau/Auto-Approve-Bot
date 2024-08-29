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
	"strings"
)

func joinRequest(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.ChatJoinRequest.Chat
	user := ctx.ChatJoinRequest.From
	button := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Add me to the group", Url: fmt.Sprintf("https://t.me/%s?startgroup=true", b.Username)},
			},
		},
	}

	text := fmt.Sprintf("Hello %s, welcome to %s!\nTap on /start .", user.FirstName, chat.Title)
	_, err := b.SendMessage(user.Id, text, &gotgbot.SendMessageOpts{ReplyMarkup: button})
	if err != nil && !strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") {
		return errors.Wrap(err, "[joinRequest] failed to send join request message")
	}

	approveStats, _ := db.IsDisabledChat(chat.Id)
	if !approveStats {
		return ext.EndGroups
	}

	_, err = b.ApproveChatJoinRequest(chat.Id, user.Id, &gotgbot.ApproveChatJoinRequestOpts{})
	_ = db.AddServedChat(chat.Id)
	if err != nil {
		return errors.Wrap(err, "[joinRequest] failed to approve chat join request")
	}

	return ext.EndGroups
}
