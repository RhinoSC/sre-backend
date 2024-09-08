package application

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	_ "github.com/mattn/go-sqlite3"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/auth"
	"github.com/RhinoSC/sre-backend/internal/handler"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/RhinoSC/sre-backend/internal/repository"
	"github.com/RhinoSC/sre-backend/internal/service"
)

type ConfigServerChi struct {
	// Address is the address to listen on
	Address   string
	JWTSecret string
}

type ServerChi struct {
	address   string
	jwtSecret string
}

func NewServerChi(cfg ConfigServerChi) *ServerChi {
	defaultCfg := ConfigServerChi{
		Address:   ":8080",
		JWTSecret: "defaultsecret",
	}

	if cfg.Address != "" {
		defaultCfg.Address = cfg.Address
	}

	if cfg.JWTSecret != "" {
		defaultCfg.JWTSecret = cfg.JWTSecret
	}

	return &ServerChi{
		address:   defaultCfg.Address,
		jwtSecret: defaultCfg.JWTSecret,
	}
}

func (s *ServerChi) Run() (err error) {

	workingDir, err := os.Getwd()
	rootDir := filepath.Join(workingDir)
	filePath := filepath.Join(rootDir, "database.db?_foreign_keys=on")

	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return
	}

	// initialize logger
	logger.InitializeLogger()

	twitch := internal.Twitch{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		ClientToken:  "",
	}

	// initilize twitch
	service.CreateFirstTime(&twitch)

	// Initialize JWT Auth
	auth.Init(s.jwtSecret)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://speedrunespanol.com", "https://speedrunespanol.com", "https://*.speedrunespanol.com","http://*.speedrunespanol.com"},
		AllowedOrigins:   []string{"https://", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.Verifier())

	router.Get("/ping", handler.PingHandler())

	router.Route("/api/v1", func(r chi.Router) {

		r.With(auth.Authenticator()).Group(func(rt chi.Router) {
			buildUserRouter(&rt, db)
			buildTeamRouter(&rt, db)
		})

		// Algunas rutas privadas y otras publicas
		r.Group(func(r chi.Router) {
			buildPrizeRouter(&r, db)
			buildEventRouter(&r, db)
			buildAdminRouter(&r, db)
			buildScheduleRouter(&r, db)
			buildRunRouter(&r, db)

			buildBidRouter(&r, db)
			buildDonationRouter(&r, db)
		})
	})

	err = http.ListenAndServe(s.address, router)

	return
}

func buildAdminRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewAdminSqlite(db)
	sv := service.NewAdminDefault(rp)
	hd := handler.NewAdminDefault(sv)

	(*router).Route("/admins", func(rt chi.Router) {

		// Public
		rt.Post("/login", hd.Login())

		// Private
		rt.With(auth.Authenticator()).Group(func(r chi.Router) {
			r.Get("/", hd.GetAll())
			r.Get("/{id}", hd.GetByID())
			r.Post("/", hd.Create())
			r.Patch("/{id}", hd.Update())
			r.Delete("/{id}", hd.Delete())

			r.Get("/validate", hd.ValidateToken())
		})
	})
}
func buildUserRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewUserSqlite(db)
	sv := service.NewUserDefault(rp)
	hd := handler.NewUserDefault(sv)

	(*router).Route("/users", func(rt chi.Router) {
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())
		rt.Get("/username/{username}", hd.GetByUsername())
		rt.Post("/", hd.Create())
		rt.Patch("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})
}

func buildEventRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewEventSqlite(db)
	sv := service.NewEventDefault(rp)
	hd := handler.NewEventDefault(sv)

	(*router).Route("/events", func(rt chi.Router) {

		// Public
		rt.Get("/info", hd.GetBasicInfo())
		rt.Get("/{id}", hd.GetByID())

		// Private
		rt.With(auth.Authenticator()).Group(func(r chi.Router) {
			r.Get("/", hd.GetAll())
			r.Post("/", hd.Create())
			r.Patch("/{id}", hd.Update())
			r.Delete("/{id}", hd.Delete())
		})
	})
}

func buildPrizeRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewPrizeSqlite(db)
	sv := service.NewPrizeDefault(rp)
	hd := handler.NewPrizeDefault(sv)

	(*router).Route("/prizes", func(rt chi.Router) {

		// public
		rt.Get("/", hd.GetAll())

		// private
		rt.With(auth.Authenticator()).Group(func(r chi.Router) {
			r.Get("/{id}", hd.GetByID())
			r.Post("/", hd.Create())
			r.Patch("/{id}", hd.Update())
			r.Delete("/{id}", hd.Delete())
		})
	})
}

func buildScheduleRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewScheduleSqlite(db)
	sv := service.NewScheduleDefault(rp)
	hd := handler.NewScheduleDefault(sv)

	(*router).Route("/schedules", func(rt chi.Router) {

		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())

		rt.With(auth.Authenticator()).Group(func(r chi.Router) {
			r.Post("/", hd.Create())
			r.Patch("/{id}", hd.Update())
			r.Delete("/{id}", hd.Delete())
		})
	})
}

func buildRunRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewRunSqlite(db)
	sv := service.NewRunDefault(rp)
	hd := handler.NewRunDefault(sv)

	(*router).Route("/runs", func(rt chi.Router) {

		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())

		rt.With(auth.Authenticator()).Group(func(r chi.Router) {
			r.Post("/", hd.Create())
			r.Patch("/{id}", hd.Update())
			r.Delete("/{id}", hd.Delete())
			r.Post("/order", hd.UpdateRunOrder())

			// Twitch
			r.Group(func(r chi.Router) {
				r.Use(httprate.Limit(13, time.Minute, httprate.WithResponseHeaders(httprate.ResponseHeaders{
					Limit:      "X-RateLimit-Limit",
					Remaining:  "X-RateLimit-Remaining",
					Reset:      "X-RateLimit-Reset",
					RetryAfter: "Retry-After",
				}), httprate.WithKeyFuncs(
					httprate.KeyByEndpoint,
				)))
				r.Get("/twitch/categories", hd.FindTwitchCategories())
				r.Get("/twitch/game", hd.FindTwitchCategoryByID())
			})
		})
	})
}

func buildTeamRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewTeamSqlite(db)
	sv := service.NewTeamDefault(rp)
	hd := handler.NewTeamDefault(sv)

	(*router).Route("/teams", func(rt chi.Router) {
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())
		rt.Post("/", hd.Create())
		rt.Patch("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})
}

func buildBidRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewBidSqlite(db)
	sv := service.NewBidDefault(rp)
	hd := handler.NewBidDefault(sv)

	(*router).Route("/bids", func(rt chi.Router) {

		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())

		rt.With(auth.Authenticator()).Group(func(r chi.Router) {
			r.Post("/", hd.Create())
			r.Patch("/{id}", hd.Update())
			r.Delete("/{id}", hd.Delete())
		})
	})
}

func buildDonationRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewDonationSqlite(db)
	sv := service.NewDonationDefault(rp)
	hd := handler.NewDonationDefault(sv)

	(*router).Route("/donations", func(rt chi.Router) {

		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())
		rt.Get("/event/{id}", hd.GetByEventID())
		rt.Post("/", hd.Create())

		rt.With(auth.Authenticator()).Group(func(r chi.Router) {
			rt.Patch("/{id}", hd.Update())
			r.Delete("/{id}", hd.Delete())
		})
	})
}
