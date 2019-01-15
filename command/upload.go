package command

import (
	"encoding/base64"
	"errors"
	"github.com/ndphu/drive-manager-cli/config"
	"github.com/ndphu/google-api-helper"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
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

			wg:=sync.WaitGroup{}
			for i := 0; i < c.NArg(); i++ {
				localPath := c.Args().Get(i)
				log.Println(localPath)
				wg.Add(1)
				go func(fp string) {
					defer wg.Done()
					err := uploadFile(srv, fp, "" , "" ,"")
					if err != nil {
						log.Println(err.Error())
					}
				}(localPath)
			}
			wg.Wait()

			return nil
		},
	}
}

func uploadFile(srv *google_api_helper.DriveService, localPath string, name string, desc string, mimeType string) error {
	uploadFile, err := os.Open(localPath)
	if err != nil {
		return err
	}

	info, err := uploadFile.Stat()
	if err != nil {
		return errors.New("error: fail to stat file at "+localPath+" by error "+err.Error())
	}

	log.Println("uploading", info.Name(), "...")
	if strings.Trim(name, " ") == "" {
		name = info.Name()
	}

	if _, err := srv.UploadFile(name, desc, mimeType, localPath); err != nil {
		return errors.New(info.Name()+" error: fail to upload by error "+err.Error())
	}
	log.Println(info.Name(), "uploaded successfully")
	return nil
}
