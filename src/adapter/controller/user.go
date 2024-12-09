package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"clean-storemap-api/src/adapter/gateway"
	accssesTypel "clean-storemap-api/src/entity"
	"clean-storemap-api/src/usecase/port"
)

type UserI interface {
	UpdateUser(c echo.Context) error
	LoginWithAuth(c echo.Context) error
	GetAuthUrl(c echo.Context) error
	SignupWithAuth(c echo.Context) error
}

type UserOutputFactory func(echo.Context) port.UserOutputPort
type UserInputFactory func(port.UserRepository, port.UserOutputPort) port.UserInputPort
type UserRepositoryFactory func(gateway.UserDriver, gateway.GoogleOAuthDriver, gateway.JwtDriver) port.UserRepository
type UserDriverFactory gateway.UserDriver
type GoogleOAuthDriverFactory gateway.GoogleOAuthDriver
type JwtDriverFactory gateway.JwtDriver

type UserController struct {
	userDriverFactory        UserDriverFactory
	googleOAuthDriverFactory GoogleOAuthDriverFactory
	jwtDriverFactory         JwtDriverFactory
	userOutputFactory        UserOutputFactory
	userInputFactory         UserInputFactory
	userRepositoryFactory    UserRepositoryFactory
}

func NewUserController(
	userDriverFactory UserDriverFactory,
	googleOAuthDriverFactory GoogleOAuthDriverFactory,
	jwtDriverFactory JwtDriverFactory,
	userOutputFactory UserOutputFactory,
	userInputFactory UserInputFactory,
	userRepositoryFactory UserRepositoryFactory,
) UserI {
	return &UserController{
		userDriverFactory:        userDriverFactory,
		googleOAuthDriverFactory: googleOAuthDriverFactory,
		jwtDriverFactory:         jwtDriverFactory,
		userOutputFactory:        userOutputFactory,
		userInputFactory:         userInputFactory,
		userRepositoryFactory:    userRepositoryFactory,
	}
}

func (uc *UserController) UpdateUser(c echo.Context) error {
	id := c.Get("userId").(string)
	// UserRequestBodyを使用すると存在しないkeyに関しても値が生成されてしまうため、UserRequestBodyにバインドさせずに取得する
	var requestBody map[string]interface{}
	// 数字 -> float64, 文字列-> stringと変換される
	if err := json.NewDecoder(c.Request().Body).Decode(&requestBody); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	updateData := make(accssesTypel.ChangeForUser)
	// 更新データを型変換しつつ格納する

	// name
	if name, ok := requestBody["name"]; ok {
		updateData["name"] = name
	}

	// email
	if email, ok := requestBody["email"]; ok {
		updateData["email"] = email
	}

	// age
	if age, ok := requestBody["age"].(string); ok {
		updateData["age"], _ = strconv.Atoi(age)
	} else if age, ok := requestBody["age"].(float64); ok {
		updateData["age"] = int(age)
	} else {
		updateData["age"] = requestBody["age"]
	}

	// sex
	if sex, ok := requestBody["sex"].(string); ok {
		updateData["sex"], _ = strconv.ParseFloat(sex, 32)
	} else if sex, ok := requestBody["sex"].(float64); ok {
		updateData["sex"] = float32(sex)
	} else {
		updateData["sex"] = requestBody["sex"]
	}

	// gender
	if gender, ok := requestBody["gender"].(string); ok {
		updateData["gender"], _ = strconv.ParseFloat(gender, 32)
	} else if gender, ok := requestBody["gender"].(float64); ok {
		updateData["gender"] = float32(gender)
	} else {
		updateData["gender"] = requestBody["gender"]
	}

	return uc.newUserInputPort(c).UpdateUser(id, updateData)
}

func (uc *UserController) LoginWithAuth(c echo.Context) error {
	codeParameter := c.QueryParam("code")
	return uc.newUserInputPort(c).LoginUser(codeParameter)
}

func (uc *UserController) GetAuthUrl(c echo.Context) error {
	accessedType := c.QueryParam("accessedType")
	if accessedType != "signup" && accessedType != "login" {
		return c.JSON(http.StatusAccepted, map[string]interface{}{"error": "accessedType is not correct"})
	}
	return uc.newUserInputPort(c).GetAuthUrl(accessedType)
}

func (uc *UserController) SignupWithAuth(c echo.Context) error {
	codeParameter := c.QueryParam("code") // パラメータの取得
	return uc.newUserInputPort(c).SignupDraft(codeParameter)
}

/* ここでpresenterにecho.Contextを渡している！起爆！！！（遅延） */
/* これによって、presenterのinterface(outputport)にecho.Contextを書かなくて良くなる */
func (uc *UserController) newUserInputPort(c echo.Context) port.UserInputPort {
	userOutputPort := uc.userOutputFactory(c)
	userDriver := uc.userDriverFactory
	googleOAuthDriver := uc.googleOAuthDriverFactory
	jwtDriver := uc.jwtDriverFactory
	userRepository := uc.userRepositoryFactory(userDriver, googleOAuthDriver, jwtDriver)
	return uc.userInputFactory(userRepository, userOutputPort)
}
