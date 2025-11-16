package v1

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/service"
	"github.com/gofiber/fiber/v3"
)

type userRouter struct {
	user service.User
	pr   service.PullRequest
}

func newUserRouter(g fiber.Router, user service.User, pr service.PullRequest) {
	r := &userRouter{
		user: user,
		pr:   pr,
	}

	g.Post("/setIsActive", r.setIsActive)
	g.Get("/getReview", r.getReview)
}

type userIsActiveInput struct {
	ID       string `json:"user_id" validate:"required"`
	IsActive bool   `json:"is_active"`
}

func (r *userRouter) setIsActive(c fiber.Ctx) error {
	var input userIsActiveInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, err := r.user.SetIsActive(c.Context(), input.ID, input.IsActive)
	if err != nil {
		return err
	}
	response := struct {
		User domain.User `json:"user"`
	}{
		User: user,
	}
	return c.JSON(response)
}

func (r *userRouter) getReview(c fiber.Ctx) error {
	userID := c.Query("user_id")

	review, err := r.pr.GetUserReview(c.Context(), userID)
	if err != nil {
		return err
	}

	response := struct {
		ID string               `json:"user_id"`
		PR []domain.PullRequest `json:"pull_requests"`
	}{
		ID: userID,
		PR: review,
	}
	return c.JSON(response)
}
