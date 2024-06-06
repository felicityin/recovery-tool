package code

import (
	"fmt"
)

var (
	Success                   = "200" //成功
	ParamErr                  = "499" //参数错误
	SystemErr                 = "500" //系统错误
	FileNotFound              = "501" //文件找不到
	MnemonicErr               = "502" //助记词错误
	EciesPrivKeyErr           = "503" //ecies私钥错误
	RsaPrivKeyErr             = "504" //rsa私钥错误
	MnemonicNotMatch          = "505" //没有与当前助记词匹配的备份数据
	RSADecryptBackupDataErr   = "506" //RSA解密备份数据失败
	EciesDecryptBackupDataErr = "507" //Ecies解密备份数据失败
	DeriveChildPrivErr        = "508" //子私钥推导失败
	DeriveChildAddressErr     = "509" //地址推导失败
)

var I18nMessage = map[string]map[string]string{
	"en": {
		"fail_prefix":             "Recovery failed: ",
		"succ_prefix":             "Recovery Successful",
		ParamErr:                  "Parameter error.",
		SystemErr:                 "System error.",
		FileNotFound:              "File not found.",
		MnemonicErr:               "Mnemonic error.",
		EciesPrivKeyErr:           "ECIES key error.",
		RsaPrivKeyErr:             "RSA key error.",
		MnemonicNotMatch:          "No backup data matching the mnemonic was found.",
		RSADecryptBackupDataErr:   "RSA decryption of backup data failed.",
		EciesDecryptBackupDataErr: "ECIES decryption of backup data failed.",
		DeriveChildPrivErr:        "Sub-private key derivation failed.",
		DeriveChildAddressErr:     "Address derivation failed.",
	},
	"zh": {
		"fail_prefix":             "恢复失败：",
		"succ_prefix":             "恢复成功",
		ParamErr:                  "参数错误",
		SystemErr:                 "系统错误",
		FileNotFound:              "文件找不到",
		MnemonicErr:               "助记词错误",
		EciesPrivKeyErr:           "Ecies密钥错误",
		RsaPrivKeyErr:             "RSA密钥错误",
		MnemonicNotMatch:          "没有找到与助记词匹配的备份数据",
		RSADecryptBackupDataErr:   "RSA解密备份数据失败",
		EciesDecryptBackupDataErr: "ECIES解密备份数据失败",
		DeriveChildPrivErr:        "子私钥推导失败",
		DeriveChildAddressErr:     "地址推导失败",
	},
}

func GetMessage(language string, code string) string {
	info, ok := I18nMessage[language]
	if !ok {
		info = I18nMessage["en"]
	}

	var message string
	if message, ok = info[code]; !ok {
		message = info[SystemErr]
	}
	if code != Success {
		return fmt.Sprintf("%s%s", info["fail_prefix"], message)
	}
	return fmt.Sprintf("%s%s", info["succ_prefix"], message)
}
