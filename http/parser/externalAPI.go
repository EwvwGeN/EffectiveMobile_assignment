package parser

import (
	"encoding/json"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"
)

func ParseFromExternalApi(carMap map[string]interface{}) (models.Car, error) {
	carJson, err := json.Marshal(carMap)
    if err != nil {
        return models.Car{}, err
    }
	var newCar models.Car
	err = json.Unmarshal(carJson, &newCar)
	if err != nil {
		return models.Car{}, err
	}
	return newCar, nil
}