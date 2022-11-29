package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	env "github.com/hashicorp/go-envparse"
)

func readFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func parseJsonIntoWallets(content []byte) (w wallets, err error) {
	w = wallets{}
	err = json.Unmarshal(content, &w)

	if err != nil {
		return
	}

	return
}

func parseEnv(content string) (config, error) {
	cfg := config{}

	configMap, err := env.Parse(strings.NewReader(content))
	if err != err {
		return config{}, err
	}

	for key, value := range configMap {
		switch key {
		case "TESTNET_MNEMONIC":
			cfg.testMnemonic = value

		case "MAINNET_MNEMONIC":
			cfg.mainMnemonic = value

		case "BOT_TOKEN":
			cfg.botToken = value

		case "CHAT_ID":
			cfg.chatId = value

		case "MINS":
			intervalMins, err := strconv.Atoi(value)
			if err != nil {
				return config{}, err
			}
			cfg.intervalMins = intervalMins

		default:
			return config{}, fmt.Errorf("key %v is invalid", key)
		}
	}

	switch {
	case cfg.testMnemonic == "":
		return config{}, fmt.Errorf("TESTNET_MNEMONIC is missing")
	case cfg.mainMnemonic == "":
		return config{}, fmt.Errorf("MAINNET_MNEMONIC is missing")
	case cfg.botToken == "":
		return config{}, fmt.Errorf("BOT_TOKEN is missing")
	case cfg.chatId == "":
		return config{}, fmt.Errorf("CHAT_ID is missing")
	case cfg.intervalMins == 0:
		return config{}, fmt.Errorf("MINS is 0")
	}

	return cfg, nil
}
