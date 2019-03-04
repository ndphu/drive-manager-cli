package login

import (
	"errors"
	"fmt"
	"github.com/ndphu/drive-manager-cli/config"
	"net/http"
)

func loginWithToken(conf *config.Config, token string, saveToken bool)  error {
	req, _ := http.NewRequest("GET", conf.BaseUrl + "/user/manage/info", nil)
	req.Header.Add("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("LOGIN_FAIL with status: %d", resp.StatusCode))
	}

	if saveToken {
		conf.Token = token
		return config.SaveConfig(conf)
	}
	return nil
}
