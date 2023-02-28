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

type config struct {
	testMnemonic string `env:"TESTNET_MNEMONIC"`
	mainMnemonic string `env:"MAINNET_MNEMONIC"`
	devMnemonic  string `env:"DEVNET_MNEMONIC"`
	qaMnemonic   string `env:"QANET_MNEMONIC"`
	devFarmName  string `env:"DEV_FARM_NAME"`
	qaFarmName   string `env:"QA_FARM_NAME"`
	testFarmName string `env:"TEST_FARM_NAME"`
	mainFarmName string `env:"MAIN_FARM_NAME"`
	botToken     string `env:"BOT_TOKEN"`
	chatID       string `env:"CHAT_ID"`
	intervalMins int    `env:"MINS"`
}

type wallet struct {
	Address   address `json:"address"`
	Threshold int     `json:"threshold"`
	Name      string  `json:"name"`
}
type wallets struct {
	Mainnet []wallet `json:"mainnet"`
	Testnet []wallet `json:"testnet"`
}

// Monitor for bot monitoring
type Monitor struct {
	env                       config
	mnemonics                 map[network]string
	farms                     map[network]string
	wallets                   wallets
	workingNodesPerNetwork    map[network][]uint32
	notWorkingNodesPerNetwork map[network][]uint32
	substrate                 map[network]client.Manager
}

// NewMonitor creates a new instance of monitor
func NewMonitor(envPath string, jsonPath string) (Monitor, error) {
	mon := Monitor{}

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

	addresses, err := parseJSONIntoWallets(jsonContent)
	if err != nil {
		return mon, err
	}

	mon.wallets = addresses
	mon.env = env

	mon.substrate = map[network]client.Manager{}

	// all needed for proxy
	for _, network := range networks {
		mon.substrate[network] = client.NewManager(SubstrateURLs[network]...)
	}

	mon.mnemonics = map[network]string{}
	mon.mnemonics[devNetwork] = mon.env.devMnemonic
	mon.mnemonics[testNetwork] = mon.env.testMnemonic
	mon.mnemonics[qaNetwork] = mon.env.qaMnemonic
	mon.mnemonics[mainNetwork] = mon.env.mainMnemonic

	mon.farms = map[network]string{}
	mon.farms[devNetwork] = mon.env.devFarmName
	mon.farms[testNetwork] = mon.env.testFarmName
	mon.farms[qaNetwork] = mon.env.qaFarmName
	mon.farms[mainNetwork] = mon.env.mainFarmName

	mon.workingNodesPerNetwork = map[network][]uint32{}
	mon.notWorkingNodesPerNetwork = map[network][]uint32{}

	return mon, nil
}

// Start starting the monitoring service
func (m *Monitor) Start() {
	ticker := time.NewTicker(time.Duration(m.env.intervalMins) * time.Minute)

	for range ticker.C {
		for network, manager := range m.substrate {

			wallets := []wallet{}
			switch network {
			case mainNetwork:
				wallets = m.wallets.Mainnet
			case testNetwork:
				wallets = m.wallets.Testnet
			}

			for _, wallet := range wallets {
				log.Debug().Msgf("monitoring for network %v, address %v", network, wallet.Address)
				err := m.sendMessage(manager, wallet)
				if err != nil {
					log.Error().Err(err).Msg("monitoring failed with error")
				}
			}
		}

		log.Debug().Msgf("monitoring proxy for all networks")
		err := m.sendProxyCheckMessage()
		if err != nil {
			log.Error().Err(err).Msg("monitoring proxy failed with error")
		}
	}
}

// getTelegramUrl returns the telegram bot api url
func (m *Monitor) getTelegramURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", m.env.botToken)
}

// sendMessage sends a message with the balance to a telegram bot
// if it is less than the tft threshold
func (m *Monitor) sendMessage(manager client.Manager, wallet wallet) error {
	balance, err := m.getBalance(manager, wallet.Address)
	if err != nil {
		return err
	}

	if balance >= float64(wallet.Threshold) {
		return nil
	}

	url := fmt.Sprintf("%s/sendMessage", m.getTelegramURL())
	body, _ := json.Marshal(map[string]string{
		"chat_id": m.env.chatID,
		"text":    fmt.Sprintf("wallet %v with address:\n%v\nhas balance = %v", wallet.Name, wallet.Address, balance),
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

// sendProxyCheckMessage checks if proxy is working against all networks
func (m *Monitor) sendProxyCheckMessage() error {
	m.notWorkingNodesPerNetwork = map[network][]uint32{}
	m.workingNodesPerNetwork = map[network][]uint32{}

	gridProxyHealthCheck, err := m.pingGridNetworks()
	if err != nil {
		return err
	}

	message := ""

	for _, network := range networks {

		if _, ok := gridProxyHealthCheck[network]; !ok {
			message += fmt.Sprintf("Proxy for %v is not working ❌\n\n", network)
			continue
		}

		message += fmt.Sprintf("Proxy for %v is working ✅\n\n", network)
	}

	url := fmt.Sprintf("%s/sendMessage", m.getTelegramURL())
	body, _ := json.Marshal(map[string]string{
		"chat_id": m.env.chatID,
		"text":    message,
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
		return errors.New("request send proxy check message failed")
	}

	defer response.Body.Close()
	return nil
}

// getBalance gets the balance in TFT for the address given
func (m *Monitor) getBalance(manager client.Manager, address address) (float64, error) {
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

// pingGridNetwork pings the different grid proxy networks
func (m *Monitor) pingGridNetworks() (map[network]bool, error) {
	gridProxyHealthCheck := map[network]bool{}

	for _, network := range networks {
		log.Debug().Msgf("pinging grid proxy for network %s", network)
		gridProxy, err := NewGridProxyClient(ProxyUrls[network])
		if err != nil {
			log.Error().Err(err).Msgf("grid proxy for %v network failed", network)
			continue
		}

		if err := gridProxy.Ping(); err != nil {
			log.Error().Err(err).Msgf("failed to ping grid proxy on network %v", network)
			continue
		}
		gridProxyHealthCheck[network] = true
	}
	return gridProxyHealthCheck, nil
}
