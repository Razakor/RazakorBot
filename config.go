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
