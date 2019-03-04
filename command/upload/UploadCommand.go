package upload

import (
	"errors"
	"github.com/urfave/cli"
)

func GetUploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "upload file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "account",
				Usage: "Upload token from Web UI",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "Path to the file to upload",
			},
			&cli.BoolFlag{
				Name:  "recursive",
				Usage: "Should be true for uploading a directory",
			},
		},
		Action: func(c *cli.Context) error {
			account := c.String("account")
			if account == "" {
				return errors.New("should provide account to upload")
			}
			file := c.String("file")
			if file == "" {
				return errors.New("should provide file to upload")
			}

			return doUpload(account, file, c.Bool("recursive"))
		},
	}
}
