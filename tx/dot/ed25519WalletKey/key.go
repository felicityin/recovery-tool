package ed25519WalletKey

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/agl/ed25519/edwards25519"
)

type WalletKey struct {
	PublicKey  []byte
	PrivateKey []byte
}

const (
	// FirstHardenedChild is the index of the firxt "harded" child key as per the
	// bip32 spec
	FirstHardenedChild = uint32(0x80000000)
	ScalarBytes        = 32
	// PublicKeyCompressedLength is the byte count of a compressed public key
	PublicKeyCompressedLength = 33
)

func NewWalletKey(Prv, Pub []byte) *WalletKey {

	return &WalletKey{
		PrivateKey: Prv,
		PublicKey:  Pub,
	}
}

type WalletExtKey struct {
	Key       *WalletKey
	ChainCode []byte // 32 bytes

}

func PruneScalar(s []byte) {
	s[0] &= 248
	s[31] &= 127 // clear top 3 bits
	s[31] |= 64  // set second highest bit
}

func WalletPubKeyFromKeyBytes(keyBytes []byte) []byte {
	digest := sha512.Sum512(keyBytes[:])
	//fmt.Println("digest len = ", len(digest))
	digest[0] &= 0xF8
	digest[31] &= 0x7F
	digest[31] |= 0x40

	var hBytes [32]byte

	copy(hBytes[:], digest[:])
	var A edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&A, &hBytes)
	var publicKeyBytes [32]byte
	A.ToBytes(&publicKeyBytes)

	return publicKeyBytes[:]
}

func GenerateWalletExtKey(seed []byte) (*WalletExtKey, error) {
	hmac := hmac.New(sha512.New, []byte("ed25519 seed"))
	_, err := hmac.Write(seed)

	if err != nil {
		return nil, err
	}
	intermediary := hmac.Sum(nil)
	//edd25519.NewKeyFromSeed(intermediary)
	// Split it into our key and chain code
	keyBytes := intermediary[:32]

	chainCode := intermediary[32:]
	var PrivateByte [32]byte
	copy(PrivateByte[:], keyBytes[:])
	var chainCodeByte [32]byte
	copy(chainCodeByte[:], chainCode[:])
	publicByte := WalletPubKeyFromKeyBytes(PrivateByte[:])
	key := NewWalletKey(PrivateByte[:], publicByte)
	return &WalletExtKey{
		Key:       key,
		ChainCode: chainCode,
	}, nil
}

func GenerateNEMWalletExtKey(seed []byte) (*WalletExtKey, error) {
	hmac := hmac.New(sha512.New, []byte("ed25519-keccak seed"))
	_, err := hmac.Write(seed)

	if err != nil {
		return nil, err
	}
	intermediary := hmac.Sum(nil)
	//edd25519.NewKeyFromSeed(intermediary)
	// Split it into our key and chain code
	keyBytes := intermediary[:32]

	chainCode := intermediary[32:]
	var PrivateByte [32]byte
	copy(PrivateByte[:], keyBytes[:])
	var chainCodeByte [32]byte
	copy(chainCodeByte[:], chainCode[:])
	publicByte := WalletPubKeyFromKeyBytes(PrivateByte[:])
	//fmt.Println(hex.EncodeToString(PrivateByte[:]))
	key := NewWalletKey(PrivateByte[:], publicByte)
	return &WalletExtKey{
		Key:       key,
		ChainCode: chainCode,
	}, nil
}

func NewWalletKeyFromMasterKey(seed []byte, derivedPath string) (*WalletExtKey, error) {
	path := strings.Replace(derivedPath, " ", "", -1)

	child, err := GenerateWalletExtKey(seed)
	if err != nil {
		return nil, err
	}

	path = path[2:]
	elements := strings.Split(path, "/")

	for _, elem := range elements {
		var hdSerializes uint32

		if strings.Index(elem, "'") == len(elem)-1 {
			elem = elem[0 : len(elem)-1]
			index, err := strconv.Atoi(elem)
			if err != nil {
				return nil, err
			}
			hdSerializes = uint32(index) + FirstHardenedChild
		} else {
			index, err := strconv.Atoi(elem)
			if err != nil {
				return nil, err
			}
			hdSerializes = uint32(index)
		}

		child, err = child.NewWalletChildKey(hdSerializes)
		if err != nil {
			return nil, err
		}

	}
	return child, nil
}

func NewNEMWalletKeyFromMasterKey(seed []byte, derivedPath string) (*WalletExtKey, error) {
	path := strings.Replace(derivedPath, " ", "", -1)

	child, err := GenerateNEMWalletExtKey(seed)
	if err != nil {
		return nil, err
	}

	path = path[2:]
	elements := strings.Split(path, "/")

	for _, elem := range elements {
		var hdSerializes uint32

		if strings.Index(elem, "'") == len(elem)-1 {
			elem = elem[0 : len(elem)-1]
			index, err := strconv.Atoi(elem)
			if err != nil {
				return nil, err
			}
			hdSerializes = uint32(index) + FirstHardenedChild
		} else {
			index, err := strconv.Atoi(elem)
			if err != nil {
				return nil, err
			}
			hdSerializes = uint32(index)
		}

		child, err = child.NewNEMWalletChildKey(hdSerializes)

		if err != nil {
			return nil, err
		}

	}
	pk := hex.EncodeToString(child.Key.PrivateKey)
	fmt.Println(pk)

	pri := hex2BaReversed(pk)
	fmt.Println(hex.EncodeToString(pri))
	publicByte := WalletPubKeyFromKeyBytes(pri[:])
	key := NewWalletKey(pri[:], publicByte)
	child.Key = key
	return child, nil
}
func hex2BaReversed(hexx string) []byte {
	data, err := hex.DecodeString(hexx)
	if err != nil {
		panic(err)
	}
	output := make([]byte, len(data))
	j := len(data) - 1
	for i := 0; i < len(data); i++ {
		output[j] = data[i]
		j--
	}
	return output
}

func (master *WalletExtKey) NewWalletChildKey(childNumber uint32) (*WalletExtKey, error) {
	// fmt.Println(N.Bytes())
	data := []byte{}
	publicKey := master.Key.PublicKey
	privateKey := master.Key.PrivateKey

	if childNumber/FirstHardenedChild == 1 {
		temp := uint8(0)
		// data = data.append(&temp)
		data = append(data[:], byte(temp))
		data = append(data[:], privateKey[:32]...)

	} else {
		temp := uint8(0)
		// data = data.append(&temp)
		data = append(data[:], byte(temp))
		data = append(data[:], publicKey[:32]...)
	}
	var path [4]byte

	// |---这部分go在强转的时候扔掉---|
	path[3] = uint8(childNumber & 0xFF)
	// 下面是右移后的数据
	path[2] = uint8((childNumber >> 8) & 0xFF)
	// 下面是右移后的数据
	path[1] = uint8((childNumber >> 16) & 0xFF)
	// 下面是右移后的数据
	path[0] = uint8((childNumber >> 24) & 0xFF)

	//fmt.Printf("%+v\n", path)
	data = append(data[:], path[:]...)
	var iData [37]byte
	copy(iData[:], data[:37])
	intermediary := HMAC_SHA512(iData[:], master.ChainCode)
	// fmt.Println("HMAC_SHA512 : ", i)
	keyBytes := intermediary[:32]
	chainCode := intermediary[32:]
	var PrivateByte [32]byte
	copy(PrivateByte[:], keyBytes[:])

	//fmt.Println("EncodePublicKey len : ", len(ec.EncodePublicKey(&ecPrivateKey.PublicKey, true)))
	// Create the key struct
	publicByte := WalletPubKeyFromKeyBytes(keyBytes[:])
	key := NewWalletKey(keyBytes[:], publicByte)
	return &WalletExtKey{
		Key:       key,
		ChainCode: chainCode,
	}, nil

}

func (master *WalletExtKey) NewNEMWalletChildKey(childNumber uint32) (*WalletExtKey, error) {
	// fmt.Println(N.Bytes())
	data := []byte{}
	publicKey := master.Key.PublicKey
	privateKey := master.Key.PrivateKey

	if childNumber/FirstHardenedChild == 1 {
		temp := uint8(0)
		// data = data.append(&temp)
		data = append(data[:], byte(temp))
		data = append(data[:], privateKey[:32]...)

	} else {
		temp := uint8(0)
		// data = data.append(&temp)
		data = append(data[:], byte(temp))
		data = append(data[:], publicKey[:32]...)
	}
	path := encode32BE(childNumber)

	//fmt.Printf("%+v\n", path)
	data = append(data[:], path[:]...)
	var iData [37]byte
	copy(iData[:], data[:37])
	intermediary := HMAC_SHA512(iData[:], master.ChainCode)
	// fmt.Println("HMAC_SHA512 : ", i)
	keyBytes := intermediary[:32]
	chainCode := intermediary[32:]
	var PrivateByte [32]byte
	copy(PrivateByte[:], keyBytes[:])

	//fmt.Println("EncodePublicKey len : ", len(ec.EncodePublicKey(&ecPrivateKey.PublicKey, true)))
	// Create the key struct
	publicByte := WalletPubKeyFromKeyBytes(keyBytes[:])
	key := NewWalletKey(keyBytes[:], publicByte)

	return &WalletExtKey{
		Key:       key,
		ChainCode: chainCode,
	}, nil

}
func encode32BE(i uint32) []byte {
	var b [4]byte
	b[0] = byte(i >> 24)
	b[1] = byte(i >> 16)
	b[2] = byte(i >> 8)
	b[3] = byte(i)
	return b[:]
}
func HMAC_SHA512(data []byte, key []byte) []byte {
	m := hmac.New(sha512.New, []byte(key))
	m.Write([]byte(data))

	return m.Sum(nil)
}
