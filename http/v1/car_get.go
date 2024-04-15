package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/httpmodels"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"
	"github.com/gorilla/mux"
)

type carOneGetter interface {
	GetOneCar(ctx context.Context, carId string) (models.Car, error)
}

type carAllGetter interface {
	GetAllCars(context.Context, models.PaginationOption, models.Filter) ([]models.Car, error)
}

type filterAdder interface {
	AddFilter(filter *models.Filter, name, value string) error
}

// query filter fields name
const (
	regNumberFieldName = "regNum"
	markFieldName = "mark"
	modelFieldName = "model"
	yearFieldName = "year"
	ownerFieldName = "owner"
)

func CarGetOne(logger *slog.Logger, carGetter carOneGetter) http.HandlerFunc {
	log := logger.With(slog.String("handler", "get_one_car"))
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("attempt to get one car")
		carId, ok := mux.Vars(r)["carId"]
		if !ok || carId == "" {
			log.Warn("failed to get car id")
			http.Error(w, "error while getting car: empty car id", http.StatusBadRequest)
			return
		}
		log.Debug("got car id", slog.String("car_id", carId))
		car, err := carGetter.GetOneCar(context.Background(), carId)
		if err != nil {
			log.Error("failed to get car", slog.String("error", err.Error()))
			http.Error(w, "error while getting car", http.StatusBadRequest)
			return
		}
		log.Debug("got car", slog.Any("car", car))
		res := &httpmodels.CarGetOneResponse{
			Car: car,
		}
		resData, err := json.Marshal(res)
		if err != nil {
			log.Error("cant encode response", slog.Any("response", res), slog.String("error", err.Error()))
			http.Error(w, "error while getting car", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resData)
	}
}

func CarGetAll(logger *slog.Logger, carGetter carAllGetter, fAdder filterAdder) http.HandlerFunc {
	log := logger.With(slog.String("handler", "get_all_cars"))
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("attempt to get all cars")
		var (
			pagOption models.PaginationOption
			filter models.Filter
		)
		limit := mux.Vars(r)["limit"]
		offset := mux.Vars(r)["offset"]
		log.Debug("got limit and offset", slog.String("limit", limit), slog.String("offset", offset))
		if limit != "" && offset != "" {
			lInt, _ := strconv.Atoi(limit)
			pagOption.Limit = lInt
			oint, _ := strconv.Atoi(offset)
			pagOption.Limit = oint
		}
		regNumber := r.URL.Query().Get(regNumberFieldName)
		err := fAdder.AddFilter(&filter, regNumberFieldName, regNumber)
		// make error terminating
		if err != nil {
			log.Warn("wrong filter",
			slog.String("field", regNumberFieldName),
			slog.String("value", regNumber),
			slog.String("error", err.Error()))
		}
		mark := r.URL.Query().Get(markFieldName)
		err = fAdder.AddFilter(&filter, markFieldName, mark)
		if err != nil {
			log.Warn("wrong filter",
			slog.String("field", markFieldName),
			slog.String("value", mark),
			slog.String("error", err.Error()))
		}
		model := r.URL.Query().Get(modelFieldName)
		err = fAdder.AddFilter(&filter, modelFieldName, model)
		if err != nil {
			log.Warn("wrong filter",
			slog.String("field", modelFieldName),
			slog.String("value", model),
			slog.String("error", err.Error()))
		}
		year := r.URL.Query().Get(yearFieldName)
		err = fAdder.AddFilter(&filter, yearFieldName, year)
		if err != nil {
			log.Warn("wrong filter",
			slog.String("field", yearFieldName),
			slog.String("value", year),
			slog.String("error", err.Error()))
		}
		owner := r.URL.Query().Get(ownerFieldName)
		err = fAdder.AddFilter(&filter, ownerFieldName, owner)
		if err != nil {
			log.Warn("wrong filter",
			slog.String("field", ownerFieldName),
			slog.String("value", owner),
			slog.String("error", err.Error()))
		}
		cars, err := carGetter.GetAllCars(context.Background(), pagOption, filter)
		if err != nil {
			log.Error("failed to get cars", slog.String("error", err.Error()))
			http.Error(w, "error while getting cars", http.StatusBadRequest)
			return
		}
		log.Debug("got cars", slog.Any("cars", cars))
		res := &httpmodels.CarGetAllResponse{
			Cars: cars,
		}
		resData, err := json.Marshal(res)
		if err != nil {
			log.Error("cant encode response", slog.Any("response", res), slog.String("error", err.Error()))
			http.Error(w, "error while getting cars", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resData)
	}
}