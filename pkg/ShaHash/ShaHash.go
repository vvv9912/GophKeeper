package ShaHash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hash(input []byte) string {
	// Создаем новый хеш SHA-256
	hasher := sha256.New()

	// Преобразуем строку в байты и передаем хеш-функции
	hasher.Write(input)

	// Получаем хеш в виде среза байтов
	hashBytes := hasher.Sum(nil)

	// Преобразуем срез байтов в строку в шестнадцатеричном формате
	hashedString := hex.EncodeToString(hashBytes)

	return hashedString
}
