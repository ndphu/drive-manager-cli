package command

import (
	"encoding/base64"
	"github.com/ndphu/drive-manager-cli/config"
	"github.com/ndphu/google-api-helper"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func UploadCommand() *cli.Command {
	return &cli.Command{
		Name:  "upload",
		Usage: "upload a file to drive account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "account",
				Usage: "account to upload",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "name of the uploaded file",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "mime-type",
				Usage: "override mime-type of the uploading file",
				Value: "",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("account") == "" {
				return cli.Exit("error: should provide account to upload", -1)
			}

			if c.NArg() == 0 {
				return cli.Exit("error: should provide file to upload", -1)
			}

			localPath := c.Args().Get(0)
			uploadFile, err := os.Open(localPath)
			if err != nil {
				return cli.Exit("error: cannot opening input file: "+err.Error(), -1)
			}

			info, err := uploadFile.Stat()
			if err != nil {
				return cli.Exit("error: fail to stat file: "+err.Error(), -1)
			}
			log.Println("prepare uploading", info.Size(), "bytes using account", c.String("account"))

			log.Println("getting account key...")

			resp, err := http.Get(config.GetConfig().BackendUrl + "/manage/driveAccount/" + c.String("account") + "/key")
			if err != nil {
				return cli.Exit("error: fail to get account key "+err.Error(), -1)
			}
			defer resp.Body.Close()

			keyEncoded, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return cli.Exit("error: fail to get account key "+err.Error(), -1)
			}

			keyDecoded, err := base64.StdEncoding.DecodeString(string(keyEncoded))
			if err != nil {
				return cli.Exit("error: fail to decode account key "+err.Error(), -1)
			}

			srv, err := google_api_helper.GetDriveService([]byte(keyDecoded))
			if err != nil {
				return cli.Exit("error: fail to initialize drive service from key "+err.Error(), -1)
			}

			quota, err := srv.GetQuotaUsage()
			if err != nil {
				return cli.Exit("error: fail to query quota "+err.Error(), -1)
			}

			if info.Size() > quota.Limit-quota.Usage {
				return cli.Exit("error: not enough free space on target account.", -1)
			}
			name := c.String("name")
			if strings.Trim(name, " ") == "" {
				name = info.Name()
			}

			uploaded, err := srv.UploadFile(name, c.String("desc"), c.String("mime-type"), localPath)
			if err != nil {
				return cli.Exit("error: fail to upload by error "+err.Error(), -1)
			}
			log.Println("file uploaded successfully", uploaded.Id)

			return nil
		},
	}
}
