package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	client "github.com/threefoldtech/substrate-client"
)

type address string
type network string

var (
	mainNetwork network = "mainnet"
	testNetwork network = "testnet"
)

var SUBSTRATE_URLS = map[network][]string{
	testNetwork: {"wss://tfchain.test.grid.tf/ws"},
	mainNetwork: {"wss://tfchain.grid.tf/ws"},
}

type config struct {
	testMnemonic string
	mainMnemonic string
	botToken     string
	chatId       string
	intervalMins int
}

type wallet struct {
	address   address
	threshold int
	name      string
}
type wallets struct {
	mainnet []wallet
	testnet []wallet
}

type monitor struct {
	env       config
	wallets   wallets
	substrate map[network]client.Manager
}

// NewMonitor creates a new instance of monitor
func NewMonitor(envPath string, jsonPath string) (monitor, error) {
	mon := monitor{}

	envContent, err := readFile(envPath)
	if err != nil {
		return mon, err
	}

	env, err := parseEnv(string(envContent))
	if err != nil {
		return mon, err
	}

	jsonContent, err := readFile(jsonPath)
	if err != nil {
		return mon, err
	}

	addresses, err := parseJsonIntoWallets(jsonContent)
	if err != nil {
		return mon, err
	}

	mon.wallets = addresses
	mon.env = env

	substrate := map[network]client.Manager{}

	if len(mon.wallets.mainnet) != 0 {
		substrate[mainNetwork] = client.NewManager(SUBSTRATE_URLS[mainNetwork]...)
	}
	if len(mon.wallets.testnet) != 0 {
		substrate[testNetwork] = client.NewManager(SUBSTRATE_URLS[testNetwork]...)
	}

	mon.substrate = substrate

	return mon, nil
}

// Start starting the monitoring service
func (m *monitor) Start() {
	ticker := time.NewTicker(time.Duration(m.env.intervalMins) * time.Minute)

	for range ticker.C {
		for network, manager := range m.substrate {

			wallets := []wallet{}
			switch network {
			case mainNetwork:
				wallets = m.wallets.mainnet
			case testNetwork:
				wallets = m.wallets.testnet
			}

			for _, address := range wallets {
				log.Debug().Msgf("monitoring for network %v, address %v", network, address)
				err := m.sendMessage(manager, address)
				if err != nil {
					log.Error().Err(err).Msg("monitoring failed with error")
				}
			}
		}
	}
}

// getTelegramUrl returns the telegram bot api url
func (m *monitor) getTelegramUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", m.env.botToken)
}

// sendMessage sends a message with the balance to a telegram bot
// if it is less than the tft threshold
func (m *monitor) sendMessage(manager client.Manager, wallet wallet) error {
	balance, err := m.getBalance(manager, wallet.address)
	if err != nil {
		return err
	}

	if balance >= float64(wallet.threshold) {
		return nil
	}

	url := fmt.Sprintf("%s/sendMessage", m.getTelegramUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": m.env.chatId,
		"text":    fmt.Sprintf("account with name:\n%v\nhas balance = %v", wallet.name, balance),
	})
	response, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return errors.New("request send message failed")
	}

	defer response.Body.Close()
	return nil
}

// getBalance gets the balance in TFT for the address given
func (m *monitor) getBalance(manager client.Manager, address address) (float64, error) {
	log.Debug().Msgf("get balance for %v", address)

	con, err := manager.Substrate()
	if err != nil {
		return 0, err
	}
	defer con.Close()

	account, err := client.FromAddress(string(address))
	if err != nil {
		return 0, err
	}

	balance, err := con.GetBalance(account)
	if err != nil {
		return 0, err
	}

	return float64(balance.Free.Int64()) / math.Pow(10, 7), nil
}
