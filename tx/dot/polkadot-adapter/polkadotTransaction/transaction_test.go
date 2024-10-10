package polkadotTransaction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"recovery-tool/tx/eddsa"
	"testing"
)

func Test_KSM_transaction(t *testing.T) {
	tx := TxStruct{
		//发送方公钥
		SenderPubkey: "86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//接收方公钥
		RecipientPubkey: "88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//发送金额（最小单位）
		Amount: 12,
		//nonce
		Nonce: 1,
		//手续费（最小单位）
		Fee: 20,
		//当前高度
		BlockHeight: 1778228,
		//当前高度区块哈希
		BlockHash: "bae19137f56d7c7bc88350131dd401c80c77ad3ffca7157bbf2d008a4d0dd8f4",
		//创世块哈希
		GenesisHash: "b0a8d493285c2df73290dfb7e61f870f17b41801197a149ca93654499ea3dafe",
		//spec版本
		SpecVersion: 1059,
		//Transaction版本
		TxVersion: 1,
	}

	// 创建空交易单和待签消息
	emptyTrans, message, err := tx.CreateEmptyTransactionAndMessage(KSM_Balannce_Transfer)
	if err != nil {
		t.Error("create failed : ", err)
		return
	}
	fmt.Println("空交易单 ： ", emptyTrans)
	fmt.Println("待签消息 ： ", message)

	// 签名
	prikey, _ := hex.DecodeString("e86bcaaab0a5aa5e3f3b0885db7e932e34eddb5a620b6bcc097a4b236a5a0354")
	signature, err := SignTransaction(message, prikey)
	if err != nil {
		t.Error("sign failed")
		return
	}
	fmt.Println("签名结果 ： ", hex.EncodeToString(signature))

	// signature, _ := hex.DecodeString("1cc69f7ba50ee37793c83d74b21f50239894e8733cdf7fd13565eded13ba97d8229fc51174035be6d4543908f58b016efd0aae137f8ad584c5540002326bc809")

	// 验签与交易单合并
	signedTrans, pass := VerifyAndCombineTransaction(KSM_Balannce_Transfer, emptyTrans, hex.EncodeToString(signature))
	if pass {
		fmt.Println("验签成功")
		fmt.Println("签名交易单 ： ", signedTrans)
	} else {
		t.Error("验签失败")
	}
}

func Test_DOT_transaction(t *testing.T) {
	tx := TxStruct{
		//发送方公钥
		SenderPubkey: "86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//接收方公钥
		RecipientPubkey: "88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//发送金额（最小单位）
		Amount: 12,
		//nonce
		Nonce: 1,
		//手续费（最小单位）
		Fee: 20,
		//当前高度
		BlockHeight: 1778228,
		//当前高度区块哈希
		BlockHash: "bae19137f56d7c7bc88350131dd401c80c77ad3ffca7157bbf2d008a4d0dd8f4",
		//创世块哈希
		GenesisHash: "b0a8d493285c2df73290dfb7e61f870f17b41801197a149ca93654499ea3dafe",
		//spec版本
		SpecVersion: 1059,
		//Transaction版本
		TxVersion: 1,
	}

	// 创建空交易单和待签消息
	emptyTrans, message, err := tx.CreateEmptyTransactionAndMessage(DOT_Balannce_Transfer)
	if err != nil {
		t.Error("create failed : ", err)
		return
	}
	fmt.Println("空交易单 ： ", emptyTrans)
	fmt.Println("待签消息 ： ", message)

	// 签名
	prikey, _ := hex.DecodeString("e86bcaaab0a5aa5e3f3b0885db7e932e34eddb5a620b6bcc097a4b236a5a0354")
	signature, err := SignTransaction(message, prikey)
	if err != nil {
		t.Error("sign failed")
		return
	}
	fmt.Println("签名结果 ： ", hex.EncodeToString(signature))

	// signature, _ := hex.DecodeString("1cc69f7ba50ee37793c83d74b21f50239894e8733cdf7fd13565eded13ba97d8229fc51174035be6d4543908f58b016efd0aae137f8ad584c5540002326bc809")

	// 验签与交易单合并
	signedTrans, pass := VerifyAndCombineTransaction(DOT_Balannce_Transfer, emptyTrans, hex.EncodeToString(signature))
	if pass {
		fmt.Println("验签成功")
		fmt.Println("签名交易单 ： ", signedTrans)
	} else {
		t.Error("验签失败")
	}
}

func Test_DOT_tx(t *testing.T) {
	tx := TxStruct{
		//发送方公钥
		SenderPubkey: "abd468b970ad280a6202063cf18f3b5893d7c59760fded8a96377765bbf83ba4", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//接收方公钥
		RecipientPubkey: "e26f921fb2cb134fdbb94bd4fa96178b25f042bc2219c8b1ea7952dabb86f2a6",
		//发送金额（最小单位）
		Amount: 30000000,
		//nonce
		Nonce: 0,
		//手续费（最小单位）
		Fee: 260000000,
		//当前高度
		BlockHeight: 0,
		//当前高度区块哈希
		BlockHash: "55cf04a4b7722492fa85037fc5a216d1a95333576a708bd74f8f6af72baa51c0",
		//创世块哈希
		GenesisHash: "91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3",
		//spec版本
		SpecVersion: 1002007,
		//Transaction版本
		TxVersion: 26,
	}

	// 创建空交易单和待签消息
	emptyTrans, message, err := tx.CreateEmptyTransactionAndMessage(DOT_Balannce_Transfer)
	if err != nil {
		t.Error("create failed : ", err)
		return
	}
	fmt.Println("空交易单 ： ", emptyTrans)
	fmt.Println("待签消息 ： ", message)

	// 签名
	privkey, _ := hex.DecodeString("093764441b70738e147516ded739e6d2c7c8e39a267ba578721b0cb872970336")

	msgBytes, err := hex.DecodeString(message)
	if err != nil {
		err = fmt.Errorf("GetSignPacket err: %s", err.Error())
		return
	}

	fmt.Println("privkey: ", privkey)
	fmt.Println("msg: ", msgBytes)

	signature, err := eddsa.Sign(privkey, msgBytes)
	if err != nil {
		return
	}
	fmt.Println("签名结果 ： ", hex.EncodeToString(signature))

	// 验签与交易单合并
	signedTrans, pass := VerifyAndCombineTransaction(DOT_Balannce_Transfer, emptyTrans, hex.EncodeToString(signature))
	fmt.Printf("tx: %s\n", emptyTrans)
	if pass {
		fmt.Println("验签成功")
		fmt.Println("签名交易单 ： ", signedTrans)
	} else {
		t.Error("验签失败")
	}
}

func Test_json(t *testing.T) {
	ts := TxStruct{
		SenderPubkey:    "123",
		RecipientPubkey: "",
		Amount:          0,
		Nonce:           0,
		Fee:             0,
		BlockHeight:     0,
		BlockHash:       "234",
		GenesisHash:     "345",
		SpecVersion:     0,
	}

	js, _ := json.Marshal(ts)

	fmt.Println(string(js))
}
