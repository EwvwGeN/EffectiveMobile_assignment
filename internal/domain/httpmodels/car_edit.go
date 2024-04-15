package httpmodels

import "github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"

type CarEditRequest struct {
	CarNewData models.CarForPatch `json:"carNewData"`
}