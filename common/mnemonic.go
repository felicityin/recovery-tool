package common

import (
	"fmt"
	"recovery-tool/crypto/hdwallet"
)

func CalcMasterPriv(menemonic string) (privKey, chainCode [32]byte, err error) {
	fmt.Println(menemonic)
	seed, err := hdwallet.NewSeed(menemonic, "", hdwallet.English)
	if err != nil {
		return privKey, chainCode, fmt.Errorf("create seed err: %s", err.Error())
	}

	privKey, chainCode = hdwallet.ComputeMastersFromSeed(seed)
	return privKey, chainCode, nil
}
