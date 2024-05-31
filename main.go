package main

import (
	"flag"
	"fmt"
	"recovery-tool/cmd"
	"recovery-tool/common"
	"time"
)

func main() {
	inputPath := flag.String("i", "./input.yaml", "the path of input parmas")
	outputPath := flag.String("o", "./output.yaml", "the path of result")
	flag.Parse()
	start := time.Now()
	err := cmd.RecoverKeysCmd(*inputPath, *outputPath)
	if err != nil {
		common.Logger.Errorf("%s", err)
	}
	duration := time.Since(start)
	fmt.Printf("Output the result to file `%s`, cost: %s \n ", *outputPath, duration.String())
}
