package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"kilimanjaro-api/database"
	"kilimanjaro-api/server/middleware"
	"net"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"kilimanjaro-api/config"
)

type Server struct {
	router *mux.Router
	server *http.Server
}

func NewServer(prefix string) (*Server, error) {
	cfg := config.GetConfig()

	mainRouter := mux.NewRouter()
	router := mainRouter.PathPrefix(prefix).Subrouter()
	router.Use(middleware.JwtAuthentication)
	if cfg.AppDebug == "true" {
		router.Use(middleware.NewLogger(middleware.LogOptions{
			Formatter: &logrus.TextFormatter{
				DisableTimestamp: true,
				ForceColors:      true,
			},
			EnableStarting: true,
		}).Logger)
	}

	database.InitDatabase()
	//utils.InitialMigration()

	s := &Server{
		router: router,
	}

	s.SetupRoutes()

	return s, nil
}

func (s *Server) ListenAndServe() error {
	cfg := config.GetConfig()

	s.server = &http.Server{
		Addr:    net.JoinHostPort(cfg.AppDomain, cfg.AppPort),
		Handler: handlers.CompressHandler(s.router),
	}

	err := s.server.ListenAndServe()

	fmt.Println("Listening on localhost")

	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
	}

	return nil
}
