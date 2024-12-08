package presenter

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseSetCookie(setCookie string) map[string]string {
	attributes := make(map[string]string)
	parts := strings.Split(setCookie, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// "key=value"形式の属性を解析
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) == 2 {
			attributes[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
		} else {
			attributes[strings.TrimSpace(keyValue[0])] = ""
		}
	}
	return attributes
}

func TestOutputUpdateResult(t *testing.T) {
	/* Arrange */
	expected := "{}\n"
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputUpdateResult()

	/* Assert */
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
}

func TestOutputLoginResult(t *testing.T) {
	/* Arrange */
	expected := "{}\n"
	token := "test_token"
	os.Setenv("JWT_TOKEN_NAME", "auth_token")
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputLoginResult(token)

	/* Assert */
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
	// レスポンスヘッダーからSet-Cookieを取得
	setCookie := rec.Header().Get("Set-Cookie")
	cookieAttributes := parseSetCookie(setCookie)
	assert.Equal(t, token, cookieAttributes[os.Getenv("JWT_TOKEN_NAME")])
}

func TestOutputAuthUrl(t *testing.T) {
	/* Arrange */
	url := "https://www.google.com"
	expected := url
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputAuthUrl(url)

	/* Assert */
	// up.OutputAuthUrlがJSONを返すこと
	if assert.NoError(t, actual) {
		assert.Equal(t, http.StatusFound, rec.Code)
		// リダイレクト先のURLが正しいこと
		assert.Equal(t, expected, rec.HeaderMap["Location"][0])
	}
}

func TestOutputSignupWithAuth(t *testing.T) {
	/* Arrange */
	t.Setenv("JWT_TOKEN_NAME", "auth_token")
	requestPath := "/editUser"
	token := "test_token"
	var expected error = nil
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputSignupWithAuth(token)

	/* Assert */
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Contains(t, rec.HeaderMap["Location"][0], requestPath)
	assert.Equal(t, expected, actual)

	// レスポンスヘッダーからSet-Cookieを取得
	setCookie := rec.Header().Get("Set-Cookie")
	cookieAttributes := parseSetCookie(setCookie)
	assert.Equal(t, token, cookieAttributes[os.Getenv("JWT_TOKEN_NAME")])
}

func TestOutputAlreadySignedup(t *testing.T) {
	/* Arrange */
	var expected error = nil
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputAlreadySignedup()

	/* Assert */
	if assert.NoError(t, actual) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, expected, actual)
	}
}

func TestOutputHasEmailInRequestBody(t *testing.T) {
	/* Arrange */
	expected := "{\"error\":\"Email is included in Request Body\"}\n"
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputHasEmailInRequestBody()

	/* Assert */
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
}
