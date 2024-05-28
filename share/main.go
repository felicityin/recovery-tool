package main

//#include "file.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"os"
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
func GoRecovery(zipPath, userMnemonic, eciesPrivKey, rsaPrivKeyPath, vaultCount, chains string) C.RSResult {
	vaultCountInt, err := strconv.Atoi(vaultCount)
	if err != nil {
		return C.RSResult{
			errMsg: C.CString("vaultCount must be a number"),
			data:   C.CString(""),
			ok:     C.FALSE,
		}
	}

	rsaBytes, err := os.ReadFile(rsaPrivKeyPath)
	if err != nil {
		return C.RSResult{
			errMsg: C.CString(fmt.Sprintf("RSA file read failed: %s", err.Error())),
			data:   C.CString(""),
			ok:     C.FALSE,
		}
	}

	chainList := strings.Split(chains, ",")
	input := cmd.RecoveryInput{
		ZipPath:      zipPath,
		UserMnemonic: userMnemonic,
		EciesPrivKey: eciesPrivKey,
		RsaPrivKey:   string(rsaBytes),
		VaultCount:   vaultCountInt,
		Chains:       chainList,
	}

	recoverResult, err := cmd.RecoverKeys(input)
	if err != nil {
		return C.RSResult{
			errMsg: C.CString(fmt.Sprintf("Recover failed: %s", err.Error())),
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

//func GoRecoveryTest(zipPath, userMnemonic, eciesPrivKey, rsaPrivKeyPath, vaultCount, chains string) (err error) {
//	vaultCountInt, err := strconv.Atoi(vaultCount)
//	if err != nil {
//		return err
//	}
//
//	rsaBytes, err := os.ReadFile(rsaPrivKeyPath)
//	if err != nil {
//		return err
//	}
//
//	chainList := strings.Split(chains, ",")
//	input := cmd.RecoveryInput{
//		ZipPath:      zipPath,
//		UserMnemonic: userMnemonic,
//		EciesPrivKey: eciesPrivKey,
//		RsaPrivKey:   string(rsaBytes),
//		VaultCount:   vaultCountInt,
//		Chains:       chainList,
//	}
//
//	recoverResult, err := cmd.RecoverKeys(input)
//	if err != nil {
//		return err
//	}
//
//	resBytes, _ := json.Marshal(recoverResult)
//	data := string(resBytes)
//	fmt.Printf("data: %s", data)
//	return nil
//}

func main() {

}
