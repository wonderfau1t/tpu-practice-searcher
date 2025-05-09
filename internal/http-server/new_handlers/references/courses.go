package references

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type ReferenceRepository interface {
	GetDepartments() ([]db_models.Department, error)
}

type Response struct {
	Departments []DepartmentDTO `json:"departments"`
}

type DepartmentDTO struct {
	ID      uint        `json:"ID"`
	Name    string      `json:"name"`
	Courses []CourseDTO `json:"courses"`
}

type CourseDTO struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}

func Courses(log *slog.Logger, repo ReferenceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Courses"
		log := log.With(slog.String("fn", fn))
		_ = log

		departments, err := repo.GetDepartments()
		if err != nil {
			return
		}

		dtos := make([]DepartmentDTO, len(departments))
		for i, department := range departments {
			dtos[i] = DepartmentDTO{
				ID:   department.ID,
				Name: department.Name,
			}
			for _, course := range department.Courses {
				courseDTO := CourseDTO{
					ID:   course.ID,
					Name: course.Name,
				}
				dtos[i].Courses = append(dtos[i].Courses, courseDTO)
			}
		}

		response := Response{Departments: dtos}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(response))
	}
}
