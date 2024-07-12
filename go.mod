module recovery-tool

go 1.20

replace github.com/btcsuite/btcutil/hdkeychain v0.0.0-20191219182022-e17c9730c422 => github.com/btcsuite/btcd/btcuti/hdkeychain v1.1.3

replace github.com/btcsuite/btcd/btcec => ./package/github.com/btcsuite/btcd/btcec/v1

replace github.com/btcsuite/btcd/btcec/v2 => ./package/github.com/btcsuite/btcd/btcec/v2

require (
	github.com/HcashOrg/bliss v0.0.0-20180719035130-f5d53c2a9b7d // indirect
	github.com/HcashOrg/hcd v0.0.0-20180816055255-f68c5e6e35cb
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/btcsuite/btcd/btcec v0.0.0-00010101000000-000000000000
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	golang.org/x/crypto v0.16.0
)

require (
	filippo.io/edwards25519 v1.0.0-rc.1
	github.com/alecthomas/gometalinter v3.0.0+incompatible
	github.com/blocktree/go-owcdrivers v1.2.27
	github.com/blocktree/go-owcrypt v1.1.14
	github.com/blocto/solana-go-sdk v1.30.0
	github.com/bnb-chain/edwards25519 v0.0.0-20231030070956-6796d47b70ba
	github.com/btcsuite/btcd v0.23.0
	github.com/btcsuite/btcd/btcec/v2 v2.2.0
	github.com/btcsuite/btcd/btcutil v1.1.3
	github.com/decred/dcrd/dcrec/edwards/v2 v2.0.3
	github.com/ecies/go/v2 v2.0.9
	github.com/ethereum/go-ethereum v1.13.5
	github.com/fbsobreira/gotron-sdk v0.0.0-20230907131216-1e824406fe8c
	github.com/imroc/req v0.2.3
	github.com/ipfs/go-log v1.0.5
	github.com/mr-tron/base58 v1.2.0
	github.com/near/borsh-go v0.3.2-0.20220516180422-1ff87d108454
	github.com/pkg/errors v0.9.1
	github.com/portto/aptos-go-sdk v0.0.0-20230807103729-9a5201cad72f
	github.com/shopspring/decimal v1.4.0
	github.com/stretchr/testify v1.8.4
	github.com/the729/lcs v0.1.5
	github.com/tidwall/gjson v1.2.1
)

require (
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/blake256 v1.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/drand/kyber v1.1.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hasura/go-graphql-client v0.9.1 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/ipfs/go-log/v2 v2.1.3 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/phoreproject/bls v0.0.0-20200525203911-a88a5ae26844 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shengdoushi/base58 v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
