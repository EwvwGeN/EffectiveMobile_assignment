package httpmodels

import "github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"

type CarGetOneResponse struct {
	Car models.Car `json:"car"`
}

type CarGetAllResponse struct {
	Cars []models.Car `json:"cars"`
}