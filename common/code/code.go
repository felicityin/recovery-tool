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

	PrivkeyIsHex         = "600"
	InvalidPrivkey       = "601"
	DstAddrNotEmpty      = "602"
	AddrNotEmpty         = "603"
	InvalidToAddr        = "604"
	AccountNotFound      = "605"
	AmountInvalid        = "606"
	CoinUnsupported      = "607"
	SolInsufficientFunds = "608"
	AptInsufficientFunds = "609"
	DotInsufficientFunds = "610"
	NetworkErr           = "611"
	FromBalanceZero      = "612"
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

		PrivkeyIsHex:         "The private key should be in hexadecimal format",
		DstAddrNotEmpty:      "The recipient's address cannot be empty",
		InvalidToAddr:        "Dest address is invalid",
		InvalidPrivkey:       "Private key is invalid",
		AmountInvalid:        "Unable to convert transfer amount to decimal",
		CoinUnsupported:      "This chain only supports main chain coin for now",
		AddrNotEmpty:         "The address cannot be empty",
		SolInsufficientFunds: "Insufficient balance to pay for transaction fee. The max tx fee is 0.00089608 sol",
		AptInsufficientFunds: "Insufficient balance to pay for transaction fee. The max tx fee is 0.002 apt",
		DotInsufficientFunds: "Insufficient balance to pay for transaction fee. The max tx fee is 1 dot",
		NetworkErr:           "Network error. Please retry later!",
		FromBalanceZero:      "The balance of the coin contract is 0",
		AccountNotFound:      "Account not found",
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

		PrivkeyIsHex:         "私钥应该为 16 进制格式",
		InvalidPrivkey:       "私钥格式不正确",
		DstAddrNotEmpty:      "目标地址不能为空",
		InvalidToAddr:        "目标地址格式不正确",
		AddrNotEmpty:         "地址不能为空",
		AccountNotFound:      "账户不存在",
		AmountInvalid:        "转账金额无法转换为小数",
		CoinUnsupported:      "该链暂时只支持主链币",
		SolInsufficientFunds: "余额不足以支付交易手续费，最大交易手续费为 0.00089608 sol",
		AptInsufficientFunds: "余额不足以支付交易手续费，最大交易手续费为 0.002 apt",
		DotInsufficientFunds: "余额不足以支付交易手续费，最大交易手续费为 1 dot",
		NetworkErr:           "网络错误，请稍后重试",
		FromBalanceZero:      "合约币余额为 0",
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
