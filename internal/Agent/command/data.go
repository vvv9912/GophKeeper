package command

import (
	"GophKeeper/internal/Agent/model"
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/logger"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strconv"
)

// CreateBinaryFile - создание бинарного файла
func (c *Cobra) CreateBinaryFile(cmd *cobra.Command, args []string) {
	var (
		Path        string
		Name        string
		Description string
	)

	if len(args) == 3 {
		Path = args[0]
		Name = args[1]
		Description = args[2]
	} else {
		logger.Log.Error("Error: Invalid number of arguments")
		return
	}
	ch := make(chan string)
	go func() {
		{
			for {
				val, ok := <-ch
				if !ok {
					return
				}
				logger.Log.Info(val)
			}
		}
	}()
	defer close(ch)
	err := c.s.CreateBinaryFile(cmd.Context(), Path, Name, Description, ch)

	if err != nil {
		logger.Log.Error("CreateBinaryFile failed", zap.Error(err))
	}

}

// CreateCredentials - создание данных логин/пароль
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
		logger.Log.Error("Error: Invalid number of arguments")
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

	logger.Log.Info("CreateCredentials success")
}

// CreateCreditCard - создание кредитной карты
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
		logger.Log.Error("CreateCreditCard failed", zap.Error(err))
	}

	logger.Log.Info("CreateCreditCard success")
}

// GetListData - Получение списка данных
func (c *Cobra) GetListData(cmd *cobra.Command, args []string) {

	resp, err := c.s.GetListData(cmd.Context())
	if err != nil {
		logger.Log.Error("GetListData failed", zap.Error(err))
	}

	var out bytes.Buffer
	err = json.Indent(&out, resp, "", "  ")

}

// GetData - Получение списка данных
func (c *Cobra) GetData(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logger.Log.Error("Error: Invalid number of arguments")
		return
	}
	userDataId, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Log.Error("Error: Invalid number of arguments", zap.Error(err))
		return
	}
	resp, err := c.s.GetData(cmd.Context(), int64(userDataId))

	if err != nil {
		logger.Log.Error("GetData failed", zap.Error(err))
	}

	logger.Log.Info(string(resp))

}

// UpdateCredentials - обновление данных логин/пароль
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
		logger.Log.Error("Error: Invalid number of arguments")
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

	logger.Log.Info(string(resp))
	logger.Log.Info("Credentials updated successfully")
}

// UpdateCreditCard - обновление данных кредитной карточки
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

	logger.Log.Info(string(resp))
	logger.Log.Info("CreditCard updated successfully")
}

// UpdateBinaryFile - обновление бинарного файл
func (c *Cobra) UpdateBinaryFile(cmd *cobra.Command, args []string) {
	var (
		UserDataId string
		Path       string
	)

	logger.Log.Debug("args", zap.Any("args", args))
	//
	if len(args) == 2 {
		UserDataId = args[0]
		Path = args[1]
	} else {
		logger.Log.Error("Error: Invalid number of arguments")
		return
	}
	ch := make(chan string)
	go func() {
		{
			for {
				val, ok := <-ch
				if !ok {
					return
				}
				logger.Log.Info(val)
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
		logger.Log.Error("UpdateBinaryFile failed", zap.Error(err))
	}

	logger.Log.Info("BinaryFile updated successfully")
}
