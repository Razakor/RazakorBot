// MIT License
// Copyright (c) 2018 Maksym Shevchuk (Razakor)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
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
	"encoding/json"
	"os"
)

// BotConfig stores some bot configuration info
type BotConfig struct {
	Version string
	Token   string
	OwnerID int
}

// NewBotConfig creates a new bot config object with default variables
func NewBotConfig() BotConfig {
	c := BotConfig{
		Version: VERSION,
		Token:   "000000000:AAAAaaAAaAAaaAaAAaaAAAAAAAAaaaaAAAa",
		OwnerID: 100000000,
	}
	err := c.ReadConfig("config.json")
	if err != nil {
		c.CreateConfig("config.json")
	}
	return c
}

// CreateConfig makes a config file with some default values
func (c *BotConfig) CreateConfig(filename string) error {
	// This file shouldn't be accessible by anyone else except bot, as there
	// will be stored sensitive data like bot token
	file, err := os.OpenFile(filename,
		os.O_CREATE|os.O_RDWR,
		0600)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	// Make output human-readable
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&c)
	return err
}

// ReadConfig reads existing configuration from file
func (c *BotConfig) ReadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// TODO: allow comments in config, somehow
	// Possible solution is to "preprocess" file before actually decoding it
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return err
}

// TODO: update config on bot update, use Version variable
