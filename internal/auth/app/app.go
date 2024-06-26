package app

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/Michman/config"
	"github.com/imirjar/Michman/internal/auth/app/middleware"
	"github.com/imirjar/Michman/internal/auth/models"
	"github.com/imirjar/Michman/internal/auth/service"
)

type Service interface {
	BuildJWTString(context.Context, models.User) (string, error)
	GetUserID(context.Context, string) int
	VerifyToken(ctx context.Context, tokenString string) error
}

type Config interface {
	GetSecret() string
	GetAuthAddr() string
}

type App struct {
	config  Config
	service Service
}

func NewApp() *App {

	return &App{
		config:  config.NewConfig(),
		service: service.NewService(),
	}
}

func (a *App) Run(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(middleware.REST())

	router.Get("/", a.Hello)

	router.Route("/token", func(token chi.Router) {
		token.Post("/create/", a.CreateJWT)
		token.Post("/view/", a.ValidateJWT)
	})

	router.Route("/user", func(auth chi.Router) {
		auth.Post("/create/", a.CreateUser)
	})

	srv := &http.Server{
		Addr:    a.config.GetAuthAddr(),
		Handler: router,
	}

	log.Printf("Run app on PORT %s", a.config.GetAuthAddr())
	return srv.ListenAndServe()
}
