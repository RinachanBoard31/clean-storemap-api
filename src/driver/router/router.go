package router

import (
	controller "clean-storemap-api/src/adapter/controller"
	"clean-storemap-api/src/driver/middleware"

	"context"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

// Validationのために必要なメソッド
type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type RouterI interface {
	Serve(ctx context.Context)
}

type Router struct {
	echo            *echo.Echo
	storeController controller.StoreI
	userController  controller.UserI
}

func NewRouter(echo *echo.Echo, storeController controller.StoreI, userController controller.UserI) RouterI {
	return &Router{
		echo:            echo,
		storeController: storeController,
		userController:  userController,
	}
}

func (router *Router) Serve(ctx context.Context) {
	// ログイン前のルーティング
	router.echo.GET("/", router.storeController.GetStores)
	router.echo.GET("/auth", router.userController.GetAuthUrl)            // Google認証用のURLを取得し返す(?accessedType=login, signup)
	router.echo.GET("/auth/signup", router.userController.SignupWithAuth) // ユーザの認証を確認し仮登録する
	router.echo.GET("/auth/login", router.userController.LoginWithAuth)

	// ログイン後のルーティング(認証が必要なパスはここより下に書く)
	// 認証のためのJWTMiddlewareを設定
	secured := router.echo.Group("")
	secured.Use(middleware.JwtAuthMiddleware())

	secured.GET("/stores/opening-hours", router.storeController.GetNearStores)
	secured.GET("/stores/favorite-ranking", router.storeController.GetTopFavoriteStores)
	secured.GET("/user/favorite-store", router.storeController.GetFavoriteStores)
	secured.POST("/user/favorite-store", router.storeController.SaveFavoriteStore)
	secured.PUT("/user", router.userController.UpdateUser)
	router.echo.Logger.Fatal(router.echo.Start(":8080"))
}
