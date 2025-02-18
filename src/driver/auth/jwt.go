package auth

import (
	"clean-storemap-api/src/entity/errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtDriver struct{}

func NewJwtDriver() *JwtDriver {
	return &JwtDriver{}
}

func (auth *JwtDriver) GenerateToken(subject string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "crean-storemap"
	claims["sub"] = subject
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 有効期限は24時間
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if err != nil {
		return "", errors.ErrFailedGenerateToken
	}
	return tokenString, nil
}
