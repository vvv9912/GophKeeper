package authorization

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Autorization struct {
	TokenExp  time.Duration
	SecretKey string
}

func NewAutorization(tokenExp time.Duration, secretKey string) *Autorization {
	return &Autorization{TokenExp: tokenExp, SecretKey: secretKey}
}

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func (a *Autorization) BuildJWTString(userId int64) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.TokenExp)),
		},
		// собственное утверждение
		UserID: userId,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(a.SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// GetUserId - возвращает id пользователя из jwt
func (a *Autorization) GetUserId(tokenString string) (int64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			//доп проверка заголовка
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(a.SecretKey), nil
		})
	if err != nil {
		return -1, nil
	}

	if !token.Valid {
		//		fmt.Println("Token is not valid")
		return -1, nil
	}

	return claims.UserID, nil
}
func Sha256Hash(input string) string {
	// Создаем новый хеш SHA-256
	hasher := sha256.New()

	// Преобразуем строку в байты и передаем хеш-функции
	hasher.Write([]byte(input))

	// Получаем хеш в виде среза байтов
	hashBytes := hasher.Sum(nil)

	// Преобразуем срез байтов в строку в шестнадцатеричном формате
	hashedString := hex.EncodeToString(hashBytes)

	return hashedString
}
