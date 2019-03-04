package main

import (
	"github.com/ndphu/drive-manager-cli/command/login"
	"github.com/ndphu/drive-manager-cli/command/upload"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		login.GetLoginCommand(),
		upload.GetUploadCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
