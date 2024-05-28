package handler

type Auth struct {
	Login    string `json:"login" `
	Password string `json:"password" `
}

type User struct {
	Login string `json:"login" `
	JWT   string `json:"jwt"`
}

type RespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
