package address

import (
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
)

var ss58Prefix = []byte("SS58PRE")
var ssPrefix = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
var encodeType = "base58"
var NetWorkByteMap = map[string]byte{
	"DOT": 0x00,
	"KSM": 0x00,
}

func PublicKeyToAddress(pub []byte, net string, addrPrefix uint8) string {
	//addrPrefix := 0
	if net == "testnet" {
		addrPrefix = 42
	}
	enc := append([]byte{byte(addrPrefix)}, pub...) // 42 测试网
	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return ""
	}
	_, err = hasher.Write(append(ss58Prefix, enc...))
	if err != nil {
		return ""
	}
	checksum := hasher.Sum(nil)

	return base58.Encode(append(enc, checksum[:2]...))
}

func PublicKeyToKsmAddress(pub []byte, net string) string {
	addrPrefix := 2
	if net == "testnet" {
		addrPrefix = 42
	}
	enc := append([]byte{byte(addrPrefix)}, pub...) // 42 测试网
	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return ""
	}
	_, err = hasher.Write(append(ss58Prefix, enc...))
	if err != nil {
		return ""
	}
	checksum := hasher.Sum(nil)

	return base58.Encode(append(enc, checksum[:2]...))
}

func PublicKeyToPlmAddress(pub []byte, net string) string {
	addrPrefix := 5
	if net == "testnet" {
		addrPrefix = 42
	}
	enc := append([]byte{byte(addrPrefix)}, pub...) // 42 测试网
	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return ""
	}
	_, err = hasher.Write(append(ss58Prefix, enc...))
	if err != nil {
		return ""
	}
	checksum := hasher.Sum(nil)

	return base58.Encode(append(enc, checksum[:2]...))
}

func PublicKeyToFisAddress(pub []byte, net string) string {
	addrPrefix := 20
	if net == "testnet" {
		addrPrefix = 42
	}
	enc := append([]byte{byte(addrPrefix)}, pub...) // 42 测试网
	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return ""
	}
	_, err = hasher.Write(append(ss58Prefix, enc...))
	if err != nil {
		return ""
	}
	checksum := hasher.Sum(nil)

	return base58.Encode(append(enc, checksum[:2]...))
}

func AddressEncode(hash []byte, addrPrefix uint8, opts ...interface{}) (string, error) {
	if len(hash) != 32 {
		hash, _ = owcrypt.CURVE25519_convert_Ed_to_X(hash)
	}
	prefix := []byte{addrPrefix}
	data := addressEncoder.CatData(prefix, hash)
	input := addressEncoder.CatData(ssPrefix, data)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]
	result := addressEncoder.EncodeData(addressEncoder.CatData(data, checkSum), encodeType, addressEncoder.BTCAlphabet)
	return result, nil
}

func AddressKsmEncode(hash []byte, opts ...interface{}) (string, error) {
	if len(hash) != 32 {
		hash, _ = owcrypt.CURVE25519_convert_Ed_to_X(hash)
	}
	prefix := []byte{2}
	data := addressEncoder.CatData(prefix, hash)
	input := addressEncoder.CatData(ssPrefix, data)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]
	result := addressEncoder.EncodeData(addressEncoder.CatData(data, checkSum), encodeType, addressEncoder.BTCAlphabet)
	return result, nil
}

func AddressPlmEncode(hash []byte, opts ...interface{}) (string, error) {
	if len(hash) != 32 {
		hash, _ = owcrypt.CURVE25519_convert_Ed_to_X(hash)
	}
	prefix := []byte{5}
	data := addressEncoder.CatData(prefix, hash)
	input := addressEncoder.CatData(ssPrefix, data)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]
	result := addressEncoder.EncodeData(addressEncoder.CatData(data, checkSum), encodeType, addressEncoder.BTCAlphabet)
	return result, nil
}

func AddressFisEncode(hash []byte, opts ...interface{}) (string, error) {
	if len(hash) != 32 {
		hash, _ = owcrypt.CURVE25519_convert_Ed_to_X(hash)
	}
	prefix := []byte{20}
	data := addressEncoder.CatData(prefix, hash)
	input := addressEncoder.CatData(ssPrefix, data)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]
	result := addressEncoder.EncodeData(addressEncoder.CatData(data, checkSum), encodeType, addressEncoder.BTCAlphabet)
	return result, nil
}

func AddressDecode(addr string, opts ...interface{}) ([]byte, error) {
	data, err := addressEncoder.Base58Decode(addr, addressEncoder.NewBase58Alphabet(addressEncoder.BTCAlphabet))
	if err != nil {
		return nil, err
	}
	pubkey := data[1 : len(data)-2]
	return pubkey, nil
}

func AddressVerify(address string, opts ...interface{}) bool {
	P2PKHPrefix := byte(0)
	P2PKHTestNetPrefix := byte(42)
	decodeBytes, err := addressEncoder.Base58Decode(address, addressEncoder.NewBase58Alphabet(addressEncoder.BTCAlphabet))
	if err != nil || len(decodeBytes) != 35 {
		return false
	}
	if decodeBytes[0] != P2PKHPrefix && decodeBytes[0] != P2PKHTestNetPrefix {
		return false
	}
	pub := decodeBytes[1 : len(decodeBytes)-2]
	prefix := []byte{0}
	if decodeBytes[0] == P2PKHTestNetPrefix {
		prefix = []byte{42}
	}
	data := append(prefix, pub...)
	input := append(ssPrefix, data...)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]

	for i := 0; i < 2; i++ {
		if checkSum[i] != decodeBytes[33+i] {
			return false
		}
	}
	if len(pub) != 32 {
		return false
	}
	return true
}

func AddressKsmVerify(address string, opts ...interface{}) bool {
	P2PKHPrefix := byte(2)
	P2PKHTestNetPrefix := byte(42)
	decodeBytes, err := addressEncoder.Base58Decode(address, addressEncoder.NewBase58Alphabet(addressEncoder.BTCAlphabet))
	if err != nil || len(decodeBytes) != 35 {
		return false
	}
	if decodeBytes[0] != P2PKHPrefix && decodeBytes[0] != P2PKHTestNetPrefix {
		return false
	}
	pub := decodeBytes[1 : len(decodeBytes)-2]
	prefix := []byte{2}
	if decodeBytes[0] == P2PKHTestNetPrefix {
		prefix = []byte{42}
	}
	data := append(prefix, pub...)
	input := append(ssPrefix, data...)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]

	for i := 0; i < 2; i++ {
		if checkSum[i] != decodeBytes[33+i] {
			return false
		}
	}
	if len(pub) != 32 {
		return false
	}
	return true
}

func AddressPlmVerify(address string, opts ...interface{}) bool {
	P2PKHPrefix := byte(5)
	P2PKHTestNetPrefix := byte(42)
	decodeBytes, err := addressEncoder.Base58Decode(address, addressEncoder.NewBase58Alphabet(addressEncoder.BTCAlphabet))
	if err != nil || len(decodeBytes) != 35 {
		return false
	}
	if decodeBytes[0] != P2PKHPrefix && decodeBytes[0] != P2PKHTestNetPrefix {
		return false
	}
	pub := decodeBytes[1 : len(decodeBytes)-2]
	prefix := []byte{5}
	if decodeBytes[0] == P2PKHTestNetPrefix {
		prefix = []byte{42}
	}
	data := append(prefix, pub...)
	input := append(ssPrefix, data...)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]

	for i := 0; i < 2; i++ {
		if checkSum[i] != decodeBytes[33+i] {
			return false
		}
	}
	if len(pub) != 32 {
		return false
	}
	return true
}
