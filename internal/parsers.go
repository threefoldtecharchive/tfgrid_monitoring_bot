package internal

import (
	"encoding/json"
	"errors"
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
	addresses := map[string][]address{}
	err = json.Unmarshal(content, &addresses)

	if err != nil {
		return
	}

	w = wallets{}
	if _, ok := addresses["mainnet"]; !ok {
		return w, errors.New("mainnet addresses is missing")
	} else {
		w.mainnet = addresses["mainnet"]
	}

	if _, ok := addresses["testnet"]; !ok {
		return w, errors.New("testnet addresses is missing")
	} else {
		w.testnet = addresses["testnet"]
	}

	return w, err
}

func parseEnv(content string) (config, error) {
	env := config{}

	testMnemonic := ""
	mainMnemonic := ""
	tftLimit := 0
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

		case "TFTS_LIMIT":
			limit, err := strconv.Atoi(value)
			if err != nil {
				return env, err
			}
			tftLimit = limit

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
	case tftLimit == 0:
		return env, fmt.Errorf("TFTS_LIMIT is 0")
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
		tftLimit:     tftLimit,
		botToken:     botToken,
		chatId:       chatId,
		intervalMins: mins,
	}, nil
}
