package upload

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ndphu/drive-manager-cli/config"
	"github.com/ndphu/google-api-helper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func doUpload(account string, file string, recursive bool) error {
	stat, err := os.Stat(file)
	if err != nil {
		log.Println("input file not found", file, err)
		return err
	}
	if stat.IsDir() && !recursive {
		errMsg := fmt.Sprintf("%s %s", file, "is a directory. Add --recursive for uploading directory")
		return errors.New(errMsg)
	}

	conf, err := config.LoadConfig()
	if err != nil {
		log.Println("fail to load config", err)
		return err
	}
	log.Println("loading account key...")

	req, err := http.NewRequest("GET", conf.BaseUrl+"/manage/driveAccount/"+account+"/key", nil)
	if err != nil {
		log.Println("fail to create request", err)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+conf.Token)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("fail to get account key", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("server return status not success", resp.StatusCode)
		return errors.New(fmt.Sprintf("Server response status: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read response body from server", err)
		return err
	}

	key, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Println("fail to decode base64 key file from server")
		return err
	}

	srv, err := google_api_helper.GetDriveService(key)
	if err != nil {
		log.Println("fail to create drive service from key file")
		return err
	}

	log.Println("account setup finished.")

	if stat.IsDir() {
		doUploadDir(srv, file)
	} else {
		doUploadFile(srv, file)
	}
	return nil
}

func doUploadDir(srv *google_api_helper.DriveService, dir string) error {
	log.Println("processing directory", dir)
	children, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.IsDir() {
			if err := doUploadDir(srv, dir + "/" + child.Name()); err != nil {
				return err
			}
		} else {
			if err := doUploadFile(srv, dir + "/" + child.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

func doUploadFile(srv *google_api_helper.DriveService,  file string) error {
	stat, _ := os.Stat(file)
	log.Println("uploading", stat.Name(), "with size", stat.Size(), "...")

	uploaded, err := srv.UploadFile(stat.Name(), "", "", file)
	if err != nil {
		log.Println("fail to upload file", err)
		return err
	}
	log.Println("file uploaded successfully", uploaded.Name, "mime type", uploaded.MimeType)
	return nil
}