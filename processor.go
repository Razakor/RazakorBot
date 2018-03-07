package main

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"gopkg.in/telegram-bot-api.v4"
)

func ProcessCommand(update tgbotapi.Update, args []string) {
	var response string
	chatID := update.Message.Chat.ID
	if strings.HasPrefix(args[0], "/len") {
		if update.Message.ReplyToMessage == nil {
			return
		}
		if update.Message.ReplyToMessage.Text == "dick" {
			response = "42"
		} else {
			response = strconv.Itoa(utf8.RuneCountInString(update.Message.ReplyToMessage.Text))
		}
	} else if strings.HasPrefix(args[0], "/rand") {
		if len(args) < 2 {
			response = strconv.Itoa(rand.Intn(19) + 1)
		} else {
			var n, b int
			n, err := strconv.Atoi(args[1])
			if err != nil {
				Bot.Send(tgbotapi.NewMessage(chatID, "Invalid arguments!"))
				return
			}
			if len(args) > 2 {
				b, err = strconv.Atoi(args[2])
				if err != nil {
					Bot.Send(tgbotapi.NewMessage(chatID, "Invalid arguments!"))
					return
				}
			}
			if n < 0 || b < 1 || b < n {
				Bot.Send(tgbotapi.NewMessage(chatID, "Invalid arguments!"))
				return
			}
			response = strconv.Itoa(b - rand.Intn((b-n)+1))
		}
	} else if strings.HasPrefix(args[0], "/words") {
		m := strings.Split(update.Message.ReplyToMessage.Text, " ")
		response = strconv.Itoa(len(m))
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	Bot.Send(msg)
}
