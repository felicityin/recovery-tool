package cmd

import (
	"archive/zip"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"recovery-tool/common/code"
	"strings"

	"github.com/HcashOrg/hcd/dcrec/edwards"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/btcsuite/btcd/btcec"
	ecies "github.com/ecies/go/v2"
	"math"
	"recovery-tool/common"
	"recovery-tool/crypto"
	"sort"
	"sync"
)

const (
	// 81/WalletType/VaultIndex/CoinType/AddressIndex
	AssetWalletPath = "81/0/%d/%d/0"
	ApiWalletPath   = "81/1/0/%d/%d"
)

type RecoveryInput struct {
	ZipPath      string   `yaml:"zip_path"`
	UserMnemonic string   `yaml:"user_mnemonic"`
	EciesPrivKey string   `yaml:"ecies_private_key"`
	RsaPrivKey   string   `yaml:"rsa_private_key"`
	VaultCount   int      `yaml:"valut_count"`
	CoinType     []int    `yaml:"coin_type"`
	Chains       []string `yaml:"chains"`
}

type DeriveResult struct {
	VaultIndex int    `yaml:"vault_index"`
	Chain      string `yaml:"chain"`
	Address    string `yaml:"address"`
	PrivKey    string `yaml:"private_key"`
}

type parsedParams struct {
	UserPrivKeyScalar *big.Int
	UserChainCode     []byte
	UserPubKey        string
	EciesPrivKey      *ecies.PrivateKey
	RsaPrivKey        *rsa.PrivateKey
}

func RecoverKeysCmd(paramsPath string, outputPath string) error {
	params := loadRecoveryParams(paramsPath)

	result, err := RecoverKeys(params)
	if err != nil {
		common.Logger.Errorf("derive keys failed")
		return err
	}

	if err = SaveResult(&result, outputPath); err != nil {
		common.Logger.Errorf("save result failed")
		return err
	}
	return nil
}

func RecoverKeys(params RecoveryInput) ([]*DeriveResult, error) {
	if err := checkParams(params); err != nil {
		return nil, code.NewI18nError(code.ParamErr, err.Error())
	}

	parsed, err := parseParams(params)
	if err != nil {
		return nil, err
	}

	hbcPrivs, err := findHbcPrivs(params.ZipPath, parsed.UserPubKey, parsed.EciesPrivKey, parsed.RsaPrivKey)
	if err != nil {
		common.Logger.Errorf("find hbc private info failed: %s", err)
		return nil, err
	}

	eddsaPrivKey := new(big.Int).SetBytes(hbcPrivs[0].PrivKey.Bytes())
	eddsaPrivKey.Add(eddsaPrivKey, hbcPrivs[1].PrivKey)
	eddsaPrivKey.Mod(eddsaPrivKey, edwards.Edwards().N)
	eddsaPrivKey.Add(eddsaPrivKey, parsed.UserPrivKeyScalar)
	eddsaPrivKey.Mod(eddsaPrivKey, edwards.Edwards().N)

	privs := &common.RootKeys{
		HbcShare0: hbcPrivs[0],
		HbcShare1: hbcPrivs[1],
		UsrShare: &common.RootKey{
			PrivKey:     parsed.UserPrivKeyScalar,
			EcdsaPubKey: crypto.ScalarBaseMult(crypto.S256(), parsed.UserPrivKeyScalar),
			EddsaPubKey: crypto.ScalarBaseMult(crypto.Edwards(), parsed.UserPrivKeyScalar),
			ChainCode:   parsed.UserChainCode[:],
		},
		// EddsaPubKey: pubKey,
		EddsaPubKey: crypto.ScalarBaseMult(crypto.Edwards(), eddsaPrivKey),
	}

	keys, err := concurrentDeriveChilds(params.VaultCount, params.Chains, privs)
	if err != nil {
		common.Logger.Errorf("derive childs failed: %s", err)
		return nil, err
	}
	return keys, nil
}

func loadRecoveryParams(path string) RecoveryInput {
	bytess, err := ioutil.ReadFile(path)
	if err != nil {
		common.Logger.Errorf("load params error: %s", err.Error())
		panic(err)
	}

	var params RecoveryInput
	if err = yaml.UnmarshalStrict(bytess, &params); err != nil {
		common.Logger.Errorf("unmarshal params error: %s", err.Error())
		panic(err)
	}
	return params
}

func checkParams(params RecoveryInput) (err error) {
	if len(params.ZipPath) <= 0 {
		return fmt.Errorf("zip_path cannot be empty")
	}

	userMnemonics := strings.Split(params.UserMnemonic, " ")
	if len(userMnemonics) != 24 {
		return fmt.Errorf("mnemonic word not 24 words")
	}
	for i, word := range userMnemonics {
		userMnemonics[i] = strings.TrimSpace(word)
	}

	if len(params.EciesPrivKey) <= 0 {
		return fmt.Errorf("ECIES key cannot be empty")
	}

	if len(params.RsaPrivKey) <= 0 {
		return fmt.Errorf("RSA key cannot be empty")
	}

	if params.VaultCount <= 0 {
		return fmt.Errorf("VaultCount must > 1")
	}

	if len(params.CoinType) <= 0 && len(params.Chains) <= 0 {
		return fmt.Errorf("coin_type and chains cannot be empty")
	}

	if len(params.Chains) > 0 {
		chainMap := make(map[string]struct{})
		for _, chainName := range params.Chains {
			if _, ok := common.ChainInfos[chainName]; !ok {
				return fmt.Errorf("unsupported chain: %s", chainName)
			}
			chainMap[chainName] = struct{}{}
		}
		chains := make([]string, len(chainMap))
		i := 0
		for chainName, _ := range chainMap {
			chains[i] = strings.TrimSpace(chainName)
			i++
		}
		params.Chains = chains
	}

	return nil
}

func parseParams(params RecoveryInput) (*parsedParams, error) {
	userPrivKey, userChainCode, err := common.CalcMasterPriv(params.UserMnemonic)
	if err != nil {
		common.Logger.Errorf("calc user priv infos failed: %s", err)
		return nil, code.NewI18nError(code.MnemonicErr, err.Error())
	}
	usrPrivKeyScalar := new(big.Int).SetBytes(userPrivKey[:])
	userPubKey := calcUserPubKey(usrPrivKeyScalar)
	common.Logger.Debugf("user pubkey: %s", userPubKey)

	eciesPrivKey, err := ecies.NewPrivateKeyFromHex(params.EciesPrivKey)
	if err != nil {
		common.Logger.Errorf("load ecies privkey failed: %s", err)
		return nil, code.NewI18nError(code.EciesPrivKeyErr, err.Error())
	}
	common.Logger.Debugf("ecies privkey: %d", eciesPrivKey.D)

	rsaPrivKey, err := crypto.ParseRsaPrivKey(params.RsaPrivKey)
	if err != nil {
		common.Logger.Errorf("parse rsa privkey failed: %s", err)
		return nil, code.NewI18nError(code.RsaPrivKeyErr, err.Error())
	}

	common.Logger.Debugf("rsa privkey: %d, %d", rsaPrivKey.Primes[0], rsaPrivKey.Primes[1])

	return &parsedParams{
		UserPrivKeyScalar: usrPrivKeyScalar,
		UserChainCode:     userChainCode[:],
		UserPubKey:        userPubKey,
		EciesPrivKey:      eciesPrivKey,
		RsaPrivKey:        rsaPrivKey,
	}, nil
}

func SaveResult(childs *[]*DeriveResult, outputPath string) error {
	yamlData, err := yaml.Marshal(childs)
	if err != nil {
		common.Logger.Errorf("yaml marshal result failed: %s", err)
		return err
	}
	err = ioutil.WriteFile(outputPath, yamlData, 0644)
	if err != nil {
		common.Logger.Errorf("unable to write data into the file")
		return err
	}
	return nil
}

func calcUserPubKey(privKey *big.Int) string {
	pubKey := crypto.ScalarBaseMult(crypto.S256(), privKey)
	pubKeyPoint := btcec.PublicKey{Curve: btcec.S256(), X: pubKey.X(), Y: pubKey.Y()}
	pubKeyBytes := pubKeyPoint.SerializeCompressed()
	return hex.EncodeToString(pubKeyBytes)
}

type encryptedTeam struct {
	HbcPrivKeys   []string `json:"hbc_private_keys"`
	HbcChainCodes []string `json:"hbc_chain_codes"`
	UserPubKey    string   `json:"user_pub_key"`
}

func findHbcPrivs(
	zipPath string,
	userPubKey string,
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

		decryptedUsrPubKey, err := decryptUserPubKey(encrypted.UserPubKey, eciesPrivKey, rsaPrivKey)
		common.Logger.Debugf("decrypted user pubkey: %s", decryptedUsrPubKey)
		if err != nil {
			return nil, err
		}

		if decryptedUsrPubKey == userPubKey {
			priv0, err := decryptHbcPriv(encrypted.HbcPrivKeys[0], encrypted.HbcChainCodes[0], eciesPrivKey, rsaPrivKey)
			if err != nil {
				return nil, err
			}
			priv1, err := decryptHbcPriv(encrypted.HbcPrivKeys[1], encrypted.HbcChainCodes[1], eciesPrivKey, rsaPrivKey)
			if err != nil {
				return nil, err
			}

			result = append(result, priv0)
			result = append(result, priv1)
			return result, nil
		}
	}
	return nil, code.NewI18nError(code.MnemonicNotMatch, "mnemonic and zip do not match")
}

func decryptUserPubKey(userPubKey string, eciesPrivKey *ecies.PrivateKey, rsaPrivKey *rsa.PrivateKey) (string, error) {
	userPubKeyBytes, err := hex.DecodeString(userPubKey)
	if err != nil {
		return "", fmt.Errorf("hex decode user pubkey error: %s", err.Error())
	}

	decryptedUsrPubKey, err := crypto.RsaDecryptOAEP(rsaPrivKey, userPubKeyBytes)
	if err != nil {
		return "", code.NewI18nError(code.RSADecryptBackupDataErr, fmt.Sprintf("rsa decrypt user pubkey error: %s", err.Error()))
	}

	decryptedUsrPubKey, err = ecies.Decrypt(eciesPrivKey, decryptedUsrPubKey)
	if err != nil {
		return "", code.NewI18nError(code.EciesDecryptBackupDataErr, fmt.Sprintf("ecies decrypt user pubkey error: %s", err.Error()))
	}

	return hex.EncodeToString(decryptedUsrPubKey), nil
}

func decryptHbcPriv(
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

	decryptedPrivKey, err := crypto.RsaDecryptOAEP(rsaPrivKey, privKeyBytes)
	if err != nil {
		return nil, code.NewI18nError(code.RSADecryptBackupDataErr, fmt.Sprintf("rsa decode privkey failed: %s", err.Error()))
	}
	decryptedChainCode, err := crypto.RsaDecryptOAEP(rsaPrivKey, chainCodeBytes)
	if err != nil {
		return nil, code.NewI18nError(code.RSADecryptBackupDataErr, fmt.Sprintf("rsa decode chaincode failed: %s", err.Error()))
	}

	decryptedPrivKey, err = ecies.Decrypt(eciesPrivKey, decryptedPrivKey)
	if err != nil {
		return nil, code.NewI18nError(code.EciesDecryptBackupDataErr, fmt.Sprintf("ecies decode privkey failed: %s", err.Error()))
	}
	decryptedChainCode, err = ecies.Decrypt(eciesPrivKey, decryptedChainCode)
	if err != nil {
		return nil, code.NewI18nError(code.EciesDecryptBackupDataErr, fmt.Sprintf("ecies decode chaincode failed: %s", err.Error()))
	}

	privateKey := new(big.Int).SetBytes(decryptedPrivKey)

	return &common.RootKey{
		PrivKey:     privateKey,
		EcdsaPubKey: crypto.ScalarBaseMult(crypto.S256(), privateKey),
		EddsaPubKey: crypto.ScalarBaseMult(crypto.Edwards(), privateKey),
		ChainCode:   decryptedChainCode,
	}, nil
}

func deriveChilds(vaultCount int, chains []string, rootKeys *common.RootKeys) ([]*DeriveResult, error) {
	var deriveResult []*DeriveResult

	for vaultIndex := 0; vaultIndex < vaultCount; vaultIndex++ {
		for _, chainName := range chains {
			coinInfo, ok := common.ChainInfos[chainName]
			if !ok {
				common.Logger.Errorf("unkonw chain: %s", chainName)
				return nil, fmt.Errorf("unkonw chain: %s", chainName)
			}

			hdPath := fmt.Sprintf(AssetWalletPath, vaultIndex, coinInfo.CoinType) // Only support asset wallet for now
			privKey, address, err := common.DeriveChild(rootKeys, hdPath, int(coinInfo.CoinType))
			if err != nil {
				return nil, fmt.Errorf("derive child failed, err: %s", err.Error())
			}

			var buf [32]byte
			privKeyBytes := privKey.FillBytes(buf[:])

			deriveResult = append(deriveResult, &DeriveResult{
				VaultIndex: vaultIndex + 1,
				Chain:      chainName,
				Address:    address,
				PrivKey:    hex.EncodeToString(privKeyBytes),
			})
		}
	}

	return deriveResult, nil
}

func deriveChilds2(vaultCount int, chains []string, rootKeys *common.RootKeys) ([]*DeriveResult, error) {
	var deriveResult []*DeriveResult

	for vaultIndex := 0; vaultIndex < vaultCount; vaultIndex++ {
		for _, chainName := range chains {
			coinInfo, _ := common.ChainInfos[chainName]

			hdPath := fmt.Sprintf(AssetWalletPath, vaultIndex, coinInfo.CoinType) // Only support asset wallet for now
			privKey, address, err := common.DeriveChild(rootKeys, hdPath, int(coinInfo.CoinType))
			if err != nil {
				return nil, fmt.Errorf("derive child failed, err: %s", err.Error())
			}

			var buf [32]byte
			privKeyBytes := privKey.FillBytes(buf[:])

			deriveResult = append(deriveResult, &DeriveResult{
				VaultIndex: vaultIndex + 1,
				Chain:      chainName,
				Address:    address,
				PrivKey:    hex.EncodeToString(privKeyBytes),
			})
		}
	}

	return deriveResult, nil
}

func concurrentDeriveChilds(vaultCount int, chains []string, rootKeys *common.RootKeys) ([]*DeriveResult, error) {
	deriveResult := make([]*DeriveResult, 0)

	var lock sync.Mutex

	var pError error

	chainTotal := len(chains)
	maxThread := 20
	currentThread := 0
	for i := 0; i <= int(math.Ceil(float64(chainTotal)/float64(maxThread))); i++ {
		wg := &sync.WaitGroup{}
		for j := 0; j < maxThread; j++ {
			if currentThread >= chainTotal {
				break
			}

			chainName := chains[currentThread]

			wg.Add(1)
			go func(chainName string) {
				defer wg.Done()

				vaultDeriveResult, err := deriveVaultChild(vaultCount, chainName, rootKeys)
				if err != nil {
					pError = err
					return
				}

				lock.Lock()
				deriveResult = append(deriveResult, vaultDeriveResult...)
				lock.Unlock()

			}(chainName)

			currentThread++
		}

		wg.Wait()
	}

	if pError != nil {
		return nil, pError
	}

	sortByVaultIndex := func(i, j int) bool {
		return deriveResult[i].VaultIndex < deriveResult[j].VaultIndex
	}

	sort.Slice(deriveResult, sortByVaultIndex)

	return deriveResult, nil
}

func deriveVaultChild(vaultCount int, chainName string, rootKeys *common.RootKeys) ([]*DeriveResult, error) {
	deriveResult := make([]*DeriveResult, 0)
	coinInfo, _ := common.ChainInfos[chainName]

	for vaultIndex := 0; vaultIndex < vaultCount; vaultIndex++ {
		hdPath := fmt.Sprintf(AssetWalletPath, vaultIndex, coinInfo.CoinType) // Only support asset wallet for now
		privKey, address, err := common.DeriveChild(rootKeys, hdPath, int(coinInfo.CoinType))
		if err != nil {
			return nil, err
		}

		var buf [32]byte
		privKeyBytes := privKey.FillBytes(buf[:])

		deriveResult = append(deriveResult, &DeriveResult{
			VaultIndex: vaultIndex + 1,
			Chain:      chainName,
			Address:    address,
			PrivKey:    hex.EncodeToString(privKeyBytes),
		})

	}
	return deriveResult, nil
}
