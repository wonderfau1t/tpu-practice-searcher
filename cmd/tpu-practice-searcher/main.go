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
	"tpu-practice-searcher/internal/http-server/handlers/companies"
	"tpu-practice-searcher/internal/http-server/handlers/companies/updatecompanyinfo"
	"tpu-practice-searcher/internal/http-server/handlers/create_new_hr"
	"tpu-practice-searcher/internal/http-server/handlers/get_all_hrs_of_company"
	"tpu-practice-searcher/internal/http-server/handlers/get_all_vacancies_of_company"
	"tpu-practice-searcher/internal/http-server/handlers/get_company_details"
	"tpu-practice-searcher/internal/http-server/handlers/get_vacancy_details"
	"tpu-practice-searcher/internal/http-server/handlers/moderator/get_all_vacancies_by_courses"
	"tpu-practice-searcher/internal/http-server/handlers/moderator/get_vacancy_info"
	"tpu-practice-searcher/internal/http-server/handlers/replies/deletereply"
	"tpu-practice-searcher/internal/http-server/handlers/replies/makereply"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_company_info"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_replies"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_vacancies"
	"tpu-practice-searcher/internal/http-server/handlers/student/get_vacancies_of_company"
	"tpu-practice-searcher/internal/http-server/handlers/student/reply_to_vacancy"
	"tpu-practice-searcher/internal/http-server/handlers/vacancies"
	"tpu-practice-searcher/internal/http-server/handlers/vacancies/createvacancywithoutcompany"
	"tpu-practice-searcher/internal/http-server/handlers/vacancies/filter"
	"tpu-practice-searcher/internal/http-server/handlers/vacancies/hide"
	"tpu-practice-searcher/internal/http-server/handlers/vacancies/search"
	"tpu-practice-searcher/internal/http-server/middlewares"
	companiesModule "tpu-practice-searcher/internal/http-server/new_handlers/companies"
	"tpu-practice-searcher/internal/http-server/new_handlers/references"
	vacanciesModule "tpu-practice-searcher/internal/http-server/new_handlers/vacancies"
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
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/vacancies", vacanciesModule.List(log, db))
		r.Put("/vacancies/{id}", vacanciesModule.Update(log, db))

		r.Put("/companies/update", companiesModule.Update(log, db))
		r.Get("/companies/requests", companiesModule.ListRequests(log, db))
		r.Patch("/companies/apply", companiesModule.Apply(log, db))
		r.Patch("/companies/reject", companiesModule.Reject(log, db))
	})

	router.Get("/references/courses", references.Courses(log, db))

	// Список вакансий
	//router.Get("/vacancies", vacanciesModule.List(log, db))
	//// Обновить информацию о вакансии
	router.Put("/vacancies/{id}", vacanciesModule.Update(log, db))

	//router.Patch("/company", handlers.)
	router.Route("/backend", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"message": "Hello world!"})
		})

		r.Route("/search", func(r1 chi.Router) {
			r1.Get("/formats", handlers.GetAllFormats(log, db))
			r1.Get("/farePaymentMethods", handlers.GetAllFarePaymentMethods(log, db))
			r1.Get("/accommodationPaymentMethods", handlers.GetAllAccommodationPaymentMethods(log, db))
			r1.Get("/courses", handlers.GetAllCourses(log, db))
		})

		r.Route("/api/v1", func(r1 chi.Router) {
			r1.Use(middlewares.AuthMiddleware)
			r1.Get("/vacancies", vacancies.GetVacancies(log, db))
			r1.Get("/vacancies/{id}", vacancies.GetVacancyDetails(log, db))
			r1.Patch("/vacancies/{id}", hide.New(log, db))
			r1.Get("/vacancies/filter", filter.New(log, db))
			r1.Get("/vacancies/search", search.New(log, db))

			r1.Post("/vacancies/{vacancyID}/replies", makereply.New(log, db))
			r1.Delete("/vacancies/{vacancyID}/replies", deletereply.New(log, db))

			r1.Get("/companies/{id}", companies.New(log, db))
			r1.Put("/companies/updateInfo", updatecompanyinfo.New(log, db))

		})

		r.Get("/vacancies", get_vacancies.New(log, db))
		// Подробная информация о вакансии
		r.Get("/vacancies/{id}", get_vacancy_details.New(log, db))
		// Вакансии определенной компании
		r.Get("/companies/{id}/vacancies", get_vacancies_of_company.New(log, db))
		// Профиль компании
		r.Get("/companies/{id}/info", get_company_info.New(log, db))

		r.Group(func(r1 chi.Router) {
			r1.Use(middlewares.AuthMiddleware)
			r1.Get("/auth", handlers.Auth(log, db))
			r1.Get("/register", handlers.RegisterStudent(log, db))
			r1.Post("/registerCompany", handlers.RegisterCompany(log, db))
			r1.Post("/addVacancy", handlers.AddVacancy(log, db))
			r1.Get("/company/vacancies", get_all_vacancies_of_company.New(log, db))
			r1.Get("/company/hrs", get_all_hrs_of_company.New(log, db))
			r1.Post("/company/createHr", create_new_hr.New(log, db))
			r1.Get("/company/info", get_company_details.New(log, db))

			r1.Post("/reply", reply_to_vacancy.New(log, db))
			r1.Get("/replies", get_replies.New(log, db))

			r1.Get("/school/vacancies", get_all_vacancies_by_courses.New(log, db))
			r1.Get("/moderator/vacancies/{id}", get_vacancy_info.New(log, db))

			r1.Post("/createVacancyWithoutCompany", createvacancywithoutcompany.New(log, db))
		})
	})
	// Справочная информация
	//router.Route("/search", func(r chi.Router) {
	//	r.Get("/categories", handlers.GetAllCategories(log, db))
	//	r.Get("/formats", handlers.GetAllFormats(log, db))
	//	r.Get("/farePaymentMethods", handlers.GetAllFarePaymentMethods(log, db))
	//	r.Get("/accommodationPaymentMethods", handlers.GetAllAccommodationPaymentMethods(log, db))
	//	r.Get("/courses", handlers.GetAllCourses(log, db))
	//})

	// Финальные эндпоинты
	//router.Route("/api/v1", func(r chi.Router) {
	//	r.Use(middlewares.AuthMiddleware)
	//	r.Get("/vacancies", vacancies.GetVacancies(log, db))
	//	r.Get("/vacancies/{id}", vacancies.GetVacancyDetails(log, db))
	//	r.Get("/vacancies/filter", filter.New(log, db))
	//	r.Get("/vacancies/search", search.New(log, db))
	//
	//	r.Post("/vacancies/{vacancyID}/replies", makereply.New(log, db))
	//	r.Delete("/vacancies/{vacancyID}/replies", deletereply.New(log, db))
	//
	//	r.Get("/companies/{id}", companies.New(log, db))
	//
	//})
	// Вакансии
	//router.Get("/vacancies", get_vacancies.New(log, db))
	//// Подробная информация о вакансии
	//router.Get("/vacancies/{id}", get_vacancy_details.New(log, db))
	//// Вакансии определенной компании
	//router.Get("/companies/{id}/vacancies", get_vacancies_of_company.New(log, db))
	//// Профиль компании
	//router.Get("/companies/{id}/info", get_company_info.New(log, db))
	////
	//router.Group(func(r chi.Router) {
	//	r.Use(middlewares.AuthMiddleware)
	//	r.Get("/auth", handlers.Auth(log, db))
	//	r.Get("/register", handlers.RegisterStudent(log, db))
	//	r.Post("/registerCompany", handlers.RegisterCompany(log, db))
	//	r.Post("/addVacancy", handlers.AddVacancy(log, db))
	//	r.Get("/company/vacancies", get_all_vacancies_of_company.New(log, db))
	//	r.Get("/company/hrs", get_all_hrs_of_company.New(log, db))
	//	r.Post("/company/createHr", create_new_hr.New(log, db))
	//	r.Get("/company/info", get_company_details.New(log, db))
	//
	//	r.Post("/reply", reply_to_vacancy.New(log, db))
	//	r.Get("/replies", get_replies.New(log, db))
	//
	//	r.Get("/school/vacancies", get_all_vacancies_by_courses.New(log, db))
	//	r.Get("/moderator/vacancies/{id}", get_vacancy_info.New(log, db))
	//})

	http.ListenAndServe("0.0.0.0:8000", router)
}
