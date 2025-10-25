package api

import (
	"github.com/aidenfine/recruiter-rater/internal/db/repositories"
	"github.com/aidenfine/recruiter-rater/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type V1Router struct {
	db *sqlx.DB
}

func NewV1Router(db *sqlx.DB) *V1Router {
	return &V1Router{
		db: db,
	}
}

func (router *V1Router) V1Router() chi.Router {
	r := chi.NewRouter()

	v1Repository := repositories.NewV1Repository(router.db)
	handler := handler.NewApiHandler(v1Repository)
	r.Get("/recruiter", handler.GetRecruiter())
	r.Get("/recruiter/{id}", handler.GetRecruiterById())
	r.Post("/add-recruiter", handler.AddNewRecruiter())

	r.Post("/reviews", handler.AddNewReview())
	r.Get("/reviews", handler.GetRecruiterReviews())
	r.Get("/reviews/recent", handler.GetMostRecentReviews())
	return r

}
