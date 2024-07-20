package application

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/RhinoSC/sre-backend/internal/handler"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/RhinoSC/sre-backend/internal/repository"
	"github.com/RhinoSC/sre-backend/internal/service"
)

type ConfigServerChi struct {
	// Address is the address to listen on
	Address string
}

type ServerChi struct {
	address string
}

func NewServerChi(cfg ConfigServerChi) *ServerChi {
	defaultCfg := ConfigServerChi{
		Address: ":8080",
	}

	if cfg.Address != "" {
		defaultCfg.Address = cfg.Address
	}

	return &ServerChi{
		address: defaultCfg.Address,
	}
}

func (s *ServerChi) Run() (err error) {
	db, err := sql.Open("sqlite3", "./database.db?foreign_keys=")
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

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", handler.PingHandler())

	router.Route("/api/v1", func(r chi.Router) {
		buildUserRouter(&r, db)
		buildEventRouter(&r, db)
		buildPrizeRouter(&r, db)
		buildScheduleRouter(&r, db)
		buildRunRouter(&r, db)
	})

	err = http.ListenAndServe(s.address, router)

	return
}

func buildUserRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewUserSqlite(db)
	sv := service.NewUserDefault(rp)
	hd := handler.NewUserDefault(sv)

	(*router).Route("/users", func(rt chi.Router) {
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetById())
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
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetById())
		rt.Post("/", hd.Create())
		rt.Patch("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})
}

func buildPrizeRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewPrizeSqlite(db)
	sv := service.NewPrizeDefault(rp)
	hd := handler.NewPrizeDefault(sv)

	(*router).Route("/prizes", func(rt chi.Router) {
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetById())
		rt.Post("/", hd.Create())
		rt.Patch("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})
}

func buildScheduleRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewScheduleSqlite(db)
	sv := service.NewScheduleDefault(rp)
	hd := handler.NewScheduleDefault(sv)

	(*router).Route("/schedules", func(rt chi.Router) {
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetById())
		rt.Post("/", hd.Create())
		rt.Patch("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})
}

func buildRunRouter(router *chi.Router, db *sql.DB) {
	rp := repository.NewRunSqlite(db)
	sv := service.NewRunDefault(rp)
	hd := handler.NewRunDefault(sv)

	(*router).Route("/runs", func(rt chi.Router) {
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())
		rt.Post("/", hd.Create())
		rt.Patch("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})
}
