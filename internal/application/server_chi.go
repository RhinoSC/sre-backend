package application

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/RhinoSC/sre-backend/internal/handler"
	"github.com/RhinoSC/sre-backend/internal/logger"
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
	db, err := sql.Open("sqlite3", "./database.db")
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

	err = http.ListenAndServe(s.address, router)

	return
}
