package v1

import (
	"avito_intership_2025/internal/service"
	"github.com/gofiber/fiber/v3"
)

func NewRouter(g fiber.Router, services *service.Services) {
	g.Get("/ping", ping)

	g.Use(errorMiddleware)

	newTeamRouter(g.Group("/team"), services.Team)
	newStatsRouter(g.Group("/stats"), services.Stats)
	newPRRouter(g.Group("/pullRequest"), services.PullRequest)
	newUserRouter(g.Group("/users"), services.User, services.PullRequest)
}

func ping(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
