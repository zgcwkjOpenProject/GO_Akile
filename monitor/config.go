package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Listen     string `json:"listen"`
	Frontend   string `json:"frontend"`
	AuthSecret string `json:"auth_secret"`
	SubRoute   string `json:"subroute"`
	InfoUri    string `json:"info_uri"`
	WsUri      string `json:"ws_uri"`
	DeleteUri  string `json:"delete_uri"`
	UpdateUri  string `json:"update_uri"`
	HookUri    string `json:"hook_uri"`
	HookToken  string `json:"hook_token"`
	TgEnable   bool   `json:"tg_enable"`
	TgToken    string `json:"tg_token"`
	TgChatID   int64  `json:"tg_chat_id"`
}

var cfg *Config

func LoadConfig() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Panic(err)
	}
	cfg = &Config{}
	err = json.Unmarshal(file, cfg)
	if err != nil {
		log.Panic(err)
	}
	return
}
