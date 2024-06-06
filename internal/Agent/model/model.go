package model

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreditCard struct {
	Name       string `json:"name"`
	ExpireAt   int    `json:"expireAt"`   //todo uint8
	CardNumber int64  `json:"cardNumber"` //todo uint
	CVV        int8   `json:"cvv"`
}
type Data struct {
	DecryptData []byte `json:"decryptData"`
}
