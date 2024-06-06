package command

import (
	"GophKeeper/internal/Agent/service"
	"GophKeeper/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Cobra struct {
	s *service.Service
}

func NewCobra(s *service.Service) *Cobra {
	return &Cobra{s: s}
}

func (c *Cobra) Start() error {
	rootCmd := &cobra.Command{
		Use:   "myapp",
		Short: "My Application",
	}

	signIn := &cobra.Command{
		Use:     "signIn",
		Short:   "Sign in",
		Example: "signIn {Login} {Password}",
		Args:    cobra.MaximumNArgs(2),
		Run:     c.SignIn,
		Aliases: []string{"signin"},
	}
	signUp := &cobra.Command{
		Use:     "signUp",
		Short:   "Sign up",
		Example: "signUp {Login} {Password}",
		Args:    cobra.MaximumNArgs(2),
		Run:     c.SignUp,
		Aliases: []string{"signup"},
	}
	//c2 := &command.Command{}

	rootCmd.AddCommand(signIn)
	rootCmd.AddCommand(signUp)
	//rootCmd.AddCommand(c2)
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Error("Root execute err", zap.Error(err))
		return err
	}
	return nil
}
