package command

import (
	"GophKeeper/internal/Agent/model"
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strconv"
	"time"
)

//	func (c *Cobra) CreateCredentials(cmd *cobra.Command, args []string) {
//		c.s.AuthService
//	}
func (c *Cobra) CreateBinaryFile(cmd *cobra.Command, args []string) {
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

	if err != nil {
		fmt.Println(err)
	}
	//todo waitGroup?
	time.Sleep(1 * time.Second)
}

func (c *Cobra) CreateCredentials(cmd *cobra.Command, args []string) {
	var (
		Login       string
		Password    string
		Description string
		Name        string
	)
	if len(args) == 4 {
		Name = args[0]
		Description = args[1]
		Login = args[2]
		Password = args[3]

	} else {
		fmt.Println("Error: Invalid number of arguments")
		return
	}

	var cred model.Credentials

	cred.Login = Login
	cred.Password = Password

	credential, err := json.Marshal(cred)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return
	}

	err = c.s.CreateCredentials(cmd.Context(), &server.ReqData{
		Name:        Name,
		Description: Description,
		Data:        credential,
	})
	if err != nil {
		logger.Log.Error("CreateCredentials failed", zap.Error(err))
	}
	fmt.Println("Credentials created successfully")
}

func (c *Cobra) CreateCreditCard(cmd *cobra.Command, args []string) {

	var (
		Name        string
		Description string

		NameBank   string
		ExpAt      string
		CardNumber string
		Cvv        string
	)
	if len(args) == 6 {
		Name = args[0]
		Description = args[1]
		NameBank = args[2]
		CardNumber = args[3]
		ExpAt = args[4]
		Cvv = args[5]
	} else {
		logger.Log.Error("Error: Invalid number of arguments")
		return
	}

	var creditCard model.CreditCard
	creditCard.Name = NameBank
	cardNum, err := strconv.ParseInt(CardNumber, 10, 64)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	creditCard.CardNumber = cardNum

	expAt, err := strconv.Atoi(ExpAt)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	creditCard.ExpireAt = expAt

	cvv, err := strconv.ParseInt(Cvv, 10, 8)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	creditCard.CVV = int8(cvv)

	data, err := json.Marshal(creditCard)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return
	}

	err = c.s.CreateCreditCard(cmd.Context(), &server.ReqData{
		Name:        Name,
		Description: Description,
		Data:        data,
	})

	if err != nil {
		logger.Log.Error("CreateCreditCatd failed", zap.Error(err))
	}
	fmt.Println("Credit card created successfully")
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

// Получение списка данных
func (c *Cobra) GetData(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Error: Invalid number of arguments")
		return
	}
	userDataId, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	resp, err := c.s.GetData(cmd.Context(), int64(userDataId))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp))

}
func (c *Cobra) UpdateCredentials(cmd *cobra.Command, args []string) {
	var (
		UserDataId string
		Login      string
		Password   string
	)
	if len(args) == 3 {
		UserDataId = args[0]
		Login = args[1]
		Password = args[2]

	} else {
		fmt.Println("Error: Invalid number of arguments")
		return
	}

	var cred model.Credentials

	cred.Login = Login
	cred.Password = Password

	credential, err := json.Marshal(cred)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return
	}
	userDataId, err := strconv.Atoi(UserDataId)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	resp, err := c.s.UpdateData(cmd.Context(), int64(userDataId), credential)
	if err != nil {
		logger.Log.Error("CreateCredentials failed", zap.Error(err))
	}

	fmt.Println(string(resp))
	fmt.Println("Credentials updated successfully")
}

func (c *Cobra) UpdateCreditCard(cmd *cobra.Command, args []string) {

	var (
		UserDataId string
		NameBank   string
		ExpAt      string
		CardNumber string
		Cvv        string
	)
	if len(args) == 5 {
		UserDataId = args[0]
		NameBank = args[1]
		CardNumber = args[2]
		ExpAt = args[3]
		Cvv = args[4]
	}
	var creditCard model.CreditCard
	creditCard.Name = NameBank
	cardNum, err := strconv.ParseInt(CardNumber, 10, 64)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	creditCard.CardNumber = cardNum

	expAt, err := strconv.Atoi(ExpAt)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	creditCard.ExpireAt = expAt

	cvv, err := strconv.ParseInt(Cvv, 10, 8)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	creditCard.CVV = int8(cvv)

	data, err := json.Marshal(creditCard)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return
	}

	userDataId, err := strconv.Atoi(UserDataId)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	resp, err := c.s.UpdateData(cmd.Context(), int64(userDataId), data)
	if err != nil {
		logger.Log.Error("CreateCredentials failed", zap.Error(err))
	}

	fmt.Println(string(resp))
	fmt.Println("CreditCard update successfully")
}
func (c *Cobra) UpdateBinaryFile(cmd *cobra.Command, args []string) {
	var (
		UserDataId string
		Path       string
	)

	fmt.Println(args)
	if len(args) == 2 {
		UserDataId = args[0]
		Path = args[1]
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
					return
				}
			}
		}
	}()
	defer close(ch)
	userDataId, err := strconv.Atoi(UserDataId)
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}

	err = c.s.UpdateBinaryFile(cmd.Context(), Path, int64(userDataId), ch)

	if err != nil {
		fmt.Println(err)
	}

	//todo waitGroup?
	time.Sleep(1 * time.Second)
	fmt.Println("File updated successfully")
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
