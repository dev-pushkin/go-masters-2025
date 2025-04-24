package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log/slog"

	"github.com/go_course_master/homework/hw_01/internal/errs"
)

// App - структура приложения
// Приложение шифрует и дешифрует данные
type App struct {
	l *slog.Logger

	// rndStrFn - функция генерации случайной строки длиной 32 байта
	// Экспортируется для возможности замены любой функцией генерации случайной строки
	RndStrFn func() string
}

func NewApp(l *slog.Logger) *App {
	return &App{
		l: l,
	}
}

type EncryptData struct {
	Secret string `json:"secret"`
	Data   string `json:"data"`
}

// Encrypt - метод шифрования данных
// Принимает строку и секрет(опционально) и возвращает зашифрованную строку и секрет
func (a *App) Encrypt(secret, data string) (*EncryptData, error) {
	if data == "" {
		return nil, errs.ErrIcomingDataIsEmpty
	}
	if secret == "" {
		secret = a.RndStrFn()
	}
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encData := gcm.Seal(nonce, nonce, []byte(data), nil)
	return &EncryptData{
		Secret: secret,
		Data:   base64.StdEncoding.EncodeToString(encData),
	}, nil

}

func (a *App) Decrypt(secretKey, encData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return "", errs.ErrIcomingDataIsEmpty
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errs.NewApiError(400, "ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
