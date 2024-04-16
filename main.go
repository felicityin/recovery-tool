package main

import (
	"flag"
	"fmt"
	"recovery-tool/cmd"
	"recovery-tool/common"
)

func main() {
	var inputPath, outputPath string
	flag.StringVar(&inputPath, "i", "./input.yaml", "the path of input parmas")
	flag.StringVar(&outputPath, "o", "./output.yaml", "the path of result")

	err := cmd.RecoverKeys(inputPath, outputPath)
	if err != nil {
		common.Logger.Errorf("%s", err)
	}
	fmt.Printf("Output the result to file `%s`\n", outputPath)
}
