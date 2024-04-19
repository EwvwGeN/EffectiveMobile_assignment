package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/config"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/httpmodels"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/validator"
	"github.com/gorilla/mux"
)

type carEditor interface {
	EditCar(context.Context, string, models.CarForPatch) error
}

func CarEdit(logger *slog.Logger, validCfg config.ValidatorConfig, cEdditor carEditor) http.HandlerFunc {
	log := logger.With(slog.String("handler", "edit_car"))
	currentYear := uint16(time.Now().Year())
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("attempt to edit a car")
		carId, ok := mux.Vars(r)["carId"]
		if !ok || carId == "" {
			log.Warn("empty car id")
			http.Error(w, "error while deleting car: empty car id", http.StatusBadRequest)
			return
		}
		log.Debug("got car id", slog.String("car_id", carId))
		req := &httpmodels.CarEditRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			http.Error(w, "error while decoding request", http.StatusBadRequest)
			return
		}
		log.Debug("got data from request", slog.Any("request_body", req))
		//TODO: add all_field_is_nil case handling
		//TODO: check from json in a for statement(?)
		if req.CarNewData.RegisterNumber != nil && !validator.ValideteByRegex(*req.CarNewData.RegisterNumber, validCfg.RegisterNumberRegex) {
			log.Info("validate error: incorrect new car register number", slog.String("registe_number", *req.CarNewData.RegisterNumber))
			http.Error(w, "not valid car register number", http.StatusBadRequest)
			return
		}
		if req.CarNewData.Mark != nil && !validator.ValideteByRegex(*req.CarNewData.Mark, validCfg.MarkRegex) {
			log.Info("validate error: incorrect new car mark", slog.String("mark", *req.CarNewData.Mark))
			http.Error(w, "not valid car mark", http.StatusBadRequest)
			return
		}
		if req.CarNewData.Model != nil && !validator.ValideteByRegex(*req.CarNewData.Model, validCfg.ModelRegex) {
			log.Info("validate error: incorrect new car model", slog.String("model", *req.CarNewData.Model))
			http.Error(w, "not valid car model", http.StatusBadRequest)
			return
		}
		if req.CarNewData.Year != nil && (*req.CarNewData.Year < 1900 || *req.CarNewData.Year > currentYear) {
			log.Info("validate error: incorrect new car year", slog.Any("year", *req.CarNewData.Year))
			http.Error(w, "not valid car year", http.StatusBadRequest)
			return
		}
		if req.CarNewData.Owner != nil && req.CarNewData.Owner.Name != nil && !validator.ValideteByRegex(*req.CarNewData.Owner.Name, validCfg.OwnerNameRegex) {
			log.Info("validate error: incorrect new owner name", slog.String("owner_name", *req.CarNewData.Owner.Name))
			http.Error(w, "not valid owner name", http.StatusBadRequest)
			return
		}
		if req.CarNewData.Owner != nil && req.CarNewData.Owner.Surname != nil && !validator.ValideteByRegex(*req.CarNewData.Owner.Surname, validCfg.OwnerSurnameRegex) {
			log.Info("validate error: incorrect new owner surname", slog.String("owner_surname", *req.CarNewData.Owner.Surname))
			http.Error(w, "not valid owner surname", http.StatusBadRequest)
			return
		}
		if req.CarNewData.Owner != nil && req.CarNewData.Owner.Patronymic != nil && !validator.ValideteByRegex(*req.CarNewData.Owner.Patronymic, validCfg.OwnerPatronymicRegex) {
			log.Info("validate error: incorrect new owner patronymic", slog.String("owner_patronymic", *req.CarNewData.Owner.Patronymic))
			http.Error(w, "not valid owner patronymic", http.StatusBadRequest)
			return
		}
		err := cEdditor.EditCar(context.Background(), carId, req.CarNewData)
		if err != nil {
			log.Warn("failed to edit the car",
			slog.String("car_id", carId),
			slog.Any("car_new_data", req.CarNewData),
			slog.String("error", err.Error()))
			http.Error(w, "error while editing the car", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}