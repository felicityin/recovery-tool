package code

import "fmt"

type I18nError struct {
	Code string
	Msg  string
}

func NewI18nError(code string, msg string) error {
	return &I18nError{
		Code: code,
		Msg:  msg,
	}
}

func (err *I18nError) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Msg)
}
