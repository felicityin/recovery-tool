package main

import (
	"flag"
	"fmt"
	"recovery-tool/cmd"
	"recovery-tool/common"
	"time"
)

func main() {
	var inputPath, outputPath string
	flag.StringVar(&inputPath, "i", "./input.yaml", "the path of input parmas")
	flag.StringVar(&outputPath, "o", "./output.yaml", "the path of result")
	start := time.Now()
	inputPath = "./input1.yaml"
	outputPath = "./output1.yaml"
	err := cmd.RecoverKeysCmd(inputPath, outputPath)
	if err != nil {
		common.Logger.Errorf("%s", err)
	}
	duration := time.Since(start)
	fmt.Printf("Output the result to file `%s`, cost: %s \n ", outputPath, duration.String())
}
