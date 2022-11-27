package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	env := config{}

	testMnemonic := ""
	mainMnemonic := ""
	botToken := ""
	chatId := ""
	mins := 0

	for _, line := range strings.Split(string(content), "\n") {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			continue
		}
		key, value := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])

		switch key {
		case "TESTNET_MNEMONIC":
			testMnemonic = value

		case "MAINNET_MNEMONIC":
			mainMnemonic = value

		case "BOT_TOKEN":
			botToken = value

		case "CHAT_ID":
			chatId = value

		case "MINS":
			intervalMins, err := strconv.Atoi(value)
			if err != nil {
				return env, err
			}
			mins = intervalMins

		default:
			return env, fmt.Errorf("key %v is invalid", key)
		}
	}

	switch {
	case testMnemonic == "":
		return env, fmt.Errorf("TESTNET_MNEMONIC is missing")
	case mainMnemonic == "":
		return env, fmt.Errorf("MAINNET_MNEMONIC is missing")
	case botToken == "":
		return env, fmt.Errorf("BOT_TOKEN is missing")
	case chatId == "":
		return env, fmt.Errorf("CHAT_ID is missing")
	case mins == 0:
		return env, fmt.Errorf("MINS is 0")
	}

	return config{
		testMnemonic: testMnemonic,
		mainMnemonic: mainMnemonic,
		botToken:     botToken,
		chatId:       chatId,
		intervalMins: mins,
	}, nil
}
