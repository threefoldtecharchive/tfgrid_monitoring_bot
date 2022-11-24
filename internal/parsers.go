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
	addresses := map[string][]map[string]interface{}{}
	err = json.Unmarshal(content, &addresses)

	if err != nil {
		return
	}

	w = wallets{}
	if _, ok := addresses["mainnet"]; !ok {
		return w, errors.New("mainnet addresses is missing")
	} else {
		for _, walletMap := range addresses["mainnet"] {
			mainnetWallet := wallet{}
			mainnetWallet.address = address(walletMap["addr"].(string))
			mainnetWallet.limit = int(walletMap["threshold"].(float64))
			w.mainnet = append(w.mainnet, mainnetWallet)
		}
	}

	if _, ok := addresses["testnet"]; !ok {
		return w, errors.New("testnet addresses is missing")
	} else {
		for _, walletMap := range addresses["testnet"] {
			testnetWallet := wallet{}
			testnetWallet.address = address(walletMap["addr"].(string))
			testnetWallet.limit = int(walletMap["threshold"].(float64))
			w.testnet = append(w.testnet, testnetWallet)
		}
	}

	return w, err
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
