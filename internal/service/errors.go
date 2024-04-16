package service

import "errors"

var (
	ErrGetCarInfo = errors.New("car with this register number didnt find")
)