package handler

// Auth - данные для аутентификации.
type Auth struct {
	Login    string `json:"login" `    // логин.
	Password string `json:"password" ` // пароль.
}

// User - данные пользователя.
type User struct {
	Login string `json:"login"` // логин.
	JWT   string `json:"jwt"`   // токен.
}

// RespError - структура ошибки.
type RespError struct {
	Code    int    `json:"code"`    // код ошибки.
	Message string `json:"message"` // текст ошибки.
}

// ReqData - данные запроса.
type ReqData struct {
	Name        string `json:"name"`        // название.
	Description string `json:"description"` // описание.
	Data        []byte `json:"data"`        // данные.
}
