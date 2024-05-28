package main

import (
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"os"
	"recovery-tool/cmd"
	"strconv"
	"strings"
	"testing"
)

// must remove import "C"
func TestGoRecoveryTest(t *testing.T) {
	bytess, err := os.ReadFile("./input1.yaml")
	if err != nil {
		panic(err)
	}

	var input cmd.RecoveryInput
	if err := yaml.UnmarshalStrict(bytess, &input); err != nil {
		panic(err)
	}

	vaultCountStr := strconv.Itoa(input.VaultCount)
	chainStr := ""
	for _, chainName := range input.Chains {
		chainStr += chainName + ","
	}
	chainStr = strings.TrimRight(chainStr, ",")
	GoRecoveryTest(input.ZipPath, input.UserMnemonic, input.EciesPrivKey, "./test/private_f5a4b26f3c2231dec42ff8c4ade8530c.key", vaultCountStr, chainStr)
}
