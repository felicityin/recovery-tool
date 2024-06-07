package main

//#include "file.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"os"
	"recovery-tool/cmd"
	"recovery-tool/common"
	"recovery-tool/common/code"
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
func GoRecovery(zipPath, userMnemonic, eciesPrivKey, rsaPrivKeyPath, vaultCount, chains, language string) C.RSResult {
	vaultCountInt, err := strconv.Atoi(vaultCount)
	if err != nil {
		return C.RSResult{
			errMsg: C.CString(code.GetMessage(language, code.VaultIndexParamErr, "RSA")),
			data:   C.CString(""),
			ok:     C.FALSE,
		}
	}

	rsaBytes, err := os.ReadFile(rsaPrivKeyPath)
	if err != nil {
		return C.RSResult{
			errMsg: C.CString(code.GetMessage(language, code.FileNotFound, "RSA")),
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
		var errMsg string
		i18nErr, ok := err.(*code.I18nError)
		if ok {
			errMsg = code.GetMessage(language, i18nErr.Code)
		} else {
			errMsg = err.Error()
		}
		return C.RSResult{
			errMsg: C.CString(errMsg),
			data:   C.CString(""),
			ok:     C.FALSE,
		}
	}

	resBytes, _ := json.Marshal(recoverResult)
	data := string(resBytes)
	msg := code.GetMessage(language, code.Success)
	return C.RSResult{
		errMsg: C.CString(msg),
		data:   C.CString(data),
		ok:     C.TRUE,
	}
}

//func GoRecoveryTest(zipPath, userMnemonic, eciesPrivKey, rsaPrivKeyPath, vaultCount, chains, language string) (err error) {
//	vaultCountInt, err := strconv.Atoi(vaultCount)
//	if err != nil {
//		return errors.New(code.ParamErrorMsg(language, code.VaultIndexParamErr))
//	}
//
//	rsaBytes, err := os.ReadFile(rsaPrivKeyPath)
//	if err != nil {
//		return errors.New(code.GetMessage(language, code.FileNotFound, "RSA"))
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
//		var errMsg string
//		i18nErr, ok := err.(*code.I18nError)
//		if ok {
//			errMsg = code.GetMessage(language, i18nErr.Code)
//		} else {
//			errMsg = err.Error()
//		}
//		return errors.New(errMsg)
//	}
//
//	resBytes, _ := json.Marshal(recoverResult)
//	data := string(resBytes)
//	fmt.Printf("data: %s", data)
//	return nil
//}

func main() {
}
