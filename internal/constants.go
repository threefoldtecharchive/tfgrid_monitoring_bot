package internal

const (
	mainNetwork network = "mainnet"
	testNetwork network = "testnet"
	devNetwork  network = "devnet"
	qaNetwork   network = "qanet"
)

var networks = []network{mainNetwork, testNetwork, devNetwork, qaNetwork}

// SubstrateURLs for substrate urls
var SubstrateURLs = map[network][]string{
	testNetwork: {"wss://tfchain.test.grid.tf/ws", "wss://tfchain.test.grid.tf:443"},
	mainNetwork: {"wss://tfchain.grid.tf/ws", "wss://tfchain.grid.tf:443"},
	devNetwork:  {"wss://tfchain.dev.grid.tf/ws", "wss://tfchain.dev.grid.tf:443"},
	qaNetwork:   {"wss://tfchain.qa.grid.tf/ws", "wss://tfchain.qa.grid.tf:443"},
}

// ProxyUrls for proxy urls
var ProxyUrls = map[network]string{
	testNetwork: "https://gridproxy.test.grid.tf/",
	mainNetwork: "https://gridproxy.grid.tf/",
	devNetwork:  "https://gridproxy.dev.grid.tf/",
	qaNetwork:   "https://gridproxy.qa.grid.tf/",
}

// NodeTwins for proxy urls
var NodeTwins = map[network][]int{
	testNetwork: {3, 4, 5, 7},
	mainNetwork: {7, 16, 22, 24},
	devNetwork:  {9, 11, 25, 21},
	qaNetwork:   {9, 10},
}
