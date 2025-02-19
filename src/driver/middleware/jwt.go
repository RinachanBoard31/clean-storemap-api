package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func JwtAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(os.Getenv("JWT_TOKEN_NAME"))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "不正なアクセスです。",
				})
			}
			claims, err := validateJwtOnSelf(cookie)
			if err != nil {
				fmt.Printf("Request is %s: %s \nError: %s\n", c.Request().Method, c.Request().URL, err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "不正なアクセスです。",
				})
			}
			userId := claims["sub"].(string)
			c.Set("userId", userId)
			return next(c)
		}
	}
}

// echoのJWTで認証も可能であるが、合っているはずなのにtoken.Validでエラーが出る時があるため手動で確かめる
func validateJwtOnSelf(cookie *http.Cookie) (jwt.MapClaims, error) {
	tokenString := cookie.Value
	jwtData := strings.Split(tokenString, ".")
	if len(jwtData) != 3 { // ヘッダー、ペイロード、シグネチャの3つに分かれているか確認
		return nil, fmt.Errorf("token is not included header, payload or signature")
	}
	header, err := decodeBase64Url(jwtData[0])
	if err != nil {
		return nil, err
	}
	payload, err := decodeBase64Url(jwtData[1])
	if err != nil {
		return nil, err
	}
	var headerData map[string]interface{}
	var payloadData map[string]interface{}
	json.Unmarshal(header, &headerData)
	json.Unmarshal(payload, &payloadData)
	if payloadData["exp"] != nil { // 有効期限の確認
		expTime := int64(payloadData["exp"].(float64))
		if expTime < time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}
	}

	// ヘッダーとペイロードを結合
	encodedData := jwtData[0] + "." + jwtData[1]
	secret := []byte(os.Getenv("JWT_SIGNING_KEY"))
	signature := createSignature(encodedData, secret)
	if jwtData[2] == signature { // シグネチャが一致しているか確認
		claims := jwt.MapClaims(payloadData)
		return claims, nil
	}
	return nil, fmt.Errorf("authentication failure")
}

// Base64エンコード
func encodeBase64Url(input []byte) string {
	encoded := base64.StdEncoding.EncodeToString(input)
	// Base64Urlエンコード（+ → - , / → _ , = を削除）
	encoded = strings.TrimRight(encoded, "=")
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	return encoded
}

// Base64エンコード
func decodeBase64Url(input string) ([]byte, error) {
	if decoded, err := base64.RawURLEncoding.DecodeString(input); err != nil {
		return nil, err
	} else {
		return decoded, nil
	}
}

// HMAC-SHA256で署名を作成
func createSignature(data string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(data))
	return encodeBase64Url(mac.Sum(nil))
}
