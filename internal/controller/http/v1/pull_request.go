package v1

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/service"
	"github.com/gofiber/fiber/v3"
)

type prRouter struct {
	pr service.PullRequest
}

func newPRRouter(g fiber.Router, pr service.PullRequest) {
	r := &prRouter{
		pr: pr,
	}

	g.Post("/create", r.create)
	g.Post("/merge", r.merge)
	g.Post("/reassign", r.reassign)
}

type prCreateInput struct {
	ID       string `json:"pull_request_id" validate:"required"`
	Name     string `json:"pull_request_name" validate:"required"`
	AuthorId string `json:"author_id" validate:"required"`
}

type prResponse struct {
	PR domain.PullRequest `json:"pr"`
}

func (r *prRouter) create(c fiber.Ctx) error {
	var input prCreateInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	pr, err := r.pr.Create(c.Context(), input.ID, input.Name, input.AuthorId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(prResponse{PR: pr})
}

type prMergeInput struct {
	ID string `json:"pull_request_id" validate:"required"`
}

func (r *prRouter) merge(c fiber.Ctx) error {
	var input prMergeInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	pr, err := r.pr.Merge(c.Context(), input.ID)
	if err != nil {
		return err
	}
	return c.JSON(prResponse{PR: pr})
}

type prReassignInput struct {
	ID         string `json:"pull_request_id" validate:"required"`
	ReviewerId string `json:"old_reviewer_id" validate:"required"`
}

func (r *prRouter) reassign(c fiber.Ctx) error {
	var input prReassignInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	newReviewer, pr, err := r.pr.Reassign(c.Context(), input.ID, input.ReviewerId)
	if err != nil {
		return err
	}
	response := struct {
		prResponse
		ReplacedBy string `json:"replaced_by"`
	}{
		prResponse: prResponse{
			PR: pr,
		},
		ReplacedBy: newReviewer,
	}
	return c.JSON(response)
}
