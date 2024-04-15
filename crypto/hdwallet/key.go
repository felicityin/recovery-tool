package hdwallet

import (
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/blocktree/go-owcdrivers/owkeychain"
	"github.com/blocktree/go-owcrypt"
)

// Key struct
type Key struct {
	opt *Options

	Mnemonic string
	Seed     string
	Net      string
}

func NewKey(opts ...Option) (*Key, error) {
	var (
		err error
		o   = newOptions(opts...)
	)
	if len(o.Mnemonic) > 0 {
		mm := strings.Replace(o.Mnemonic, " ", "", -1)

		if ok, _ := regexp.MatchString(`[\p{Hangul}]`, mm); ok {
			o.Language = "korean"

		}
		if ok, _ := regexp.MatchString(`[a-zA-Z]`, mm); ok {
			o.Language = "english"
		}
		if ok, _ := regexp.MatchString(`[\p{Han}]`, mm); ok {
			o.Language = "chinese_simplified"
		}
	}
	if len(o.Seed) <= 0 {
		o.Seed, err = NewSeed(o.Mnemonic, o.Password, o.Language)
	}
	if err != nil {
		return nil, err
	}

	key := &Key{
		opt:      o,
		Mnemonic: o.Mnemonic,
		Seed:     hex.EncodeToString(o.Seed),
		Net:      o.Net,
	}

	err = key.init()
	if err != nil {
		return nil, err
	}

	return key, nil
}
func (k *Key) init() error {
	return nil
}

// GetChildKey return a key from master key
// params: [Purpose], [CoinType], [Account], [Change], [AddressIndex], [Path]

func DerivePathFromSeed(masterSeed []byte, path string) (prv []byte, err error) {
	pkey, err := owkeychain.DerivedPrivateKeyWithPath(masterSeed, path, owcrypt.ECC_CURVE_SECP256K1)
	if err != nil {
		return nil, err
	} else {
		return pkey.GetPrivateKeyBytes()
	}
}
