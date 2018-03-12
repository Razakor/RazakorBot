package main

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"gopkg.in/telegram-bot-api.v4"
)

type command func(args []string, update tgbotapi.Update) string

// commands is a map of command strings and function for invocation
var commands map[string]command

func commandLen(args []string, update tgbotapi.Update) string {
	if update.Message.ReplyToMessage == nil {
		return "Invalid command invocation"
	}
	if update.Message.ReplyToMessage.Text == "dick" {
		return "42"
	} else {
		return strconv.Itoa(utf8.RuneCountInString(update.Message.ReplyToMessage.Text))
	}
}

func commandRand(args []string, update tgbotapi.Update) string {
	if len(args) < 2 {
		return strconv.Itoa(rand.Intn(19) + 1)
	} else {
		var b int
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return "Invalid command invocation"
		}
		if len(args) > 2 {
			b, err = strconv.Atoi(args[2])
			if err != nil {
				b = 1
			}
		} else {
			b = 1
		}
		if n < 1 || b < 1 || n == b {
			return "Invalid command invocation"
		}
		if n < b {
			return strconv.Itoa(b - rand.Intn((b-n)+1))
		}
		return strconv.Itoa(n - rand.Intn((n-b)+1))
	}
}

func commandWords(args []string, update tgbotapi.Update) string {
	m := strings.Split(update.Message.ReplyToMessage.Text, " ")
	return strconv.Itoa(len(m))
}

// ProcessCommand processes received update and executes command if it is valid
func ProcessCommand(update tgbotapi.Update) {
	args := strings.Split(update.Message.Text, " ")
	var response string
	//chatID := update.Message.Chat.ID
	if strings.Contains(args[0], "@") {
		response = strings.Split(args[0], "@")[0]
	} else {
		response = args[0]
	}
	if val, ok := commands[response]; ok {
		response = val(args, update)
	} else {
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	Bot.Send(msg)
}
