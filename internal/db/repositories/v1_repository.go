package repositories

import (
	"database/sql"
	"errors"

	"github.com/aidenfine/recruiter-rater/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type V1RepositoryInterface interface {
	GetRecruiterByLinkedinUsername(username string) (*models.RecruiterResponse, error)
	GetRecruiterById(id string) (*models.RecruiterResponse, error)
	AddNewRecruiter(recruiter *models.RecruiterPayload) error
	SearchRecruiters(searchTerm string) ([]models.RecruiterResponse, error)
	AddNewReview(payload *models.ReviewPayload) error
	GetAllRecruiterReviews(recruiterId uuid.UUID, limit int) (*[]models.ReviewResponse, error)
}

type V1Repository struct {
	db *sqlx.DB
}

func NewV1Repository(db *sqlx.DB) *V1Repository {
	return &V1Repository{
		db: db,
	}
}

func (vr *V1Repository) GetRecruiterByLinkedinUsername(username string) (*models.RecruiterResponse, error) {
	query := `
	SELECT id, name, first_name, last_name, linkedin_username, job_title, image_url, rating, current_company, created_at, verified
	FROM recruiters WHERE linkedin_username = $1 AND verified = true
	`

	var recruiter models.RecruiterResponse
	err := vr.db.Get(&recruiter, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not_found")
		}
		return nil, err
	}

	return &recruiter, nil
}

func (vr *V1Repository) GetRecruiterById(id string) (*models.RecruiterResponse, error) {
	query := `
	SELECT id, name, first_name, last_name, linkedin_username, job_title, image_url, rating, current_company, created_at, verified
	FROM recruiters
	WHERE id = $1
	AND verified = true
	`
	var recruiter models.RecruiterResponse
	err := vr.db.Get(&recruiter, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not_found")
		}
		return nil, err
	}

	return &recruiter, nil
}

func (vr *V1Repository) AddNewRecruiter(recruiter *models.RecruiterPayload) error {
	query := `INSERT INTO recruiters (name, first_name, last_name, linkedin_username, job_title, current_company) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := vr.db.Exec(query, recruiter.Name, recruiter.FirstName, recruiter.LastName, recruiter.Linkedin, recruiter.JobTitle, recruiter.CurrentCompany)
	if err != nil {
		return err
	}
	return nil
}

func (vr *V1Repository) SearchRecruiters(searchTerm string) ([]models.RecruiterResponse, error) {
	query := `
	SELECT id, name, first_name, last_name, linkedin_username, job_title, image_url, rating, current_company, created_at, verified
	FROM recruiters 
	WHERE verified = true 
	AND (
		linkedin_username ILIKE $1 
		OR name ILIKE $1 
		OR first_name ILIKE $1 
		OR last_name ILIKE $1
		OR current_company ILIKE $1
	)
	AND verified = true
	ORDER BY rating DESC NULLS LAST
	LIMIT 50
	`

	var recruiters []models.RecruiterResponse
	searchPattern := "%" + searchTerm + "%"
	err := vr.db.Select(&recruiters, query, searchPattern)
	if err != nil {
		return nil, err
	}

	if len(recruiters) == 0 {
		return nil, errors.New("not_found")
	}

	return recruiters, nil
}

func (vr *V1Repository) AddNewReview(payload *models.ReviewPayload) error {
	query := `
	INSERT INTO reviews (recruiter_id, rating, description) VALUES ($1, $2, $3)
	`
	_, err := vr.db.Exec(query, payload.RecruiterId, payload.Rating, payload.Description)
	if err != nil {
		return err
	}
	return nil
}
func (vr *V1Repository) GetAllRecruiterReviews(recruiterId uuid.UUID, limit int) (*[]models.ReviewResponse, error) {
	query := `
	SELECT id, recruiter_id, created_at, rating, description, thumbs_down_count, thumbs_up_count
	FROM reviews
	WHERE recruiter_id = $1
	ORDER BY created_at DESC
	LIMIT $2
	`

	items := []models.ReviewResponse{}
	err := vr.db.Select(&items, query, recruiterId, limit)
	if err != nil {
		return nil, err
	}
	return &items, nil
}
