package v1

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/service"
	"errors"
	"github.com/gofiber/fiber/v3"
)

type teamRouter struct {
	team service.Team
}

func newTeamRouter(g fiber.Router, team service.Team) {
	r := &teamRouter{
		team: team,
	}

	g.Post("/add", r.add)
	g.Get("/get", r.get)
}

func (r *teamRouter) add(c fiber.Ctx) error {
	var input domain.Team

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := r.team.Create(c.Context(), input); err != nil {
		if errors.Is(err, domain.ErrTeamExists) {
			return c.Status(fiber.StatusBadRequest).JSON(newErrResp(CodeTeamExists, input.Name+" already exists"))
		}
		return err
	}
	response := struct {
		Team domain.Team `json:"team"`
	}{
		Team: input,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (r *teamRouter) get(c fiber.Ctx) error {
	team, err := r.team.Find(c.Context(), c.Query("team_name"))
	if err != nil {
		return err
	}
	return c.JSON(team)
}
