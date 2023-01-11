module github.com/rawdaGastan/tfgrid_monitoring_bot

go 1.19

require (
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/rs/zerolog v1.26.0
	github.com/spf13/cobra v1.6.1
	github.com/threefoldtech/go-rmb v0.2.2
	github.com/threefoldtech/substrate-client v0.0.0-20230111124412-c251782691c3
	github.com/threefoldtech/substrate-client-main v0.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	golang.org/x/net v0.0.0-20220403103023-749bd193bc2b // indirect
)

require (
	github.com/ChainSafe/go-schnorrkel v1.0.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/centrifuge/go-substrate-rpc-client/v4 v4.0.5 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/base58 v1.0.3 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/ethereum/go-ethereum v1.10.17 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gtank/merlin v0.1.1 // indirect
	github.com/gtank/ristretto255 v0.1.2 // indirect
	github.com/hashicorp/go-envparse v0.1.0
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6 // indirect
	github.com/mimoo/StrobeGo v0.0.0-20210601165009-122bf33a46e0 // indirect
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.8.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vedhavyas/go-subkey v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/sys v0.1.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)

replace github.com/centrifuge/go-substrate-rpc-client/v4 v4.0.5 => github.com/threefoldtech/go-substrate-rpc-client/v4 v4.0.6-0.20220927094755-0f0d22c73cc7

replace github.com/threefoldtech/substrate-client-main v0.0.1 => github.com/threefoldtech/substrate-client v0.0.0-20230103083438-c6d86a2e6338
