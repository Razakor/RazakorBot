// MIT License
// Copyright (c) 2018 Maksym Shevchuk (Razakor)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"

	"gopkg.in/telegram-bot-api.v4"
)

// logRotate rotates all logs, saving them for future use
func logRotate() {
	files, err := ioutil.ReadDir("logs/")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("logs/", 0740)
			return
		} else {
			log.Fatal(err)
		}
	}
	var logs []string
	for _, i := range files {
		if strings.Contains(i.Name(), "bot.") &&
			strings.Contains(i.Name(), ".log") {
			logs = append(logs, i.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(logs)))
	for _, i := range logs {
		arr := strings.Split(i, ".")
		num, err := strconv.Atoi(arr[len(arr)-2])
		if err != nil {
			log.Fatal(err)
		}
		os.Rename("logs/"+i, "logs/bot."+strconv.Itoa(num+1)+".log")
	}
}

func main() {
	log.Println("RazakorBot")
	log.Println("Copyright (C) 2018  " +
		"Maksym Shevchuk (Razakor)" +
		"\n\tSee LICENSE file for more info")

	{ // Initialize logging to file and stdout
		// UTC time
		log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
		logRotate()
		logFile, err := os.OpenFile("logs/bot.0.log",
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0660)
		if err != nil {
			log.Panicln(err)
		}
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
	}

	log.Println("Loading config")
	config := NewBotConfig()

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message == nil || update.Message.Text == "" {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if strings.HasPrefix(update.Message.Text, "/") {
				go func(args []string) {
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
								bot.Send(tgbotapi.NewMessage(chatID, "Invalid arguments!"))
								return
							}
							if len(args) > 2 {
								b, err = strconv.Atoi(args[2])
								if err != nil {
									bot.Send(tgbotapi.NewMessage(chatID, "Invalid arguments!"))
									return
								}
							}
							if n < 0 || b < 1 || b < n {
								bot.Send(tgbotapi.NewMessage(chatID, "Invalid arguments!"))
								return
							}
							response = strconv.Itoa(b - rand.Intn((b-n)+1))
						}
					} else if strings.HasPrefix(args[0], "/words") {
						m := strings.Split(update.Message.ReplyToMessage.Text, " ")
						response = strconv.Itoa(len(m))
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
					bot.Send(msg)
				}(strings.Split(update.Message.Text, " "))
			}
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	fmt.Println()
	log.Println("Interrupt signal caught, exiting!")
}
