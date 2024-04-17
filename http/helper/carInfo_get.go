package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"
)

type parser func (map[string]interface{}) (models.Car, error)

func GetCarInfoGetter(logger *slog.Logger, sourceUrl string, parseFunc parser) (func(context.Context, string) (models.Car, error), error) {
	parsedUrl, err := url.Parse(sourceUrl)
	if err != nil {
		return nil, err
	}
	log := logger.With("handler", "car_info_getter")
	return func(ctx context.Context, carRegisteNum string) (models.Car, error) {
		newUrl := *parsedUrl
		values := newUrl.Query()
		values.Set("regNum", carRegisteNum)
		newUrl.RawQuery = values.Encode()
		resp, err := http.Get(newUrl.String())
		if err != nil {
			log.Error("filed to execute get request", slog.String("error", err.Error()))
			return models.Car{}, err
		}
		var carMap map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&carMap)
		if err != nil {
			log.Error("failed to decode response body", slog.String("error", err.Error()))
			return models.Car{}, err
		}
		resp.Body.Close()
		if len(carMap) == 0 {
			return models.Car{}, fmt.Errorf("empty car info response")
		}
		log.Debug("got car map", slog.Any("car_map", carMap))
		car, err := parseFunc(carMap)
		if err != nil {
			return models.Car{}, err
		}
		return car, nil
	}, nil
}