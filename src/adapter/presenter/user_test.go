package presenter

import (
	"fmt"
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

func TestOutputLoginWithAuth(t *testing.T) {
	/* Arrange */
	t.Setenv("JWT_TOKEN_NAME", "auth_token")
	token := "test_token"
	var expected error = nil
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputLoginWithAuth(token)

	/* Assert */
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, actual)

	// レスポンスヘッダーからSet-Cookieを取得
	setCookie := rec.Header().Get("Set-Cookie")
	cookieAttributes := parseSetCookie(setCookie)
	assert.Equal(t, token, cookieAttributes[os.Getenv("JWT_TOKEN_NAME")])
}
func TestOutputAuthUrl(t *testing.T) {
	/* Arrange */
	url := "https://www.google.com"
	expected := fmt.Sprintf("{\"url\":\"%s\"}\n", url)

	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputAuthUrl(url)

	/* Assert */
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOutputSignupWithAuth(t *testing.T) {
	/* Arrange */
	t.Setenv("JWT_TOKEN_NAME", "auth_token")
	token := "test_token"
	var expected error = nil
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputSignupWithAuth(token)

	/* Assert */
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, actual)

	// レスポンスヘッダーからSet-Cookieを取得
	setCookie := rec.Header().Get("Set-Cookie")
	cookieAttributes := parseSetCookie(setCookie)
	assert.Equal(t, token, cookieAttributes[os.Getenv("JWT_TOKEN_NAME")])
}

func TestOutputAlreadySignedup(t *testing.T) {
	/* Arrange */
	expected := "{\"error\":\"Already exist favorite store\"}\n"
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputAlreadySignedup()

	/* Assert */
	assert.Equal(t, http.StatusConflict, rec.Code)
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
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

func TestOutputNotRegistered(t *testing.T) {
	/* Arrange */
	error := "not_registered"
	expected := fmt.Sprintf("{\"error\":\"%s\"}\n", error)
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputNotRegistered()

	/* Assert */
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
	assert.Equal(t, http.StatusNotFound, rec.Code)

}
