package presenter

import (
	"clean-storemap-api/src/entity/errors"
	"clean-storemap-api/src/usecase/port"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type UserPresenter struct {
	c echo.Context
}

func NewUserOutputPort(c echo.Context) port.UserOutputPort {
	return &UserPresenter{c: c}
}

func (up *UserPresenter) HandleError(err error) error {
	if errors.Is(err, errors.ErrUserNotFound) {
		return up.c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Userが登録されていません。"})
	} else if errors.Is(err, errors.ErrUserAlreadyExists) {
		return up.c.JSON(http.StatusConflict, map[string]interface{}{"error": "そのemailはすでに登録されています。"})
	} else if errors.Is(err, errors.ErrInvalidUserData) {
		return up.c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "不正なデータが含まれています。"})
	} else if errors.Is(err, errors.ErrInvalidEmail) {
		return up.c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Emailが適した形ではありません。"})
	} else if errors.Is(err, errors.ErrAgeBelowZero) {
		return up.c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "年齢が0未満です。"})
	}

	// 認証関連のエラー
	if errors.Is(err, errors.ErrUnAuthorizedForUser) {
		return up.c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "認証に失敗しました。"})
	}
	return up.c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "予期しないエラーが発生しました"})
}

func (up *UserPresenter) OutputUpdateResult() error {
	return up.c.JSON(http.StatusOK, map[string]interface{}{})
}

func (up *UserPresenter) OutputLoginResult(token string) error {
	cookie := createAuthCookie(token)
	up.c.SetCookie(cookie)
	return up.c.JSON(http.StatusOK, map[string]interface{}{})
}

func (up *UserPresenter) OutputAuthUrl(url string) error {
	return up.c.Redirect(http.StatusFound, url)
}

func (up *UserPresenter) OutputSignupWithAuth(token string) error {
	url := os.Getenv("FRONT_URL") + "/editUser" // 認証以外のユーザ情報を入力するページ
	cookie := createAuthCookie(token)
	up.c.SetCookie(cookie)
	return up.c.Redirect(http.StatusFound, url)
}

func (up *UserPresenter) OutputAlreadySignedup() error {
	url := os.Getenv("FRONT_URL") // すでに登録済みの場合はトップページにリダイレクト
	return up.c.Redirect(http.StatusFound, url)
}

func (up *UserPresenter) OutputHasEmailInRequestBody() error {
	errMsg := "Email is included in Request Body"
	return up.c.JSON(http.StatusBadRequest, map[string]interface{}{"error": errMsg})
}

func createAuthCookie(token string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = os.Getenv("JWT_TOKEN_NAME")
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour) // 24時間有効
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode // クロスサイトリクエストを許可
	cookie.HttpOnly = false
	cookie.Secure = false
	return cookie
}
