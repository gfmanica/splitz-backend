package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gfmanica/splitz-backend/service/bill"
	"github.com/gfmanica/splitz-backend/service/ride"
	"github.com/gfmanica/splitz-backend/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(add string, db *sql.DB) *APIServer {
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

	billStore := bill.NewStore(s.db)
	billHandler := bill.NewHandler(billStore, userStore)
	billHandler.RegisterRoutes(subrouter)

	rideStore := ride.NewStore(s.db)
	rideHandler := *ride.NewHandler(rideStore, userStore)
	rideHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
