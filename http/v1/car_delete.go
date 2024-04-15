package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type carDeleter interface {
	DeleteCar(context.Context, string) error
}

func CarDelete(logger *slog.Logger, cDeleter carDeleter) http.HandlerFunc {
	log := logger.With(slog.String("handler", "delete_car"))
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("attempt to delete car")
		carId, ok := mux.Vars(r)["carId"]
		if !ok || carId == "" {
			log.Warn("empty car id")
			http.Error(w, "error while deleting car: empty car id", http.StatusBadRequest)
			return
		}
		log.Debug("got car id", slog.String("car_id", carId))
		err := cDeleter.DeleteCar(context.Background(), carId)
		if err != nil {
			log.Warn("failed to delete the car", slog.String("car_id", carId), slog.String("error", err.Error()))
			http.Error(w, "error while deleting car", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}