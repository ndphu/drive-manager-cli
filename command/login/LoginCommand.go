package login

import (
	"errors"
	"github.com/ndphu/drive-manager-cli/config"
	"github.com/urfave/cli"
)

func GetLoginCommand() cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "login to pre-registered user",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "token",
				Usage: "JWT token provided from Web UI",
			},
		},
		Action: func(c *cli.Context) error {
			jwtToken := c.String("token")
			if jwtToken == "" {
				return errors.New("should provide token to login")
			}
			conf, err := config.LoadConfig()
			if err != nil {
				return err
			}
			return loginWithToken(conf, jwtToken, true)
		},
	}
}
