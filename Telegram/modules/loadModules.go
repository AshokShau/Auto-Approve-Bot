/*
 * Copyright © 2024 AshokShau <github.com/AshokShau>
 */

package modules

import (
	"fmt"
	"github.com/AshokShau/Auto-Approve-Bot/Telegram/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/chatjoinrequest"
	"html"
)

func errorHandler(bot *gotgbot.Bot, _ *ext.Context, err error) ext.DispatcherAction {
	msg := fmt.Sprintf("%s", html.EscapeString(err.Error()))
	if _, err = bot.SendMessage(config.OwnerId, msg, &gotgbot.SendMessageOpts{ParseMode: "HTML"}); err != nil {
		_ = fmt.Errorf("failed to send error message: %v", err)
		return ext.DispatcherActionNoop
	}

	return ext.DispatcherActionNoop
}

var Dispatcher = newDispatcher()

func newDispatcher() *ext.Dispatcher {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{Error: errorHandler, MaxRoutines: 50})
	loadModules(dispatcher)

	return dispatcher
}

func loadModules(d *ext.Dispatcher) {
	d.AddHandler(handlers.NewCommand("start", start))
	d.AddHandler(handlers.NewCommand("ping", ping))
	d.AddHandler(handlers.NewCommand("stats", stats))
	d.AddHandler(handlers.NewCommand("broadcast", broadCast))

	d.AddHandler(handlers.Command{
		AllowChannel: true,
		Response:     autoApprove,
		Triggers:     []rune{'/'},
		Command:      "autoapprove",
	})

	d.AddHandler(handlers.NewChatJoinRequest(chatjoinrequest.All, joinRequest))
	d.AddHandler(handlers.NewCallback(callbackquery.Prefix("app_"), autoApproveCallback))
}
