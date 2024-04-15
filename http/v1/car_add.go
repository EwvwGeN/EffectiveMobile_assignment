package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/httpmodels"
)

type carAdder interface {
	AddCar(context.Context, []string) error
}

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
		err := cAdder.AddCar(context.Background(), req.RegisterNumbers)
		if err != nil {
			log.Error("failed to add car", slog.Any("register_numbers", req.RegisterNumbers), slog.String("error", err.Error()))
			http.Error(w, "error while adding car", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}