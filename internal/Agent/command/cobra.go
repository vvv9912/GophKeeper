package command

import (
	"GophKeeper/internal/Agent/service"
	"GophKeeper/pkg/logger"
	"context"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Cobra struct {
	s       *service.Service
	rootCmd *cobra.Command
}

func NewCobra(s *service.Service) *Cobra {

	return &Cobra{s: s, rootCmd: &cobra.Command{}}
}

func (c *Cobra) Start(ctx context.Context) error {

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
		Use:     "createBinaryFile",
		Short:   "Create binary file",
		Example: "createFile {Path} {Name} {Description}",
		Args:    cobra.MaximumNArgs(3),
		Run:     c.CreateBinaryFile,
		Aliases: []string{"createbinaryfile"},
	}

	CreateCredentials := &cobra.Command{
		Use:     "createCredentials",
		Short:   "Create credentials",
		Example: "createCredentials {Name} {Description} {Login} {Password} ",
		Args:    cobra.MaximumNArgs(4),
		Run:     c.CreateCredentials,
		Aliases: []string{"createcredentials"},
	}

	CreateCreditCard := &cobra.Command{
		Use:     "createCreditCard",
		Short:   "Create credit card",
		Example: "createCreditCard {Name} {Description} {Name} {ExpireAt} {CardNumber} {CVV}",
		Args:    cobra.MaximumNArgs(6),
		Run:     c.CreateCreditCard,
		Aliases: []string{"createcreditcard"},
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

	UpdateCredentials := &cobra.Command{
		Use:     "updateCredentials",
		Short:   "Update credentials",
		Example: "updateCredentials {UserDataId} {Login} {Password}",
		Args:    cobra.MaximumNArgs(3),
		Run:     c.UpdateCredentials,
		Aliases: []string{"updatecredentials"},
	}

	UpdateCreditCard := &cobra.Command{
		Use:     "updateCreditCard",
		Short:   "Update credit card",
		Example: "updateCreditCard {UserDataId} {ExpireAt} {CardNumber} {CVV}",
		Args:    cobra.MaximumNArgs(5),
		Run:     c.UpdateCreditCard,
		Aliases: []string{"updatecreditcard"},
	}

	UpdateBinaryFile := &cobra.Command{
		Use:     "updateBinaryFile",
		Short:   "Update binary file",
		Example: "updateBinaryFile {UserDataId} {Path} ",
		Args:    cobra.MaximumNArgs(2),
		Run:     c.UpdateBinaryFile,
		Aliases: []string{"updatebinaryfile"},
	}

	c.rootCmd.AddCommand(signIn)
	c.rootCmd.AddCommand(signUp)
	c.rootCmd.AddCommand(CreateFile)
	c.rootCmd.AddCommand(ListData)
	c.rootCmd.AddCommand(GetData)
	c.rootCmd.AddCommand(CreateCredentials)
	c.rootCmd.AddCommand(CreateCreditCard)
	c.rootCmd.AddCommand(UpdateCredentials)
	c.rootCmd.AddCommand(UpdateCreditCard)
	c.rootCmd.AddCommand(UpdateBinaryFile)
	c.rootCmd.SetContext(ctx)
	if err := c.rootCmd.Execute(); err != nil {
		logger.Log.Error("Root execute err", zap.Error(err))
		return err
	}
	return nil
}
