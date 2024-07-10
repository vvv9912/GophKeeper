package authorization

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestSha256Hash(t *testing.T) {
	input := []byte("Hello, World!")
	expectedHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"

	result := Sha256Hash(string(input))

	if !reflect.DeepEqual(result, string(expectedHash)) {
		t.Errorf("Expected hash %s, but got %s", expectedHash, result)
	}
}

func TestBuildJWTString(t *testing.T) {
	// Создаем объект Autorization
	auth := &Autorization{
		SecretKey: "supersecretkey",
		TokenExp:  time.Hour,
	}

	// Вызываем функцию BuildJWTString с userId = 123
	tokenString, err := auth.BuildJWTString(123)
	if err != nil {
		t.Errorf("Error building JWT token: %v", err)
	}

	// Проверяем, что строка токена не пустая
	if tokenString == "" {
		t.Error("Empty token string")
	}
}
func TestBuildJWTStringBad(t *testing.T) {
	// Создаем объект Autorization
	auth := &Autorization{
		SecretKey: "",
		TokenExp:  time.Hour,
	}

	// Вызываем функцию BuildJWTString с userId = 123
	tokenString, err := auth.BuildJWTString(123)
	if err != nil {
		t.Errorf("Error building JWT token: %v", err)
	}

	// Проверяем, что строка токена не пустая
	if tokenString == "" {
		t.Error("Empty token string")
	}
}

func TestGetUserId(t *testing.T) {
	// Создаем объект Autorization
	auth := &Autorization{
		SecretKey: "supersecretkey",
		TokenExp:  time.Hour,
	}

	// Создаем токен
	token := jwt.New(jwt.SigningMethodHS256)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.TokenExp)),
		},
		UserID: 456,
	}
	token.Claims = claims

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(auth.SecretKey))
	if err != nil {
		t.Errorf("Error signing token: %v", err)
	}

	// Получаем userId из токена
	userId, err := auth.GetUserId(tokenString)
	if err != nil {
		t.Errorf("Error getting userId: %v", err)
	}

	// Проверяем, что userId совпадает с ожидаемым значением
	if userId != 456 {
		t.Errorf("Incorrect userId. Expected: 456, Got: %d", userId)
	}
}
func TestGetUserIdBad(t *testing.T) {
	// Создаем объект Autorization
	auth := &Autorization{
		SecretKey: "supersecretkey",
		TokenExp:  time.Hour,
	}

	// Создаем токен
	token := jwt.New(jwt.SigningMethodHS256)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.TokenExp)),
		},
		UserID: 456,
	}
	token.Claims = claims

	// Подписываем токен
	_, err := token.SignedString(nil)
	assert.Error(t, err)

}
