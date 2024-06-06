package main

// must remove import "C"
//func TestGoRecoveryTest(t *testing.T) {
//	bytess, err := os.ReadFile("../input1.yaml")
//	if err != nil {
//		panic(err)
//	}
//
//	var input cmd.RecoveryInput
//	if err := yaml.UnmarshalStrict(bytess, &input); err != nil {
//		panic(err)
//	}
//
//	vaultCountStr := strconv.Itoa(input.VaultCount)
//	chainStr := ""
//	for _, chainName := range input.Chains {
//		chainStr += chainName + ","
//	}
//	chainStr = strings.TrimRight(chainStr, ",")
//	GoRecoveryTest(input.ZipPath, input.UserMnemonic, input.EciesPrivKey, "./test/RSAKet", vaultCountStr, chainStr, "zh")
//}
