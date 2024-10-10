package extrinsic

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_Extrinsic(t *testing.T) {
	alice, _ := hex.DecodeString("beff0e5d6f6e6e6d573d3044f3e2bfb353400375dc281da3337468d4aa527908")
	bob, _ := hex.DecodeString("cfa75b7b120968299f32ca2ff641dce4749f4fd0756f656275eeb286e9264690")
	aliceb := [32]byte{}
	copy(aliceb[:], alice)

	bobb := [32]byte{}
	copy(bobb[:], bob)
	transfer := NewTransfer(aliceb, bobb, 100000000000, 6)
	parse_hex := "612d82bc053d1b4729057688ecb1ebf62745d817ddd9b595bc822f5f2ba0e41a"
	sk, _ := hex.DecodeString(parse_hex)
	ext, err := transfer.AsSignedExtrinsic(sk)
	if err != nil {
		fmt.Println(err.Error())
	}
	enc, err := ext.Encode()
	if err != nil {
		fmt.Println(err.Error())
	}
	// r := &bytes.Buffer{}
	// r.Write(enc)
	fmt.Println("enc = ", hex.EncodeToString(enc))
}
