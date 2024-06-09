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

	CreateFile := &cobra.Command{
		Use:     "createFile",
		Short:   "Create file",
		Example: "createFile {Path} {Name} {Description}",
		Args:    cobra.MaximumNArgs(3),
		Run:     c.CreateFile,
		Aliases: []string{"createfile"},
	}

	ListData := &cobra.Command{
		Use:     "listData",
		Short:   "List data",
		Example: "listData",
		Args:    cobra.MaximumNArgs(0),
		Run:     c.GetListData,
		Aliases: []string{"listdata"},
	}

	GetData := &cobra.Command{
		Use:     "getData",
		Short:   "Get data",
		Example: "getData",
		Args:    cobra.MaximumNArgs(1),
		Run:     c.GetData,
		Aliases: []string{"getdata"},
	}
	//c2 := &command.Command{}

	rootCmd.AddCommand(signIn)
	rootCmd.AddCommand(signUp)
	rootCmd.AddCommand(CreateFile)
	rootCmd.AddCommand(ListData)
	rootCmd.AddCommand(GetData)
	//rootCmd.AddCommand(c2)
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Error("Root execute err", zap.Error(err))
		return err
	}
	return nil
}
