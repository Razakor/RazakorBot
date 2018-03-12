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
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

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

	// Setting global pointer to simplify access from everywhere
	Bot = bot

	// Registering commands
	commands = map[string]command{
		"/start": command(commandStart),
		"/ping":  command(commandPing),
		"/len":   command(commandLen),
		"/rand":  command(commandRand),
		"/words": command(commandWords),
		"/help":  command(commandHelp),
	}

	// Registering help pages
	helpPages = map[string]string{
		"/start": "Print some information about bot and all commands.",
		"/help": "<command> Prints help page for a specified command.\n" +
			"Some commands can accept additional argument. Required arguments are " +
			"specified <like this>, and optional arguments are specified [like " +
			"this]. Braces are ought to be ommited, they just denote if argument " +
			"is required or not.",
		"/ping":  "Do a pong!",
		"/len":   "Measure a length of replied message **in symbols**.",
		"/words": "Measures a length of replied message **in words**.",
		"/rand": "Generate random number in range [1;20] i.e. \"throws a d20\"\n" +
			"/rand <a> [b] generates a number in specified range. " +
			"If b is ommitted or invalid, a number in rande [1;a] will be " +
			"generated. If both number specified correctly, a range will be picked " +
			"automatically: if a > b, then range is [a;b]; if b > a, then range is " +
			"[b;a] obviously.",
	}

	go func() {
		for update := range updates {
			if update.Message == nil || update.Message.Text == "" {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if strings.HasPrefix(update.Message.Text, "/") {
				go ProcessCommand(update)
			}
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	fmt.Println()
	log.Println("Interrupt signal caught, exiting!")
}
