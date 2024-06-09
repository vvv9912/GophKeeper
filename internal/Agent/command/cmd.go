package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

// todo
func (c *Cobra) Demon(cmd *cobra.Command, args []string) {
	fmt.Println("N called")
	// Записываем во временный файлик pid или в sqlite
	// при завершении процессах убираем этот  pid
	//слздаем signal который убьет процесс os.Signal
	jwt, err := c.s.SignIn(cmd.Context(), "login", "passw")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jwt)
}
