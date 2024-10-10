package address

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/blocktree/go-owcrypt/eddsa"
)

func Test_Address(t *testing.T) {
	//parse_hex := "deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363"

	parse_hex := "612d82bc053d1b4729057688ecb1ebf62745d817ddd9b595bc822f5f2ba0e41a"
	sk, _ := hex.DecodeString(parse_hex)
	pubkey, _ := eddsa.ED25519_genPub(sk)
	pubkey, _ = hex.DecodeString("deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363")

	//fmt.Println("pubkey = ", hex.EncodeToString(pubkey))
	fmt.Println("address = ", PublicKeyToAddress(pubkey, "", 0))

	address, _ := AddressEncode(pubkey, 0)
	fmt.Println("address = ", address)

}
