package model

// Credentials содержит информацию о логине и пароле пользователя.
type Credentials struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}

// CreditCard содержит информацию о кредитной карте.
type CreditCard struct {
	Name       string `json:"name"`       // Имя владельца карты
	ExpireAt   int    `json:"expireAt"`   // Срок действия карты (TODO: изменить на uint8)
	CardNumber int64  `json:"cardNumber"` // Номер карты (TODO: изменить на uint)
	CVV        int8   `json:"cvv"`        // CVV код карты
}

// Data содержит зашифрованные данные.
type Data struct {
	DecryptData []byte `json:"decryptData"` // Расшифрованные данные
}
