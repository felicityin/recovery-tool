module recovery-tool

go 1.20

replace github.com/btcsuite/btcutil/hdkeychain v0.0.0-20191219182022-e17c9730c422 => github.com/btcsuite/btcd/btcuti/hdkeychain v1.1.3

replace github.com/btcsuite/btcd/btcec => ./package/github.com/btcsuite/btcd/btcec/v1

replace github.com/btcsuite/btcd/btcec/v2 => ./package/github.com/btcsuite/btcd/btcec/v2

require (
	github.com/HcashOrg/bliss v0.0.0-20180719035130-f5d53c2a9b7d // indirect
	github.com/HcashOrg/hcd v0.0.0-20180816055255-f68c5e6e35cb
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/btcsuite/btcd/btcec v0.0.0-00010101000000-000000000000
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.16.0
)

require (
	github.com/alecthomas/gometalinter v3.0.0+incompatible
	github.com/blocktree/go-owcdrivers v1.2.27
	github.com/blocktree/go-owcrypt v1.1.14
	github.com/btcsuite/btcd v0.23.0
	github.com/btcsuite/btcd/btcec/v2 v2.2.0
	github.com/btcsuite/btcd/btcutil v1.1.3
	github.com/decred/dcrd/dcrec/edwards/v2 v2.0.3
	github.com/ecies/go/v2 v2.0.9
	github.com/ethereum/go-ethereum v1.13.5
	github.com/fbsobreira/gotron-sdk v0.0.0-20230907131216-1e824406fe8c
)

require (
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 // indirect
	github.com/dchest/blake256 v1.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/drand/kyber v1.1.4 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/phoreproject/bls v0.0.0-20200525203911-a88a5ae26844 // indirect
	github.com/shengdoushi/base58 v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
