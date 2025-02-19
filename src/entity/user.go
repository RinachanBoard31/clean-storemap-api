package model

import (
	"clean-storemap-api/src/entity/errors"

	"regexp"
)

// type UserId string

type User struct {
	Id     string // uuidを使用
	Name   string
	Email  string
	Age    int     // xx代として表記する(60代以上は全て60とする)
	Sex    float32 // -1.0(男性)~1.0(女性)で表現する。中性、無回答は0となる。
	Gender float32 // -1.0(男性)~1.0(女性)で表現する。中性、無回答は0となる。
}

type UserCredentials struct {
	Email string
}

type ChangeForUser map[string]interface{}

func emailValid(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.ErrInvalidEmail
	}
	return nil
}

func AgeValid(age int) error {
	if age < 0 {
		return errors.ErrAgeBelowZero
	}
	return nil
}

func AgeFormat(age int) int {
	// ageValidで0未満はエラーとなるので0未満は扱わない。
	if age >= 60 {
		return 60
	}
	return (age / 10) * 10
}

func SexFormat(sex float32) float32 {
	if sex < -1.0 {
		return -1.0
	}
	if sex > 1.0 {
		return 1.0
	}
	return sex
}

func GenderFormat(gender float32) float32 {
	if gender < -1.0 {
		return -1.0
	}
	if gender > 1.0 {
		return 1.0
	}
	return gender
}

func NewUser(name string, email string, age int, sex float32, gender float32) (*User, error) {
	// バリデーションのチェック
	emailValidError := emailValid(email)
	ageValidError := AgeValid(age)
	if emailValidError != nil && ageValidError != nil {
		return nil, errors.ErrInvalidEmailAndAge
	}
	// userの作成
	user := &User{
		Name:   name,
		Email:  email,
		Age:    AgeFormat(age),
		Sex:    SexFormat(sex),
		Gender: GenderFormat(gender),
	}
	return user, nil
}

func NewUserCredentials(email string) (*UserCredentials, error) {
	// バリデーションのチェック
	if err := emailValid(email); err != nil {
		return &UserCredentials{}, err
	}

	user := &UserCredentials{
		Email: email,
	}
	return user, nil
}
