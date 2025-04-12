package handlers

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type GetAllAccommodationPaymentMethodsController interface {
	GetAllAccommodationPaymentMethods() ([]db_models.PaymentForAccommodation, error)
}

type GetAllAccommodationPaymentMethodsResult struct {
	PaymentMethods []db_models.PaymentForAccommodation `json:"paymentMethods"`
}

func GetAllAccommodationPaymentMethods(log *slog.Logger, db GetAllAccommodationPaymentMethodsController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.GetAllAccommodationPaymentMethods"
		log := log.With(slog.String("fn", fn))

		methods, err := db.GetAllAccommodationPaymentMethods()
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		result := GetAllAccommodationPaymentMethodsResult{PaymentMethods: methods}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
