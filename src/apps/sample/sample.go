package main

import (
	"common"
	"fmt"
	_ "net/http/pprof"
	"os"
)

func app_main() {
	if err := common.ParseCommandAndFile(); err != nil {
		fmt.Fprintf(os.Stderr, "ParseCommandAndFile fail, err=[%v]\n", err)
		os.Exit(-1)
	}
}

func main() {
	common.BaseMain(app_main)
}
