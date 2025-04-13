package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"net/http"
	"os"
	"tpu-practice-searcher/internal/config"
	"tpu-practice-searcher/internal/http-server/handlers"
	"tpu-practice-searcher/internal/http-server/handlers/create_new_hr"
	"tpu-practice-searcher/internal/http-server/handlers/get_all_hrs_of_company"
	"tpu-practice-searcher/internal/http-server/handlers/get_all_vacancies_of_company"
	"tpu-practice-searcher/internal/http-server/handlers/get_company_details"
	"tpu-practice-searcher/internal/http-server/handlers/get_vacancy_details"
	"tpu-practice-searcher/internal/http-server/handlers/moderator/get_all_vacancies_by_courses"
	"tpu-practice-searcher/internal/http-server/handlers/moderator/get_vacancy_info"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_company_info"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_replies"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_vacancies"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_vacancies_of_company"
	"tpu-practice-searcher/internal/http-server/handlers/student/reply_to_vacancy"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/logger"
	"tpu-practice-searcher/internal/storage/postgresql"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("Config successfully loaded")

	db, err := postgresql.SetupStorage(cfg.Storage)
	if err != nil {
		log.Error("Failed to set storage")
		os.Exit(1)
	}
	log.Info("Storage successfully set")

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "Hello world!"})
	})

	// Справочная информация
	router.Route("/search", func(r chi.Router) {
		r.Get("/categories", handlers.GetAllCategories(log, db))
		r.Get("/formats", handlers.GetAllFormats(log, db))
		r.Get("/farePaymentMethods", handlers.GetAllFarePaymentMethods(log, db))
		r.Get("/accommodationPaymentMethods", handlers.GetAllAccommodationPaymentMethods(log, db))
		r.Get("/courses", handlers.GetAllCourses(log, db))
	})

	// Вакансии
	router.Get("/vacancies", get_vacancies.New(log, db))
	// Подробная информация о вакансии
	router.Get("/vacancies/{id}", get_vacancy_details.New(log, db))
	// Вакансии определенной компании
	router.Get("/companies/{id}/vacancies", get_vacancies_of_company.New(log, db))
	// Профиль компании
	router.Get("/companies/{id}/info", get_company_info.New(log, db))

	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Get("/auth", handlers.Auth(log, db))
		r.Get("/register", handlers.RegisterStudent(log, db))
		r.Post("/registerCompany", handlers.RegisterCompany(log, db))
		r.Post("/addVacancy", handlers.AddVacancy(log, db))
		r.Get("/company/vacancies", get_all_vacancies_of_company.New(log, db))
		r.Get("/company/hrs", get_all_hrs_of_company.New(log, db))
		r.Post("/company/createHr", create_new_hr.New(log, db))
		r.Get("/company/info", get_company_details.New(log, db))

		r.Post("/reply", reply_to_vacancy.New(log, db))
		r.Get("/replies", get_replies.New(log, db))

		r.Get("/school/vacancies", get_all_vacancies_by_courses.New(log, db))
		r.Get("/moderator/vacancies/{id}", get_vacancy_info.New(log, db))
	})

	http.ListenAndServe("0.0.0.0:8000", router)
}
