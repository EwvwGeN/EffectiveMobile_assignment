package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/httpmodels"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/service"
)

type carAdder interface {
	AddCar(context.Context, []string) error
}

// @summary Добавить машину
// @tags Car
// @description Добавление машины по ее регистрационному номеру
// @id Car_add
// @accept json
// @produce plain
// @Param regNums body []string true "Регистрационные номера машины" SchemaExample({\n\r "regNums": ["string"]\n\r}) 
// @Router /api/cars/add [post]
// @Success 201
// @Failure 400
//
func CarAdd(logger *slog.Logger, cAdder carAdder) http.HandlerFunc {
	log := logger.With(slog.String("handler", "add_cars"))
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("attempt to add a cars")
		req := &httpmodels.CarAddRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			http.Error(w, "error while decoding request", http.StatusBadRequest)
			return
		}
		log.Debug("got data from request", slog.Any("request_body", req))
		if len(req.RegisterNumbers) == 0 {
			log.Warn("empty register numbers")
			http.Error(w, "error while adding car: empty register numbers", http.StatusBadRequest)
			return
		}
		err := cAdder.AddCar(context.Background(), req.RegisterNumbers)
		if err != nil {
			if errors.Is(err, service.ErrGetCarInfo) {
				log.Warn("failed to get cars info", slog.Any("register_numbers", req.RegisterNumbers), slog.String("error", err.Error()))
				http.Error(w, "error while adding car", http.StatusBadRequest)
				return
			}
			log.Error("failed to add car", slog.Any("register_numbers", req.RegisterNumbers), slog.String("error", err.Error()))
			http.Error(w, "error while adding car", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}