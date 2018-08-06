package main

import (
	"os"

	"github.com/olegsu/iris/pkg/cli"
)

func main() {
	app := cli.CreateApplication()
	app.Run(os.Args)
}
