package extrinsic

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/ed25519"

	"recovery-tool/tx/dot/ed25519WalletKey"
	"recovery-tool/tx/dot/polkadot-adapter/common"
)

const (
	AuthoritiesChangeType = 0
	TransferType          = 1
	IncludeDataType       = 2
	StorageChangeType     = 3
	// TODO: implement when storage changes trie is completed
	//ChangesTrieConfigUpdateType = 4
)

// Extrinsic represents a runtime Extrinsic
type Extrinsic interface {
	Type() int
	Encode() ([]byte, error)
	Decode(r io.Reader) error
}

// Transfer represents a runtime Transfer
type Transfer struct {
	from   [32]byte
	to     [32]byte
	amount uint64
	nonce  uint64
}

// NewTransfer returns a Transfer
func NewTransfer(from, to [32]byte, amount, nonce uint64) *Transfer {
	return &Transfer{
		from:   from,
		to:     to,
		amount: amount,
		nonce:  nonce,
	}
}

// Encode returns the SCALE encoding of the Transfer
func (t *Transfer) Encode() ([]byte, error) {
	enc := []byte{}

	buf := make([]byte, 8)

	enc = append(enc, t.from[:]...)
	enc = append(enc, t.to[:]...)

	binary.LittleEndian.PutUint64(buf, t.amount)
	enc = append(enc, buf...)

	binary.LittleEndian.PutUint64(buf, t.nonce)
	enc = append(enc, buf...)

	return enc, nil
}

// Decode decodes the SCALE encoding into a Transfer
func (t *Transfer) Decode(r io.Reader) (err error) {
	t.from, err = common.ReadHash(r)
	if err != nil {
		return err
	}

	t.to, err = common.ReadHash(r)
	if err != nil {
		return err
	}

	t.amount, err = ReadUint64(r)
	if err != nil {
		return err
	}

	t.nonce, err = ReadUint64(r)
	if err != nil {
		return err
	}

	return nil
}
func ReadUint64(r io.Reader) (uint64, error) {
	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

// AsSignedExtrinsic returns a TransferExt that includes the transfer and a signature.
func (t *Transfer) AsSignedExtrinsic(prikey []byte) (*TransferExt, error) {
	enc, err := t.Encode()
	if err != nil {
		return nil, err
	}
	publicKey := ed25519WalletKey.WalletPubKeyFromKeyBytes(prikey)
	privateByte := []byte{}
	fmt.Println(hex.EncodeToString(publicKey))
	privateByte = append(privateByte[:], prikey[:]...)
	privateByte = append(privateByte[:], publicKey[:]...)
	sig := ed25519.Sign(privateByte[:], enc)
	// sig, err := nil, nil
	// //sig, err := key.Sign(enc)
	// if err != nil {
	// 	return nil, err
	// }
	fmt.Println("sig = ", hex.EncodeToString(sig))
	sigb := [64]byte{}
	copy(sigb[:], sig)

	return NewTransferExt(t, sigb), nil
}

// TransferExt represents an Extrinsic::Transfer
type TransferExt struct {
	transfer  *Transfer
	signature [64]byte
}

// NewTransferExt returns a TransferExt
func NewTransferExt(transfer *Transfer, signature [64]byte) *TransferExt {
	return &TransferExt{
		transfer:  transfer,
		signature: signature,
	}
}

// Type returns TransferType
func (e *TransferExt) Type() int {
	return TransferType
}

// Encode returns the SCALE encoding of the TransferExt
func (e *TransferExt) Encode() ([]byte, error) {
	enc := []byte{TransferType}

	tenc, err := e.transfer.Encode()
	if err != nil {
		return nil, err
	}

	enc = append(enc, tenc...)
	enc = append(enc, e.signature[:]...)

	return enc, nil
}

// Decode decodes the SCALE encoding into a TransferExt
func (e *TransferExt) Decode(r io.Reader) error {
	e.transfer = new(Transfer)
	err := e.transfer.Decode(r)
	if err != nil {
		return err
	}

	_, err = r.Read(e.signature[:])
	if err != nil {
		return err
	}

	return nil
}
