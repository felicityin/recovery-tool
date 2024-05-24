package main

//#include <file.h>
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

//export GetKey
func GetKey(youKey string) *C.char {
	fmt.Println("hello, dart")
	theKey := "This is input key: " + youKey
	return C.CString(theKey)
}

//export SumTest
func SumTest(a, b int) C.int {
	return C.int(a + b)
}

//export GetRSResult
func GetRSResult(s string) C.RSResult {
	return C.RSResult{
		errMsg: C.CString(""),
		data:   C.CString(s),
		ok:     C.TRUE,
	}
}

//export GoRecovery
func GoRecovery(zipPath, userMnemonic, eciesPrivKey, rsaPrivKey, vaultCount, coinTypes string) C.RSResult {
	vaultCountInt, err := strconv.Atoi(vaultCount)
	if err != nil {
		return C.RSResult{
			errMsg: C.CString(err.Error()),
			data:   C.CString(""),
			ok:     C.FALSE,
		}
	}

	coinTypesList := strings.Split(coinTypes, ",")
	coinTypesSlice := make([]int, len(coinTypesList))
	for i, coinType := range coinTypesList {
		coinTypeInt, err := strconv.Atoi(coinType)
		if err != nil {
			return C.RSResult{
				errMsg: C.CString("coinTypes invalid"),
				data:   C.CString(""),
				ok:     C.FALSE,
			}
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
		return C.RSResult{
			errMsg: C.CString(fmt.Sprintf("recover fail: %s", err.Error())),
			data:   C.CString(""),
			ok:     C.FALSE,
		}
	}

	resBytes, _ := json.Marshal(recoverResult)
	data := string(resBytes)
	return C.RSResult{
		errMsg: C.CString(""),
		data:   C.CString(data),
		ok:     C.TRUE,
	}
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
