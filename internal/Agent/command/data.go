package command

import (
	"github.com/spf13/cobra"
)

func (c *Cobra) CreateCredentials(cmd *cobra.Command, args []string) {
	c.s.AuthService
}
