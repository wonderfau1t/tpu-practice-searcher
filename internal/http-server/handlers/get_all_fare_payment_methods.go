package handlers

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type GetAllFarePaymentMethodController interface {
	GetAllFarePaymentMethods() ([]db_models.FarePayment, error)
}

type GetAllFarePaymentMethodsResult struct {
	PaymentMethods []db_models.FarePayment `json:"paymentMethods"`
}

func GetAllFarePaymentMethods(log *slog.Logger, db GetAllFarePaymentMethodController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.GetAllFarePaymentMethods"
		log := log.With(slog.String("fn", fn))

		methods, err := db.GetAllFarePaymentMethods()
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		result := GetAllFarePaymentMethodsResult{PaymentMethods: methods}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
