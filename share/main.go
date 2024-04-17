package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
*/

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"os"
	"recovery-tool/cmd"
	"strconv"
	"strings"
)

type RSResult struct {
	Success bool                `json:"success"`
	ErrMsg  string              `json:"errMsg"`
	Data    []*cmd.DeriveResult `json:"data"`
}

func (r *RSResult) ToJson() string {
	b, _ := json.Marshal(r)
	return string(b)
}

//export PrintHello
func PrintHello(name string) string {
	return "Hello " + name
}

//export GoRecovery
func GoRecovery(zipPath, userMnemonic, eciesPrivKey, rsaPrivKey, vaultCount, coinTypes string) string {
	rs := &RSResult{}

	vaultCountInt, err := strconv.Atoi(vaultCount)
	if err != nil {
		rs.Success = false
		rs.ErrMsg = "vaultCount not int"
		return rs.ToJson()
	}

	coinTypesList := strings.Split(coinTypes, ",")
	coinTypesSlice := make([]int, len(coinTypesList))
	for i, coinType := range coinTypesList {
		coinTypeInt, err := strconv.Atoi(coinType)
		if err != nil {
			rs.Success = false
			rs.ErrMsg = "coinTypes invalid"
			return rs.ToJson()
		}
		coinTypesSlice[i] = coinTypeInt
	}

	input := cmd.RecoveryInput{
		ZipPath:      zipPath,
		UserMnemonic: userMnemonic,
		EciesPrivKey: eciesPrivKey,
		RsaPrivKey:   rsaPrivKey,
		VaultCount:   vaultCountInt,
		CoinType:     coinTypesSlice,
	}

	recoverResult, err := cmd.RecoverKeys(input)
	if err != nil {
		rs.Success = false
		rs.ErrMsg = fmt.Sprintf("recover fail: %s", err.Error())
		return rs.ToJson()
	}

	rs.Success = true
	rs.Data = recoverResult
	return rs.ToJson()

}

func main() {
	//test
	bytess, err := os.ReadFile("../input.yaml")
	if err != nil {
		panic(err)
	}

	var input cmd.RecoveryInput
	if err := yaml.UnmarshalStrict(bytess, &input); err != nil {
		panic(err)
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
