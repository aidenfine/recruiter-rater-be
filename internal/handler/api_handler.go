package handler

import (
	"fmt"
	"net/http"

	"github.com/aidenfine/recruiter-rater/internal/common"
	"github.com/aidenfine/recruiter-rater/internal/db/repositories"
	"github.com/aidenfine/recruiter-rater/internal/models"
	"github.com/aidenfine/recruiter-rater/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ApiHandler struct {
	V1Repository repositories.V1RepositoryInterface
}

func NewApiHandler(vr repositories.V1RepositoryInterface) *ApiHandler {
	return &ApiHandler{
		V1Repository: vr,
	}
}

var validate = validator.New()

func (h *ApiHandler) GetRecruiter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := utils.ParseQueryParams[models.GetRecruiterQueryParams](r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := validate.Struct(params); err != nil {
			common.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if params.Url {
			recruiter, err := h.V1Repository.SearchRecruiters(params.Search)
			if err != nil {
				if err.Error() == "not_found" {
					fmt.Printf("[RECRUITER] Recruiter with request: %s not found prompt to create a new recruiter", params.Search)
					fmt.Printf("[LOG] Request made by %s", r.RemoteAddr)
					common.WriteJSON(w, http.StatusNotFound, "NOT_FOUND")
					return
				} else {
					fmt.Print(err.Error())
					common.Error(w, http.StatusInternalServerError, "internal server error")
					return
				}
			}
			common.WriteJSON(w, http.StatusOK, recruiter)
		} else if !params.Url {
			return
		} else {
			common.Error(w, http.StatusBadRequest, "unexpected query param")
			return
		}

	}
}

func (h *ApiHandler) GetRecruiterById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			common.Error(w, http.StatusBadRequest, "missing id")
			return
		}

		recruiter, err := h.V1Repository.GetRecruiterById(id)
		if err != nil {
			fmt.Print(err.Error())
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		fmt.Printf("[RECRUITER] Getting recruiter %s, id: %s", recruiter.Name, recruiter.Id)
		fmt.Printf("[LOG] Request made by %s", r.RemoteAddr)

		common.WriteJSON(w, http.StatusOK, recruiter)

	}
}
func (h *ApiHandler) AddNewRecruiter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := common.DecodeJSONBody[models.RecruiterPayload](w, r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, "failed to decode json")
			return
		}

		if err := validate.Struct(payload); err != nil {
			common.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.V1Repository.AddNewRecruiter(&payload)
		if err != nil {
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		fmt.Println("[RECRUITER] Someone made a request to add a new recruiter with params %v\n", payload)
		fmt.Printf("[LOG] Request made by %s", r.RemoteAddr)

		common.WriteJSON(w, http.StatusCreated, "OK")

	}
}
func (h *ApiHandler) GetRecruiterReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := utils.ParseQueryParams[models.GetReviewsQueryParams](r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, "failed to decode json")
			return
		}

		fmt.Printf("[REVIEWS] Getting recruiter reviews for %s", params.RecruiterId)
		fmt.Printf("[LOG] Request made by %s", r.RemoteAddr)

		if err := validate.Struct(params); err != nil {
			common.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		recruiterUUID, err := uuid.Parse(params.RecruiterId)
		if err != nil {
			fmt.Println("error parsing uuid", err.Error())
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		items, err := h.V1Repository.GetAllRecruiterReviews(recruiterUUID, params.Limit)
		if err != nil {
			fmt.Println("error during query \n", err.Error())
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		common.WriteJSON(w, http.StatusOK, items)

	}
}

func (h *ApiHandler) AddNewReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := common.DecodeJSONBody[models.ReviewPayload](w, r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, "failed to decode json")
			return
		}
		if err := validate.Struct(payload); err != nil {
			common.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		recruiter, err := h.V1Repository.GetRecruiterById(payload.RecruiterId)
		if err != nil {
			fmt.Println(err.Error())
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		err = h.V1Repository.AddNewReview(&payload)
		if err != nil {
			fmt.Print(err.Error())
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// do this after successful added review
		updatedRecruiterRatingCount := recruiter.NumberOfRatings + 1
		updatedRecruiterRatingSum := recruiter.RatingSum + payload.Rating
		updatedRecruiterRating := updatedRecruiterRatingSum / updatedRecruiterRatingCount

		err = h.V1Repository.UpdateRecruiterRating(updatedRecruiterRatingSum, updatedRecruiterRatingCount, updatedRecruiterRating, payload.RecruiterId)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error when updating recruiter rating")
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		common.WriteJSON(w, http.StatusCreated, "OK")
	}
}

func (h *ApiHandler) GetMostRecentReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recentReviews, err := h.V1Repository.GetMostRecentReviews(6)
		if err != nil {
			fmt.Println(err.Error())
			common.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		common.WriteJSON(w, http.StatusOK, recentReviews)
	}
}
