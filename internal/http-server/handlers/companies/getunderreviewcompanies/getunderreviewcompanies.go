package getunderreviewcompanies

import (
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
)

type Storage interface {
	GetUnderReviewCompanies() ([]db_models.Company, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.getunderreviewcompanies.New"
		log := log.With(slog.String("fn", fn))

		_ = log
	}
}
