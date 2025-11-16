package app

import (
	"avito_intership_2025/config"
	httpv1 "avito_intership_2025/internal/controller/http/v1"
	"avito_intership_2025/internal/repo"
	"avito_intership_2025/internal/service"
	"avito_intership_2025/pkg/postgres"
	"avito_intership_2025/pkg/postgres/txmanager"
	"avito_intership_2025/pkg/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("config init error")
	}
	setLogger(cfg.Log.Level, cfg.Log.Output)

	pg, err := postgres.NewPG(cfg.PG.Url)
	if err != nil {
		log.Fatal().Err(err).Msg("init postgres error")
	}
	defer pg.Close()

	tx := txmanager.NewManager(pg)

	d := &service.ServicesDependencies{
		Tx:          tx,
		User:        repo.NewUserRepo(tx),
		Team:        repo.NewTeamRepo(tx),
		PullRequest: repo.NewPRRepo(tx),
		Stats:       repo.NewStatsRepo(tx),
	}

	services := service.NewServices(d)

	h := fiber.New(fiber.Config{
		StructValidator: validator.New(),
	})

	httpv1.NewRouter(h, services)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	handlerCh := make(chan error, 1)
	go func() {
		handlerCh <- h.Listen(net.JoinHostPort("", cfg.HTTP.Port))
	}()

	log.Info().Msgf("app started, listen port %s", cfg.HTTP.Port)

	select {
	case s := <-interrupt:
		log.Info().Msgf("app signal %s", s.String())
	case err = <-handlerCh:
		log.Err(err).Msg("http server error")
	}

	if err = h.Shutdown(); err != nil {
		log.Err(err).Msg("http server shutdown error")
	}

	log.Info().Msg("app shutdown with exit code 0")
}

func init() {
	if _, ok := os.LookupEnv("HTTP_PORT"); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("load env file error")
		}
	}
}
