package ton

type Account struct {
	Nonce uint64 `json:"nonce"`
	//账户是否激活
	IsActive bool `json:"isActive"`
	//当前账户状态
	Status  string `json:"status"`
	Balance string `json:"balance"`
}

type JettonWalletsRes struct {
	Error         string `json:"error"`
	JettonWallets []struct {
		Address           string `json:"address"`
		Balance           string `json:"balance"`
		Owner             string `json:"owner"`
		Jetton            string `json:"jetton"`
		LastTransactionLt string `json:"last_transaction_lt"`
		CodeHash          string `json:"code_hash"`
		DataHash          string `json:"data_hash"`
	} `json:"jetton_wallets"`
}

type MasterBlockChainInfo struct {
	Last struct {
		Workchain              int64  `json:"workchain"`
		Shard                  string `json:"shard"`
		Seqno                  int64  `json:"seqno"`
		RootHash               string `json:"root_hash"`
		FileHash               string `json:"file_hash"`
		GlobalID               int64  `json:"global_id"`
		Version                int64  `json:"version"`
		AfterMerge             bool   `json:"after_merge"`
		BeforeSplit            bool   `json:"before_split"`
		AfterSplit             bool   `json:"after_split"`
		WantMerge              bool   `json:"want_merge"`
		WantSplit              bool   `json:"want_split"`
		KeyBlock               bool   `json:"key_block"`
		VertSeqnoIncr          bool   `json:"vert_seqno_incr"`
		Flags                  int64  `json:"flags"`
		GenUtime               string `json:"gen_utime"`
		StartLt                string `json:"start_lt"`
		EndLt                  string `json:"end_lt"`
		ValidatorListHashShort int64  `json:"validator_list_hash_short"`
		GenCatchainSeqno       int64  `json:"gen_catchain_seqno"`
		MinRefMcSeqno          int64  `json:"min_ref_mc_seqno"`
		PrevKeyBlockSeqno      int64  `json:"prev_key_block_seqno"`
		VertSeqno              int64  `json:"vert_seqno"`
		MasterRefSeqno         int64  `json:"master_ref_seqno"`
		RandSeed               string `json:"rand_seed"`
		CreatedBy              string `json:"created_by"`
		TxCount                int64  `json:"tx_count"`
		MasterchainBlockRef    struct {
			Workchain int64  `json:"workchain"`
			Shard     string `json:"shard"`
			Seqno     int64  `json:"seqno"`
		} `json:"masterchain_block_ref"`
		PrevBlocks []struct {
			Workchain int64  `json:"workchain"`
			Shard     string `json:"shard"`
			Seqno     int64  `json:"seqno"`
		} `json:"prev_blocks"`
	} `json:"last"`
}

type WalletInfo struct {
	Balance             string `json:"balance"`
	WalletType          string `json:"wallet_type"`
	Seqno               uint64 `json:"seqno"`
	WalletID            uint64 `json:"wallet_id"`
	LastTransactionLt   string `json:"last_transaction_lt"`
	LastTransactionHash string `json:"last_transaction_hash"`
	Status              string `json:"status"`
}

type SendTxParams struct {
	Boc string `json:"boc"`
}

type SendTxResult struct {
	Error       string `json:"error"`
	MessageHash string `json:"message_hash"`
	Detail      []struct {
		Loc  []interface{} `json:"loc"`
		Msg  string        `json:"msg"`
		Type string        `json:"type"`
	} `json:"detail"`
}
