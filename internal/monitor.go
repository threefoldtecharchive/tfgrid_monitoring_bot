package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	client "github.com/threefoldtech/substrate-client"
	mainClient "github.com/threefoldtech/substrate-client-main"
)

type address string
type network string

type config struct {
	testMnemonic string `env:"TESTNET_MNEMONIC"`
	mainMnemonic string `env:"MAINNET_MNEMONIC"`
	devMnemonic  string `env:"DEVNET_MNEMONIC"`
	qaMnemonic   string `env:"QANET_MNEMONIC"`
	farmName     string `env:"FARM_NAME"`
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
	wallets                   wallets
	workingNodesPerNetwork    map[network][]uint32
	notWorkingNodesPerNetwork map[network][]uint32
	substrate                 map[network]client.Manager
	mainSubstrate             map[network]mainClient.Manager
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
	mon.mainSubstrate = map[network]mainClient.Manager{}

	// all needed for proxy
	for _, network := range networks {
		if network == mainNetwork || network == testNetwork {
			mon.mainSubstrate[network] = mainClient.NewManager(SubstrateURLs[network]...)
		}
		mon.substrate[network] = client.NewManager(SubstrateURLs[network]...)
	}

	mon.mnemonics = map[network]string{}
	mon.mnemonics[devNetwork] = mon.env.devMnemonic
	mon.mnemonics[testNetwork] = mon.env.testMnemonic
	mon.mnemonics[qaNetwork] = mon.env.qaMnemonic
	mon.mnemonics[mainNetwork] = mon.env.mainMnemonic

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
	versions, err := m.systemVersion(context.Background())
	if err != nil {
		return err
	}

	message := ""

	for _, network := range networks {

		if _, ok := versions[network]; !ok {
			notWorkingTestedNodes := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(m.notWorkingNodesPerNetwork[network])), ", "), "[]")
			message += fmt.Sprintf("Proxy for %v is not working ❌\nNodes tested but failed: %v\n\n", network, notWorkingTestedNodes)
			continue
		}

		workingTestedNodes := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(m.workingNodesPerNetwork[network])), ", "), "[]")
		message += fmt.Sprintf("Proxy for %v is working ✅\nNodes successfully tested: %v\n\n", network, workingTestedNodes)
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

type version struct {
	ZOS   string `json:"zos"`
	ZInit string `json:"zinit"`
}

// systemVersion executes system version cmd
func (m *Monitor) systemVersion(ctx context.Context) (map[network]version, error) {
	versions := map[network]version{}

	for _, network := range networks {
		log.Debug().Msgf("get system version for network %v", network)

		identity, err := NewIdentityFromSr25519Phrase(m.mnemonics[network])
		if err != nil {
			log.Error().Err(err).Msgf("creating new identity for %v network failed", network)
			continue
		}

		con, err := m.substrate[network].Substrate()
		if err != nil {
			log.Error().Err(err).Msgf("substrate connection for %v network failed", network)
			continue
		}
		defer con.Close()

		twinID, err := con.GetTwinByPubKey(identity.PublicKey())
		if err != nil {
			log.Error().Err(err).Msgf("returning twin ID for %v network failed", network)
			continue
		}

		devProxyBus, err := NewProxyBus(ProxyUrls[network], twinID, *con, identity, true)
		if err != nil {
			log.Error().Err(err).Msgf("proxy bus for %v network failed", network)
			continue
		}

		farmID, err := con.GetFarmByName(m.env.farmName)
		if err != nil {
			log.Error().Err(err).Msgf("cannot get farm ID for farm '%s'", m.env.farmName)
			continue
		}

		farmNodes, err := con.GetNodes(farmID)
		if err != nil {
			log.Error().Err(err).Msgf("cannot get farm nodes for farm %d", farmID)
			continue
		}

		rand.Shuffle(len(farmNodes), func(i, j int) { farmNodes[i], farmNodes[j] = farmNodes[j], farmNodes[i] })
		var randomNodes []uint32
		if len(farmNodes) < 4 {
			randomNodes = farmNodes[:]
		} else {
			randomNodes = farmNodes[:4]
		}

		for _, NodeID := range randomNodes {

			// substrate for main and test is different
			if network == mainNetwork || network == testNetwork {
				mainCon, err := m.mainSubstrate[network].Substrate()
				if err != nil {
					log.Error().Err(err).Msgf("substrate connection for %v network failed", network)
					continue
				}
				defer mainCon.Close()

				ver, err := m.checkNodeSystemVersionMainnet(ctx, mainCon, devProxyBus, NodeID, network)
				if err != nil {
					log.Error().Err(err).Msgf("check node %d failed", NodeID)
					continue
				}

				versions[network] = ver
				continue
			}

			ver, err := m.checkNodeSystemVersion(ctx, con, devProxyBus, NodeID, network)
			if err != nil {
				log.Error().Err(err).Msgf("check node %d failed", NodeID)
				continue
			}

			versions[network] = ver
		}
	}

	return versions, nil
}

func (m *Monitor) checkNodeSystemVersion(ctx context.Context, con *client.Substrate, proxyBus *ProxyBus, NodeID uint32, net network) (version, error) {
	const cmd = "zos.system.version"
	var ver version

	node, err := con.GetNode(NodeID)
	if err != nil {
		m.notWorkingNodesPerNetwork[net] = append(m.notWorkingNodesPerNetwork[net], NodeID)
		return ver, fmt.Errorf("cannot get node %d. failed with error: %w", NodeID, err)
	}

	err = proxyBus.Call(ctx, uint32(node.TwinID), cmd, nil, &ver)
	if err != nil {
		m.notWorkingNodesPerNetwork[net] = append(m.notWorkingNodesPerNetwork[net], NodeID)
		return ver, fmt.Errorf("proxy bus getting system version for %v network failed using node twin %v with node ID %v. failed with error: %w", net, node.TwinID, NodeID, err)
	}

	m.workingNodesPerNetwork[net] = append(m.workingNodesPerNetwork[net], NodeID)
	return ver, nil
}

func (m *Monitor) checkNodeSystemVersionMainnet(ctx context.Context, con *mainClient.Substrate, proxyBus *ProxyBus, NodeID uint32, net network) (version, error) {
	const cmd = "zos.system.version"
	var ver version

	node, err := con.GetNode(NodeID)
	if err != nil {
		m.notWorkingNodesPerNetwork[net] = append(m.notWorkingNodesPerNetwork[net], NodeID)
		return ver, fmt.Errorf("cannot get node %d. failed with error: %w", NodeID, err)
	}

	err = proxyBus.Call(ctx, uint32(node.TwinID), cmd, nil, &ver)
	if err != nil {
		m.notWorkingNodesPerNetwork[net] = append(m.notWorkingNodesPerNetwork[net], NodeID)
		return ver, fmt.Errorf("proxy bus getting system version for %v network failed using node twin %v with node ID %v. failed with error: %w", net, node.TwinID, NodeID, err)
	}

	m.workingNodesPerNetwork[net] = append(m.workingNodesPerNetwork[net], NodeID)
	return ver, nil
}
