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
	})
}
