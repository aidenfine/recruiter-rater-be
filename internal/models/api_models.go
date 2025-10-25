package models

import (
	"time"

	"github.com/google/uuid"
)

type GetRecruiterQueryParams struct {
	Search string `query:"search" required:"true"`
	Url    bool   `query:"url" required:"true"`
}
type GetReviewsQueryParams struct {
	RecruiterId string `query:"recruiterId" required:"true"`
	Limit       int    `query:"limit" required:"true"`
}

type RecruiterResponse struct {
	Id               uuid.UUID `db:"id" json:"id"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	Name             string    `db:"name" json:"name"`
	FirstName        string    `db:"first_name" json:"firstName"`
	LastName         string    `db:"last_name" json:"lastName"`
	LinkedInUsername string    `db:"linkedin_username" json:"linkedinUsername"`
	JobTitle         *string   `db:"job_title" json:"jobTitle"`
	ImageUrl         *string   `db:"image_url" json:"imageUrl"`
	Rating           float32   `db:"rating" json:"rating"`
	CurrentCompany   *string   `db:"current_company" json:"currentCompany"`
	Verified         bool      `db:"verified" json:"verified"`
	RatingSum        int       `db:"rating_sum" json:"rating_sum"`
	NumberOfRatings  int       `db:"number_of_ratings" json:"number_of_ratings"`
}
type RecruiterPayload struct {
	Name           string  `db:"name" json:"name"`
	FirstName      string  `db:"first_name" json:"firstName"`
	LastName       string  `db:"last_name" json:"lastName"`
	Linkedin       string  `db:"linkedin_username" json:"linkedin"`
	JobTitle       *string `db:"job_title" json:"jobTitle"`
	CurrentCompany *string `db:"current_company" json:"currentCompany"`
}
type ReviewResponse struct {
	Id              uuid.UUID `db:"id" json:"id"`
	RecruiterId     uuid.UUID `db:"recruiter_id" json:"recruiterId"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	Rating          int       `db:"rating" json:"rating" validate:"min=0,max=5,required"`
	Description     string    `db:"description" json:"description" validate:"max=500"`
	ThumbsDownCount int       `db:"thumbs_down_count" json:"thumbsDownCount"`
	ThumbsUpCount   int       `db:"thumbs_up_count" json:"thumbs_down_count"`
}
type ReviewPayload struct {
	RecruiterId string `db:"recruiter_id" json:"recruiterId" validate:"required"`
	Rating      int    `db:"rating" json:"rating" validate:"min=0,max=5,required"`
	Description string `db:"description" json:"description" validate:"max=500"`
}
