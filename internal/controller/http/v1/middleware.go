package v1

import (
	"avito_intership_2025/internal/domain"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

const (
	CodeTeamExists  = "TEAM_EXISTS"
	CodePRExists    = "PR_EXISTS"
	CodePRMerged    = "PR_ASSIGNED"
	CodeNotAssigned = "NOT_ASSIGNED"
	CodeNoCandidate = "NO_CANDIDATE"
	CodeNotFound    = "NOT_FOUND"
)

type errBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errBody `json:"error"`
}

func newErrResp(code, message string) errorResponse {
	return errorResponse{
		Error: errBody{
			Code:    code,
			Message: message,
		},
	}
}

func errorMiddleware(c fiber.Ctx) error {
	err := c.Next()
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(domain.ErrNotFound, err):
		return c.Status(fiber.StatusNotFound).JSON(newErrResp(CodeNotFound, "resource not found"))

	case errors.Is(err, domain.ErrPRExists):
		return c.Status(fiber.StatusConflict).JSON(newErrResp(CodePRExists, "PR id already exists"))

	case errors.Is(err, domain.ErrPRMerged):
		return c.Status(fiber.StatusConflict).JSON(newErrResp(CodePRMerged, "cannot reassign on merged PR"))

	case errors.Is(err, domain.ErrReviewerNotAssigned):
		return c.Status(fiber.StatusConflict).JSON(newErrResp(CodeNotAssigned, "reviewer is not assigned to this PR"))

	case errors.Is(err, domain.ErrNoCandidate):
		return c.Status(fiber.StatusConflict).JSON(newErrResp(CodeNoCandidate, "no active replacement candidate in team"))

	}
	log.Err(err).Msg("error middleware")
	return c.SendStatus(fiber.StatusInternalServerError)
}
