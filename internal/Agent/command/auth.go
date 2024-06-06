package command

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func (c *Cobra) SignIn(cmd *cobra.Command, args []string) {
	Login, Password := auth(args)

	fmt.Println("Login:", Login)
	fmt.Println("Password:", Password)

	Login = strings.ToLower(Login)

	jwt, err := c.s.SignIn(cmd.Context(), Login, Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("The user has been successfully logged in")
	fmt.Println(jwt)
}

func (c *Cobra) SignUp(cmd *cobra.Command, args []string) {
	Login, Password := auth(args)

	fmt.Println("Login:", Login)
	fmt.Println("Password:", Password)

	Login = strings.ToLower(Login)

	jwt, err := c.s.SignUp(cmd.Context(), Login, Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("The user has been successfully registered")
	fmt.Println(jwt)
}

func auth(args []string) (string, string) {
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

			if len(Login) > 20 || len(Login) < 5 {
				fmt.Println("Error Login") //todo
				continue
			}
			match, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]*$", Login)
			if err != nil {
				fmt.Println("Error Login")
				continue
			}
			if !match {
				fmt.Println("Enter Login") //todo
				continue
			}

			fmt.Print("Password: ")

			passwordByte, err := gopass.GetPasswdMasked()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(Login) > 20 || len(Login) < 5 {
				fmt.Println("Error Password")
				continue
			}
			Password = string(passwordByte)
			match, err = regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]*$", Password)
			if err != nil {
				fmt.Println("Error Password")
				continue
			}
			if !match {
				fmt.Println("Enter Password") //todo
				continue
			}
			break
		}
	}
	return Login, Password
}
