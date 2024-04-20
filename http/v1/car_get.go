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

type filterAdder func(filter *models.Filter, name, value string) error

// query filter fields name
const (
	regNumberFieldName = "reg_num"
	markFieldName = "mark"
	modelFieldName = "model"
	yearFieldName = "year"
	ownerNameFieldName = "owner_name"
	ownerSurnameFieldName = "owner_surname"
	ownerPatronymicFieldName = "owner_patronymic"
)

// @summary Получить данные машины
// @tags Car
// @description Получение данных машины по ее идентификатору
// @id Car_get_one
// @produce json
// @Param carId path string true "Идентификатор машины"
// @Router /api/car/{carId} [get]
// @Success 200 {object} httpmodels.CarGetOneResponse
// @Failure 400
//
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

// @summary Получить данные машин с фильтром и пагинацией
// @tags Car
// @description Получение данных машин с фильтром и пагинацией
// @description
// @description Фильтр включает в себя 2 необязательных и 1 обязательный параметр
// @description col_name=UnionCondition:Operator:Value
// @description Где:
// @description UnionCondition(необязательный) - условия включения с другими фильтрами or/and (по умолчанию and)
// @description Operator(необязательный) - логический оператор (eq,neq,gt,get,lt,let,like) (по умолчанию eq)
// @description Value(обязательный) - само значение для фильтра
// @id Car_get_all
// @produce json
// @Param reg_nums query []string false "Фильтр для поля регистрационного номера" example(like:X123XX150) collectionFormat(multi)
// @Param marks query []string false "Фильтр для поля марки" example(or:like:Lada) collectionFormat(multi)
// @Param models query []string false "Фильтр для поля модели" collectionFormat(multi)
// @Param year query []string false "Фильтр для поля года" example(and:gt:2001) collectionFormat(multi)
// @Param owner_name query []string false "Фильтр для поля имени владельца" collectionFormat(multi)
// @Param owner_surname query []string false "Фильтр для поля фамилии владельца" collectionFormat(multi)
// @Param owner_patronymic query []string false "Фильтр для поля отчества владельца" collectionFormat(multi)
// @Param limit query integer false "Количество записей на странице" minimum(1)
// @Param offset query integer false "Количество пропущенных записей"
// @Router /api/cars [get]
// @Success 200 {object} httpmodels.CarGetAllResponse
// @Failure 400
//
func CarGetAll(logger *slog.Logger, carGetter carAllGetter, fAdder filterAdder) http.HandlerFunc {
	log := logger.With(slog.String("handler", "get_all_cars"))
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("attempt to get all cars")
		var (
			pagOption models.PaginationOption
			filter models.Filter
			err error
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
		queries := r.URL.Query()
		regNumbers := queries[regNumberFieldName]
		if len(regNumbers) != 0 {
			for i := 0; i < len(regNumbers); i++ {
				// make error terminating
				err = fAdder(&filter, regNumberFieldName, regNumbers[i])
				if err != nil {
					log.Warn("wrong filter",
					slog.String("field", regNumberFieldName),
					slog.String("value", regNumbers[i]),
					slog.String("error", err.Error()))
				}
			}
		}

		marks := queries[markFieldName]
		if len(marks) != 0 {
			for i := 0; i < len(marks); i++ {
				err = fAdder(&filter, markFieldName, marks[i])
				if err != nil {
					log.Warn("wrong filter",
					slog.String("field", markFieldName),
					slog.String("value", marks[i]),
					slog.String("error", err.Error()))
				}
			}
		}

		models := queries[modelFieldName]
		if len(models) != 0 {
			for i := 0; i < len(models); i++ {
				err = fAdder(&filter, modelFieldName, models[i])
				if err != nil {
					log.Warn("wrong filter",
					slog.String("field", modelFieldName),
					slog.String("value", models[i]),
					slog.String("error", err.Error()))
				}
			}
		}

		year := queries[yearFieldName]
		if len(year) != 0 {
			for i := 0; i < len(year); i++ {
				err = fAdder(&filter, yearFieldName, year[i])
				if err != nil {
					log.Warn("wrong filter",
					slog.String("field", yearFieldName),
					slog.String("value", year[i]),
					slog.String("error", err.Error()))
				}
			}
		}

		ownerNames := queries[ownerNameFieldName]
		if len(ownerNames) != 0 {
			for i := 0; i < len(ownerNames); i++ {
				err = fAdder(&filter, ownerNameFieldName, ownerNames[i])
				if err != nil {
					log.Warn("wrong filter",
					slog.String("field", ownerNameFieldName),
					slog.String("value", ownerNames[i]),
					slog.String("error", err.Error()))
				}
			}
		}

		ownerSurnames := queries[ownerSurnameFieldName]
		if len(ownerSurnames) != 0 {
			for i := 0; i < len(ownerSurnames); i++ {
				err = fAdder(&filter, ownerSurnameFieldName, ownerSurnames[i])
				if err != nil {
					log.Warn("wrong filter",
					slog.String("field", ownerSurnameFieldName),
					slog.String("value", ownerSurnames[i]),
					slog.String("error", err.Error()))
				}
			}
		}

		ownerPatronymics := queries[ownerPatronymicFieldName]
		if len(ownerPatronymics) != 0 {
			for i := 0; i < len(ownerPatronymics); i++ {
				err = fAdder(&filter, ownerPatronymicFieldName, ownerPatronymics[i])
				if err != nil {
					log.Warn("wrong filter",
					slog.String("field", ownerPatronymicFieldName),
					slog.String("value", ownerPatronymics[i]),
					slog.String("error", err.Error()))
				}
			}
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