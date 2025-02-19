package presenter

import (
	"clean-storemap-api/src/entity/errors"
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

func TestHandleError(t *testing.T) {
	/* Arrange */
	testCases := []struct {
		input        error
		expectedCode int
		expected     string
	}{
		{input: errors.ErrUserNotFound, expectedCode: http.StatusNotFound, expected: "{\"error\":\"Userが登録されていません。\"}\n"},
		{input: errors.ErrUserAlreadyExists, expectedCode: http.StatusConflict, expected: "{\"error\":\"そのemailはすでに登録されています。\"}\n"},
		{input: errors.ErrUserInvalidData, expectedCode: http.StatusBadRequest, expected: "{\"error\":\"不正なデータが含まれています。\"}\n"},
		{input: errors.ErrInvalidEmail, expectedCode: http.StatusBadRequest, expected: "{\"error\":\"Emailが適した形ではありません。\"}\n"},
		{input: errors.ErrAgeBelowZero, expectedCode: http.StatusBadRequest, expected: "{\"error\":\"年齢が0未満です。\"}\n"},
		{input: errors.ErrUnAuthorized, expectedCode: http.StatusUnauthorized, expected: "{\"error\":\"認証に失敗しました。\"}\n"},
		{input: errors.ErrUnsupported, expectedCode: http.StatusInternalServerError, expected: "{\"error\":\"予期しないエラーが発生しました\"}\n"},
	}
	/* Act, Assert */
	for _, tc := range testCases {
		// recが1ループごとに切れないのでrouterは毎回作り直す。
		c, rec := newRouter()
		up := &UserPresenter{c: c}
		actual := up.HandleError(tc.input)
		if assert.NoError(t, actual) {
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expected, rec.Body.String())
		}
	}
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
