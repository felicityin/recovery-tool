package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"recovery-tool/cmd"
	"recovery-tool/common"
)

func main() {
	recoverCmd := flag.NewFlagSet("recover", flag.ExitOnError)
	inputPath := recoverCmd.String("i", "./input.yaml", "The path of input parmas")
	outputPath := recoverCmd.String("o", "./output.yaml", "The path of result")

	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	address := balanceCmd.String("addr", "", "address")
	coin := balanceCmd.String("coin", "", "Coin contract address. For sol, refer to https://solscan.io/leaderboard/token")
	chain := balanceCmd.String("chain", "sol", "Chain name, can be sol, apt or dot")
	url := balanceCmd.String("url", "https://api.mainnet-beta.solana.com", "url")

	transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	fromkey := transferCmd.String("fromkey", "", "Private key")
	toAddress := transferCmd.String("to", "", "Address")
	amount := transferCmd.String("amount", "", "Amount")
	coinAddress := transferCmd.String("coin", "", "Coin contract address. For sol, refer to https://solscan.io/leaderboard/token")
	chainName := transferCmd.String("chain", "sol", "Chain name, can be sol, apt or dot")
	chainUrl := transferCmd.String("url", "https://api.mainnet-beta.solana.com", "url")
	memo := transferCmd.String("memo", "", "Memo")

	if len(os.Args) < 2 {
		fmt.Println("expected 'recover', 'balance' or 'transfer' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "recover":
		recoverCmd.Parse(os.Args[2:])

		start := time.Now()
		err := cmd.RecoverKeysCmd(*inputPath, *outputPath)
		if err != nil {
			common.Logger.Errorf("%s", err)
			os.Exit(1)
		}
		duration := time.Since(start)
		fmt.Printf("Output the result to file `%s`, cost: %s \n ", *outputPath, duration.String())
	case "balance":
		balanceCmd.Parse(os.Args[2:])

		amount, err := cmd.GetBalance(*chain, *url, *address, *coin)
		if err != nil {
			common.Logger.Errorf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("balance: %s\n", amount.Balance)
		fmt.Printf("decimals: %s\n", amount.Decimals)
		fmt.Printf("amount: %s\n", amount.Amount)
	case "transfer":
		transferCmd.Parse(os.Args[2:])

		txHash, err := cmd.Transfer(*chainName, *chainUrl, *fromkey, *toAddress, *amount, *coinAddress, *memo)
		if err != nil {
			common.Logger.Errorf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("tx: %s/%s\n", cmd.Scan(*chainName), txHash)
	default:
		fmt.Println("expected 'recover', 'balance' or 'transfer' subcommands")
		os.Exit(1)
	}
}
