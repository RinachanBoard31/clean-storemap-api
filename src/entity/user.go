package model

import (
	"errors"
	"regexp"
)

type User struct {
	Id     int
	Name   string
	Email  string
	Age    int     // xx代として表記する(60代以上は全て60とする)
	Sex    float32 // -1.0(男性)~1.0(女性)で表現する。中性、無回答は0となる。
	Gender float32 // -1.0(男性)~1.0(女性)で表現する。中性、無回答は0となる。
}

func emailValid(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("emailではありません")
	}
	return nil
}

func ageValid(age int) error {
	if age < 0 {
		return errors.New("年齢が0未満です。")
	}
	return nil
}

func ageFormat(age int) int {
	// ageValidで0未満はエラーとなるので0未満は扱わない。
	if age >= 60 {
		return 60
	}
	return (age / 10) * 10
}

func sexFormat(sex float32) float32 {
	if sex < -1.0 {
		return -1.0
	}
	if sex > 1.0 {
		return 1.0
	}
	return sex
}

func genderFormat(gender float32) float32 {
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
	ageValidError := ageValid(age)
	if err := errors.Join(emailValidError, ageValidError); err != nil {
		return nil, err
	}
	// userの作成
	user := &User{
		Name:   name,
		Email:  email,
		Age:    ageFormat(age),
		Sex:    sexFormat(sex),
		Gender: genderFormat(gender),
	}
	return user, nil
}
