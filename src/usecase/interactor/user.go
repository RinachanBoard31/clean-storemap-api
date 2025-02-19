package interactor

import (
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/entity/errors"
	port "clean-storemap-api/src/usecase/port"
)

type UserInteractor struct {
	userRepository port.UserRepository
	userOutputPort port.UserOutputPort
}

func NewUserInputPort(userRepository port.UserRepository, userOutputPort port.UserOutputPort) port.UserInputPort {
	return &UserInteractor{
		userRepository: userRepository,
		userOutputPort: userOutputPort,
	}
}

func (ui *UserInteractor) UpdateUser(id string, updateData model.ChangeForUser) error {
	// emailを更新しようとした場合にはエラーを返す
	if _, ok := updateData["email"]; ok {
		return ui.userOutputPort.OutputUserHandleError(errors.ErrUserInvalidData)
	}

	// 整形する
	if age, ok := updateData["age"].(int); ok {
		if err := model.AgeValid(age); err != nil {
			return ui.userOutputPort.OutputUserHandleError(err)
		}
		updateData["age"] = model.AgeFormat(age)
	}

	if sex, ok := updateData["sex"].(float32); ok {
		updateData["sex"] = model.SexFormat(sex)
	}

	if gender, ok := updateData["gender"].(float32); ok {
		updateData["gender"] = model.GenderFormat(gender)
	}

	// userが存在するか確認
	user, err := ui.userRepository.Get(id)
	if err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}

	if err := ui.userRepository.Update(user, updateData); err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}
	return ui.userOutputPort.OutputUpdateResult()
}

func (ui *UserInteractor) LoginUser(userCredentials *model.UserCredentials) error {
	user, err := ui.userRepository.FindBy(userCredentials)
	if err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}
	token, err := ui.userRepository.GenerateAccessToken(user.Id)
	if err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}
	// urlのクエリパラメータにidを付与してそのidをユーザの更新時に受け取りどのユーザを更新するかを判別する
	return ui.userOutputPort.OutputLoginResult(token)
}

func (ui *UserInteractor) GetAuthUrl() error {
	url := ui.userRepository.GenerateAuthUrl()
	return ui.userOutputPort.OutputAuthUrl(url)
}

func (ui *UserInteractor) SignupDraft(code string) error {
	email, err := ui.userRepository.GetUserInfoWithAuthCode(code)
	if err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}

	// 先にemailのみで登録する(仮登録)
	user, err := model.NewUser("", email, 0, 0.0, 0.0)
	if err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}
	// 存在しない場合にerrが返ってくるため、nilであればすでに存在しているということ
	if err := ui.userRepository.Exist(user); err == nil {
		// すでに登録されている場合はログイン画面に遷移させる
		if err := ui.userOutputPort.OutputAlreadySignedup(); err != nil {
			return ui.userOutputPort.OutputUserHandleError(err)
		}
		return ui.userOutputPort.OutputUserHandleError(err)
	}
	if user, err = ui.userRepository.Create(user); err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}
	token, err := ui.userRepository.GenerateAccessToken(user.Id)
	if err != nil {
		return ui.userOutputPort.OutputUserHandleError(err)
	}

	// urlのクエリパラメータにidを付与してそのidをユーザの更新時に受け取りどのユーザを更新するかを判別する
	return ui.userOutputPort.OutputSignupWithAuth(token)
}
