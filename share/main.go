package main

//#include "file.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"recovery-tool/cmd"
	"recovery-tool/common"
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

//export GetChainList
func GetChainList() *C.char {
	chainList, _ := json.Marshal(common.ChainList)
	return C.CString(string(chainList))
}

//export GoRecovery
func GoRecovery(zipPath, userMnemonic, eciesPrivKey, rsaPrivKey, vaultCount, chains string) C.RSResult {
	vaultCountInt, err := strconv.Atoi(vaultCount)
	if err != nil {
		return C.RSResult{
			errMsg: C.CString(err.Error()),
			data:   C.CString(""),
			ok:     C.FALSE,
		}
	}

	chainList := strings.Split(chains, ",")
	input := cmd.RecoveryInput{
		ZipPath:      zipPath,
		UserMnemonic: userMnemonic,
		EciesPrivKey: eciesPrivKey,
		RsaPrivKey:   rsaPrivKey,
		VaultCount:   vaultCountInt,
		Chains:       chainList,
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

}
