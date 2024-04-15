package cmd

import (
	"archive/zip"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"recovery-tool/common"
	"recovery-tool/crypto"

	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/btcsuite/btcd/btcec"
	ecies "github.com/ecies/go/v2"
)

const (
	// 81/WalletType/VaultIndex/CoinType/AddressIndex
	AssetWalletPath = "81/0/%d/%d/0"
	ApiWalletPath   = "81/1/0/%d/%d"
)

type RecoveryInput struct {
	ZipPath         string `yaml:"zip_path"`
	UserMnemonic    string `yaml:"user_mnemonic"`
	UserMnemonicPwd string `yaml:"user_mnemonic_passwd"`
	EciesPrivKey    string `yaml:"ecies_private_key"`
	RsaPrivKey      string `yaml:"coincover_private_key"`
	VaultCount      int    `yaml:"valut_count"`
	CoinType        int    `yaml:"coin_type"`
}

type DeriveResult struct {
	VaultIndex int    `yaml:"vault_index"`
	CoinType   string `yaml:"coin_type"`
	Address    string `yaml:"address"`
	PrivKey    string `yaml:"private_key"`
}

var params RecoveryInput

func RecoverKeys(paramsPath string, outputPath string) error {
	err := loadRecoveryParams(paramsPath)
	if err != nil {
		common.Logger.Errorf("load params failed: %s", err)
		return err
	}

	userPrivKey, userChainCode, err := common.CalcMasterPriv(params.UserMnemonic)
	if err != nil {
		common.Logger.Errorf("calc user priv infos failed: %s", err)
		return err
	}
	usrPrivKey := new(big.Int).SetBytes(userPrivKey[:])

	eciesPrivKey, err := ecies.NewPrivateKeyFromHex(params.EciesPrivKey)
	if err != nil {
		common.Logger.Errorf("load ecies privkey failed: %s", err)
		return err
	}
	rsaPrivKey, err := crypto.ParseRsaPrivKey(params.RsaPrivKey)
	if err != nil {
		common.Logger.Errorf("parse rsa privkey failed: %s", err)
		return err
	}

	encryptedUserPubKey, err := encryptUsrPubKey(userPrivKey[:], eciesPrivKey, rsaPrivKey)
	if err != nil {
		common.Logger.Errorf("encrypt user privkey failed: %s", err)
		return err
	}
	hbcPrivs, err := findHbcPrivs(params.ZipPath, encryptedUserPubKey, eciesPrivKey, rsaPrivKey)
	if err != nil {
		common.Logger.Errorf("find hbc private info failed: %s", err)
		return err
	}

	privs := &common.RootKeys{
		HbcShare0: hbcPrivs[0],
		HbcShare1: hbcPrivs[1],
		UsrShare: &common.RootKey{
			PrivKey:   usrPrivKey,
			PubKey:    calcPubKey(usrPrivKey),
			ChainCode: userChainCode[:],
		},
	}
	childs, err := deriveChilds(params.VaultCount, params.CoinType, privs)
	if err != nil {
		common.Logger.Errorf("derive childs failed: %s", err)
		return err
	}

	yamlData, err := yaml.Marshal(&childs)
	if err != nil {
		common.Logger.Errorf("yaml marshal result failed: %s", err)
		return err
	}
	err = ioutil.WriteFile(outputPath, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}

	return nil
}

func loadRecoveryParams(path string) error {
	bytess, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = yaml.UnmarshalStrict(bytess, &params)
	return err
}

func encryptUsrPubKey(userPrivKey []byte, eciesPrivKey *ecies.PrivateKey, rsaPrivKey *rsa.PrivateKey) (string, error) {
	userPubKey := crypto.ScalarBaseMult(crypto.S256(), new(big.Int).SetBytes(userPrivKey))
	userPublicKey := btcec.PublicKey{Curve: btcec.S256(), X: userPubKey.X(), Y: userPubKey.Y()}
	userPublicKeyBytes := userPublicKey.SerializeCompressed()

	encryptedUsrPubKey, err := ecies.Encrypt(eciesPrivKey.PublicKey, userPublicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("ecies Encrypt userPublicKey failed, err: %s", err.Error())
	}

	encryptedUsrPubKey, err = crypto.RsaEncryptOAEP(&rsaPrivKey.PublicKey, encryptedUsrPubKey)
	if err != nil {
		return "", fmt.Errorf("rsa Encrypt userPublicKey failed, err: %s", err.Error())
	}

	return hex.EncodeToString(encryptedUsrPubKey), nil
}

type encryptedTeam struct {
	HbcPrivKeys   []string `json:"hbc_private_keys"`
	HbcChainCodes []string `json:"hbc_chain_codes"`
	UserPubKey    string   `json:"user_pub_key"`
}

func findHbcPrivs(
	zipPath string,
	encryptedUserPubKey string,
	eciesPrivKey *ecies.PrivateKey,
	rsaPrivKey *rsa.PrivateKey,
) ([]*common.RootKey, error) {
	zf, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zf.Close()

	var result []*common.RootKey

	for _, file := range zf.File {
		fileBytes, err := common.ReadAll(file)
		if err != nil {
			return nil, err
		}

		encrypted := encryptedTeam{}
		err = json.Unmarshal(fileBytes, &encrypted)
		if err != nil {
			return nil, fmt.Errorf("unmarshal team failed: %s", err.Error())
		}

		if encrypted.UserPubKey == encryptedUserPubKey {
			priv0, err := getHbcPriv(encrypted.HbcPrivKeys[0], encrypted.HbcChainCodes[0], eciesPrivKey, rsaPrivKey)
			if err != nil {
				return nil, err
			}
			priv1, err := getHbcPriv(encrypted.HbcPrivKeys[1], encrypted.HbcChainCodes[1], eciesPrivKey, rsaPrivKey)
			if err != nil {
				return nil, err
			}

			result = append(result, priv0)
			result = append(result, priv1)
			return result, nil
		}
	}
	return nil, fmt.Errorf("mnemonic or zip is invalid")
}

func getHbcPriv(
	privKey, chainCode string,
	eciesPrivKey *ecies.PrivateKey,
	rsaPrivKey *rsa.PrivateKey,
) (*common.RootKey, error) {
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, fmt.Errorf("hex decode privkey failed: %s", err.Error())
	}
	chainCodeBytes, err := hex.DecodeString(chainCode)
	if err != nil {
		return nil, fmt.Errorf("hex decode chaincode failed: %s", err.Error())
	}

	decryptedPrivKey, err := crypto.RSADecryptOAEP(rsaPrivKey, privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("rsa decode privkey failed: %s", err.Error())
	}
	decryptedChainCode, err := crypto.RSADecryptOAEP(rsaPrivKey, chainCodeBytes)
	if err != nil {
		return nil, fmt.Errorf("rsa decode chaincode failed: %s", err.Error())
	}

	decryptedPrivKey, err = ecies.Decrypt(eciesPrivKey, decryptedPrivKey)
	if err != nil {
		return nil, fmt.Errorf("ecies decode privkey failed: %s", err.Error())
	}
	decryptedChainCode, err = ecies.Decrypt(eciesPrivKey, decryptedChainCode)
	if err != nil {
		return nil, fmt.Errorf("ecies decode chaincode failed: %s", err.Error())
	}

	privateKey := new(big.Int).SetBytes(decryptedPrivKey)

	return &common.RootKey{
		PrivKey:   privateKey,
		PubKey:    calcPubKey(privateKey),
		ChainCode: decryptedChainCode,
	}, nil
}

func calcPubKey(privKey *big.Int) *ecdsa.PublicKey {
	pubPoint := crypto.ScalarBaseMult(btcec.S256(), privKey)
	return &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     pubPoint.X(),
		Y:     pubPoint.Y(),
	}
}

func deriveChilds(vaultCount int, coinType int, rootKeys *common.RootKeys) ([]*DeriveResult, error) {
	var deriveResult []*DeriveResult

	for vaultIndex := 0; vaultIndex < vaultCount; vaultIndex++ {
		hdPath := fmt.Sprintf(AssetWalletPath, vaultIndex, coinType) // Only support asset wallet for now
		privKey, address, err := common.DeriveChild(rootKeys, hdPath)
		if err != nil {
			return nil, fmt.Errorf("derive child failed, err: %s", err.Error())
		}

		var buf [32]byte
		privKeyBytes := privKey.FillBytes(buf[:])

		deriveResult = append(deriveResult, &DeriveResult{
			VaultIndex: vaultIndex,
			CoinType:   common.SwitchChain(uint32(coinType)),
			Address:    address,
			PrivKey:    hex.EncodeToString(privKeyBytes),
		})
	}

	return deriveResult, nil
}
