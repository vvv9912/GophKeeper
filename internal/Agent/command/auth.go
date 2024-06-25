package command

import (
	"GophKeeper/pkg/logger"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"regexp"
	"strings"
)

// SignIn - вход в аккаунт
func (c *Cobra) SignIn(cmd *cobra.Command, args []string) {
	Login, Password := auth(args)
	log.Println("Login:", Login)
	log.Println("Password:", Password)

	Login = strings.ToLower(Login)

	jwt, err := c.s.SignIn(cmd.Context(), Login, Password)
	if err != nil {
		logger.Log.Error("Sign in failed", zap.Error(err))
		return
	}
	logger.Log.Info("Sign in success")

	logger.Log.Debug("jwt", zap.String("jwt", jwt))
}

// SignUp - регистрация
func (c *Cobra) SignUp(cmd *cobra.Command, args []string) {
	Login, Password := auth(args)

	log.Println("Login:", Login)
	log.Println("Password:", Password)

	Login = strings.ToLower(Login)

	jwt, err := c.s.SignUp(cmd.Context(), Login, Password)
	if err != nil {
		logger.Log.Error("Sign Up  failed", zap.Error(err))
		return
	}
	logger.Log.Info("The user has been successfully registered")
	logger.Log.Debug("jwt", zap.String("jwt", jwt))
}

func auth(args []string) (string, string) {

	// Ограничение по длине символов для логин, пароля
	const (
		LenLoginMax = 20
		LenLoginMin = 5

		LenPasswordMax = 20
		LenPasswordMin = 5
	)

	var (
		Login    string
		Password string
	)

	if len(args) == 2 {
		Login = args[0]
		Password = args[1]
	} else {

		for {
			fmt.Print("Enter Login: ")
			fmt.Scanln(&Login)

			if len(Login) > LenLoginMax || len(Login) < LenLoginMin {
				logger.Log.Error("Error Login")
				continue
			}

			match, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]*$", Login)
			if err != nil {
				logger.Log.Error("Error Login")
				continue
			}
			if !match {
				logger.Log.Error("Error Login")
				continue
			}

			logger.Log.Debug("Password: ", zap.String("Password", Login))

			passwordByte, err := gopass.GetPasswdMasked()
			if err != nil {
				logger.Log.Error("Error Password")
				continue
			}

			if len(Login) > LenPasswordMax || len(Login) < LenPasswordMin {
				logger.Log.Error("Error Login")
				continue
			}
			Password = string(passwordByte)
			match, err = regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]*$", Password)
			if err != nil {
				logger.Log.Error("Error Password")
				continue
			}
			if !match {
				logger.Log.Error("Error Password")
				continue
			}
			break
		}

	}
	return Login, Password
}
