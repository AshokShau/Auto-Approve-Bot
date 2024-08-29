/*
 * Copyright © 2024 AshokShau <github.com/AshokShau>
 */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/AshokShau/Auto-Approve-Bot/Telegram/modules"
	"io"
	"net/http"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

const (
	statusCodeSuccess = 200
)

// Bot Handles all incoming traffic from webhooks.
func Bot(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	split := strings.Split(url, "/")
	if len(split) < 2 {
		fmt.Println(w, "url path too short")
		w.WriteHeader(statusCodeSuccess)

		return
	}

	botToken := split[len(split)-2]

	bot, _ := gotgbot.NewBot(botToken, &gotgbot.BotOpts{DisableTokenCheck: false})

	var update gotgbot.Update

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error reading request body: %v", err)
		w.WriteHeader(statusCodeSuccess)
		return
	}

	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("failed to unmarshal body ", err)
		w.WriteHeader(statusCodeSuccess)

		return
	}

	bot.Username = split[len(split)-1]
	err = modules.Dispatcher.ProcessUpdate(bot, &update, map[string]any{})
	if err != nil {
		fmt.Printf("error while processing update: %v", err)
	}

	w.WriteHeader(statusCodeSuccess)
}
