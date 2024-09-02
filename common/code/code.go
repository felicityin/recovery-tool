package code

import (
	"fmt"
	"strings"
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
	MnemonicNot24Words        = "510"
	ChainNameNotEmpty         = "511"
	VaultCountErr             = "512"
	ChainParamErr             = "513"
	EciesKeyNotEmpty          = "514"
	FileFormatErr             = "515"
	VaultIndexParamErr        = "516"
	RSAKeyNotEmpty            = "517"
	FailedToParseDataErr      = "518"

	PrivkeyIsHex           = "600"
	PrivkeyInvalid         = "601"
	DstAddrNotEmpty        = "602"
	SrcAddrNotEmpty        = "603" // get balance
	DstAddrInvalid         = "604"
	SrcAccountNotFound     = "605"
	DstAccountNotFound     = "606"
	SrcCoinAccountNotFound = "607"
	DstCoinAccountNotFound = "608"
	AmountInvalid          = "609"
	CoinUnsupported        = "610"
	SolInsufficientFunds   = "611"
	AptInsufficientFunds   = "612"
	DotInsufficientFunds   = "613"
	NetworkErr             = "614"
)

var I18nMessage = map[string]map[string]string{
	"en": {
		"fail_prefix":             "Recovery failed: ",
		"succ_prefix":             "Recovery successful",
		ParamErr:                  "Parameter error.",
		SystemErr:                 "System error.",
		FileNotFound:              "File not found.",
		MnemonicErr:               "Mnemonic error.",
		EciesPrivKeyErr:           "ECIES key error.",
		RsaPrivKeyErr:             "RSA key error.",
		MnemonicNotMatch:          "No backup data matching the mnemonic phrase.",
		RSADecryptBackupDataErr:   "RSA descryption of backup data failed.",
		EciesDecryptBackupDataErr: "ECIES descryption of backup data failed.",
		DeriveChildPrivErr:        "Child private key derivation failed.",
		DeriveChildAddressErr:     "Address derivation failed.",
		MnemonicNot24Words:        "Mnemonic must be 24 words.",
		ChainNameNotEmpty:         "Chain name should not be empty.",
		VaultCountErr:             "Wallet quantity must be greater or equal than 1.",
		ChainParamErr:             "Chain parameter error",
		EciesKeyNotEmpty:          "ECIES key should not be empty.",
		RSAKeyNotEmpty:            "RSA key should not be empty.",
		FileFormatErr:             "File format error.",
		VaultIndexParamErr:        "Vault index param error.",
		FailedToParseDataErr:      "Failed to parse backup data.",

		PrivkeyIsHex:           "The private key should be in hexadecimal format.",
		PrivkeyInvalid:         "The private key format is wrong, please re-enter.",
		DstAddrNotEmpty:        "The target address cannot be empty, please re-enter.",
		SrcAddrNotEmpty:        "The address cannot be empty, please re-enter.",
		DstAddrInvalid:         "The receiving address format is incorrect, please re-enter.",
		DstAccountNotFound:     "The receiving account does not exist, please check and try again",
		SrcAccountNotFound:     "The sending account does not exist, please check and try again",
		SrcCoinAccountNotFound: "The sending token address does not exist, please check and try again.",
		DstCoinAccountNotFound: "The receiving token address does not exist, please check and try again.",
		AmountInvalid:          "Please enter the correct amount.",
		CoinUnsupported:        "This chain does not support non-main chain currency.",
		SolInsufficientFunds:   "Insufficient gas fee (the current maximum transaction fee on the chain is 0.00089608 sol).",
		AptInsufficientFunds:   "Insufficient gas fee (the current maximum transaction fee on the chain is 0.002 apt).",
		DotInsufficientFunds:   "Insufficient gas fee (the current maximum transaction fee on the chain is 1 dot).",
		NetworkErr:             "Network error, please try again later.",
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
		MnemonicNot24Words:        "助记词必须为24个单词",
		ChainNameNotEmpty:         "链名不能为空",
		VaultCountErr:             "钱包数量必须大于等于1",
		ChainParamErr:             "链 参数错误",
		EciesKeyNotEmpty:          "ECIES 密钥不能为空",
		RSAKeyNotEmpty:            "RSA 密钥不能为空",
		FileFormatErr:             "文件格式错误",
		VaultIndexParamErr:        "钱包数量 参数错误",
		FailedToParseDataErr:      "解析备份数据失败",

		PrivkeyIsHex:           "私钥应该为 16 进制格式",
		PrivkeyInvalid:         "私钥格式错误，请重新填写",
		DstAddrNotEmpty:        "目标地址不能为空，请重新填写",
		SrcAddrNotEmpty:        "地址不能为空，请重新填写",
		DstAddrInvalid:         "接收地址格式不正确，请重新填写",
		DstAccountNotFound:     "接收账户不存在，请检查后重试",
		SrcAccountNotFound:     "发送账户不存在，请检查后重试",
		DstCoinAccountNotFound: "接收代币地址不存在，请检查后重试",
		SrcCoinAccountNotFound: "发送代币地址不存在，请检查后重试",
		AmountInvalid:          "请输入正确数量",
		CoinUnsupported:        "该链暂不支持非主链币发送",
		SolInsufficientFunds:   "网络费用不足（当前链上最大交易手续费为 0.00089608 sol）",
		AptInsufficientFunds:   "网络费用不足（当前链上最大交易手续费为 0.002 apt）",
		DotInsufficientFunds:   "网络费用不足（当前链上最大交易手续费为 1 dot）",
		NetworkErr:             "网络错误，请稍后重试",
	},
}

func GetMessage(language string, code string, arg ...string) string {
	info, ok := I18nMessage[language]
	if !ok {
		info = I18nMessage["en"]
	}

	var message string
	if message, ok = info[code]; !ok {
		message = info[SystemErr]
	}
	if code != Success {
		return fmt.Sprintf("%s%s%s", info["fail_prefix"], strings.Join(arg, ""), message)
	}
	return fmt.Sprintf("%s", info["succ_prefix"])
}

func GetMsg(language string, code string, arg ...string) string {
	info, ok := I18nMessage[language]
	if !ok {
		info = I18nMessage["en"]
	}

	var message string
	if message, ok = info[code]; !ok {
		message = info[SystemErr]
	}
	if code != Success {
		return message
	}
	return ""
}

func ParamErrorMsg(language, code string, arg ...string) string {
	info, ok := I18nMessage[language]
	if !ok {
		info = I18nMessage["en"]
	}

	var message string
	if message, ok = info[code]; !ok {
		message = info[SystemErr]
	}

	return fmt.Sprintf("%s%s", info["fail_prefix"], fmt.Sprintf(message, arg))
}
