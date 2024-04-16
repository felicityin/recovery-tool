package main

import (
	"flag"
	"fmt"
	"os"
	"recovery-tool/cmd"
	"recovery-tool/common"
)

const (
	taskRecover = "recover"
)

func main() {
	recoverCmd := flag.NewFlagSet(taskRecover, flag.ExitOnError)
	recoveryParamsPath := recoverCmd.String("i", "./recovery.yaml", "the path of input parmas")
	recoveryOutputPath := recoverCmd.String("o", "./recovery_output.yaml", "the path of result")

	if len(os.Args) < 2 {
		fmt.Printf("expect '%s' subcommand\n", taskRecover)
		os.Exit(1)
	}

	switch os.Args[1] {
	case taskRecover:
		err := cmd.RecoverKeys(*recoveryParamsPath, *recoveryOutputPath)
		if err != nil {
			common.Logger.Errorf("%s", err)
		}
		fmt.Printf("Output the result to file `%s`\n", *recoveryOutputPath)
	default:
		fmt.Printf("expect '%s' subcommand\n", taskRecover)
		os.Exit(1)
	}
}
