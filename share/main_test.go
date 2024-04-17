package main

import (
	"fmt"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"os"
	"recovery-tool/cmd"
	"strconv"
	"strings"
	"testing"
)

// must remove import "C"
func TestGoRecovery(t *testing.T) {
	bytess, err := os.ReadFile("../input.yaml")
	if err != nil {
		t.Error(err)
	}

	var input cmd.RecoveryInput
	if err := yaml.UnmarshalStrict(bytess, &input); err != nil {
		t.Error(err)
	}
	input.ZipPath = "../test/134_archive.zip"
	vaultCountStr := strconv.Itoa(input.VaultCount)
	coinTypeStr := ""
	for _, c := range input.CoinType {
		coinTypeStr += strconv.Itoa(c) + ","
	}
	coinTypeStr = strings.TrimRight(coinTypeStr, ",")
	res := GoRecovery(input.ZipPath, input.UserMnemonic, input.EciesPrivKey, input.RsaPrivKey, vaultCountStr, coinTypeStr)
	fmt.Printf("res: %v \n", res)
}
