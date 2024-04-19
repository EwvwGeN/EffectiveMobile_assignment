package service

import "errors"

var (
	ErrGetCarInfo = errors.New("car with this register number didnt find")
	ErrAddCar = errors.New("failed to save cars")
	ErrGetCar = errors.New("failed to get car")
	ErrEditCar = errors.New("failed to edit car")
	ErrDeleteCar = errors.New("failed to delete car")
)