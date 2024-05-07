package main

import (
	"os"

	"github.com/Gustrb/jbm/src/cli"
)

func main() {
	cli, err := cli.CreateCLI(os.Args)
	if err != nil {
		os.Exit(1)
		return
	}

	if err := cli.Run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
