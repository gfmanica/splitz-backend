package api

import (
	"log"
	"net/http"

	"github.com/gfmanica/splitz-backend/service/user"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type APIServer struct {
	addr string
	db   *pgx.Conn
}

func NewAPIServer(add string, db *pgx.Conn) *APIServer {
	return &APIServer{
		addr: add,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
