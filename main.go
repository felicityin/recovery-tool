package main

import (
	"flag"
	"fmt"
	"recovery-tool/cmd"
	"recovery-tool/common"
)

func main() {
	inputPath := flag.String("i", "./input.yaml", "the path of input parmas")
	outputPath := flag.String("o", "./output.yaml", "the path of result")
	flag.Parse()
	err := cmd.RecoverKeysCmd(*inputPath, *outputPath)
	if err != nil {
		common.Logger.Errorf("%s", err)
	}
	fmt.Printf("Output the result to file `%s`\n", *outputPath)
}
