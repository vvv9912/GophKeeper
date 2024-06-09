package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

//	func (c *Cobra) CreateCredentials(cmd *cobra.Command, args []string) {
//		c.s.AuthService
//	}
func (c *Cobra) CreateFile(cmd *cobra.Command, args []string) {
	var (
		Path        string
		Name        string
		Description string
	)
	fmt.Println(args)
	if len(args) == 3 {
		Path = args[0]
		Name = args[1]
		Description = args[2]
	} else {
		fmt.Println("Error: Invalid number of arguments")
		return
	}
	ch := make(chan string)
	go func() {
		{
			for {
				val, ok := <-ch
				if ok {
					fmt.Println(val)
				} else {
					fmt.Println("Канал закрыт")
				}
			}
		}
	}()
	defer close(ch)
	err := c.s.CreateFile(cmd.Context(), Path, Name, Description, ch)
	time.Sleep(1 * time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

// Получение списка данных
func (c *Cobra) GetListData(cmd *cobra.Command, args []string) {

	resp, err := c.s.GetListData(cmd.Context())
	if err != nil {
		fmt.Println(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, resp, "", "  ")
	fmt.Println(out.String())
}

//func (c *Cobra) CreateCredentials(cmd *cobra.Command, args []string) {
//	var (
//		Name        string
//		Description string
//	)
//	if len(args) == 2 {
//		Name = args[0]
//		Description = args[1]
//	} else {
//		fmt.Println("Error: Invalid number of arguments")
//		return
//	}
//	err := c.s.CreateCredentials(cmd.Context(), Name, Description)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
