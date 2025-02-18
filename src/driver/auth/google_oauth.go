package auth

import (
	"clean-storemap-api/src/entity/errors"
	"context"
	"os"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthDriver struct{}

func NewGoogleOAuthDriver() *GoogleOAuthDriver {
	return &GoogleOAuthDriver{}
}

func newGoogleOauthConfig() *oauth2.Config {
	// Google認証を行うためのリダイレクト先のURL
	redirectURL := os.Getenv("BACKEND_URL") + "/auth/signup"
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	conf := &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "email"}, // emailを取得するためのスコープ
		Endpoint:     google.Endpoint,
	}
	return conf
}

func (oauth *GoogleOAuthDriver) GenerateUrl() string {
	//	認証情報を取得
	config := newGoogleOauthConfig()
	// URLの生成
	return config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// GoogleのOAuth認証を行い、ユーザー情報を取得する
func getProfile(code string) (map[string]interface{}, error) {
	config := newGoogleOauthConfig()
	ctx := context.Background()
	// 認証情報を取得
	oauth2Token, err := config.Exchange(ctx, code)
	if err != nil {
		return make(map[string]interface{}), err
	}
	// IDトークンの取得
	rawIdToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return make(map[string]interface{}), err
	}
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return make(map[string]interface{}), err
	}
	// IDトークンの検証
	idToken, err := provider.Verifier(&oidc.Config{ClientID: config.ClientID}).Verify(ctx, rawIdToken)
	if err != nil {
		return make(map[string]interface{}), err
	}
	// ユーザー情報の取得
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return make(map[string]interface{}), err
	}
	return profile, nil
}

func (oauth *GoogleOAuthDriver) GetEmail(code string) (string, error) {
	profile, err := getProfile(code)
	if err != nil {
		return "", errors.ErrGetEmailFailed
	}
	return profile["email"].(string), nil
}
