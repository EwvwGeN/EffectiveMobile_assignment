package service

import (
	"context"
	"log/slog"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"
)

type carService struct {
	log     *slog.Logger
	carRepo carRepo
	carInfoGetter carInfoGetter
}

type carRepo interface {
	SaveCars(context.Context, []models.Car) (error)
	GetCarById(context.Context, string) (models.Car, error)
	GetCarsWithFilterAndPagination(context.Context, models.PaginationOption, models.Filter) ([]models.Car, error)
	UpdateCarById(context.Context, string, models.CarForPatch) (error)
	DeleteCarById(context.Context, string) (error)
}

type carInfoGetter func(context.Context, string) (models.Car, error)

func NewCarService(logger *slog.Logger, cRepo carRepo, cGetter carInfoGetter) *carService {
	return &carService{
		log: logger.With(slog.String("service", "car")),
		carRepo:  cRepo,
		carInfoGetter: cGetter,
	}
}
func (cs *carService) AddCar(ctx context.Context, regNumbers []string) error {
	var (
		carList []models.Car
		car models.Car
		err error
	)
	cs.log.Info("attempt to add a car")
	cs.log.Debug("got cars register numbers", slog.Any("register_numbers", regNumbers))
	for _, regNumber := range regNumbers {
		car, err = cs.carInfoGetter(ctx, regNumber)
		if err != nil {
			cs.log.Warn("failed to get car info", slog.String("register_number", regNumber))
			break
		}
		carList = append(carList, car)
	}
	if err != nil  {
		return ErrGetCarInfo
	}
	cs.log.Debug("got cars info", slog.Any("cars_info", carList))
	err = cs.carRepo.SaveCars(ctx, carList)
	if err != nil  {
		cs.log.Error("failed to save cars", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (cs *carService) GetOneCar(ctx context.Context, carId string) (models.Car, error) {
	cs.log.Info("attempt to get car by id")
	cs.log.Debug("got car id", slog.String("car_id", carId))
	car, err := cs.carRepo.GetCarById(ctx, carId)
	if err != nil {
		cs.log.Error("failed to get car by id", slog.String("car_id", carId), slog.String("error", err.Error()))
		return models.Car{}, err
	}
	return car, nil
}

func (cs *carService) GetAllCars(ctx context.Context, pOption models.PaginationOption, filter models.Filter) ([]models.Car, error) {
	cs.log.Info("attempt to get all cars with filter")
	cs.log.Debug("got filter and pagination options", slog.Any("pagination_option", pOption), slog.Any("filter", filter))
	carList, err := cs.carRepo.GetCarsWithFilterAndPagination(ctx, pOption, filter)
	if err != nil {
		cs.log.Error("failed to get cars",
		slog.Any("pagination_option", pOption),
		slog.Any("filter", filter),
		slog.String("error", err.Error()))
		return nil, err
	}
	return carList, err
}
func (cs *carService) EditCar(ctx context.Context, carId string, newData models.CarForPatch) error {
	cs.log.Info("attempt to edit car by id")
	cs.log.Debug("got car id", slog.String("car_id", carId))
	err := cs.carRepo.UpdateCarById(ctx, carId, newData)
	if err != nil {
		cs.log.Error("failed to edit car by id", slog.String("car_id", carId), slog.String("error", err.Error()))
		return err
	}
	return nil
}
func (cs *carService) DeleteCar(ctx context.Context, carId string) error {
	cs.log.Info("attempt to edit car by id")
	cs.log.Debug("got car id", slog.String("car_id", carId))
	err := cs.carRepo.DeleteCarById(ctx, carId)
	if err != nil {
		cs.log.Error("failed to delete car by id", slog.String("car_id", carId), slog.String("error", err.Error()))
		return err
	}
	return nil
}
