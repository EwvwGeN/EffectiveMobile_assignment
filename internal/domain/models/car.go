package models

type Car struct {
	Id             int    `json:"carId"`
	RegisterNumber string `json:"regNum"`
	Mark           string `json:"mark"`
	Model          string `json:"model"`
	Year           uint16 `json:"year"`
	// we will store full name in record of car, of course better to create
	// individual table for owners and store there only id of owner object
	Owner Owner `json:"owner"`
}

type CarForPatch struct {
	Id             *string        `json:"carId"`
	RegisterNumber *string        `json:"regNum"`
	Mark           *string        `json:"mark"`
	Model          *string        `json:"model"`
	Year           *uint16        `json:"year"`
	Owner          *OwnerForPatch `json:"owner"`
}