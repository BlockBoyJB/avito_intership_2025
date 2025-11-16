package v1

import (
	"avito_intership_2025/internal/service"
	"github.com/gofiber/fiber/v3"
)

type statsRouter struct {
	stats service.Stats
}

func newStatsRouter(g fiber.Router, stats service.Stats) {
	r := &statsRouter{
		stats: stats,
	}

	g.Get("", r.get)
}

func (r *statsRouter) get(c fiber.Ctx) error {
	var (
		result any
		err    error
	)
	switch c.Query("filter") {
	case "reviewers":
		result, err = r.stats.Reviewers(c.Context())
	case "author":
		result, err = r.stats.ByAuthor(c.Context())
	case "team":
		result, err = r.stats.ByTeam(c.Context())
	default:
		return c.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.JSON(result)
}
