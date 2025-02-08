package ride

import (
	// "fmt"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gfmanica/splitz-backend/service/auth"
	"github.com/gfmanica/splitz-backend/types"
	"github.com/gfmanica/splitz-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.RideStore
	userStore types.UserStore
}

func NewHandler(store types.RideStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/ride", auth.WithJWTAuth(h.handleGetRides, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/ride/{id}", auth.WithJWTAuth(h.handleGetRide, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/ride", auth.WithJWTAuth(h.handleCreateRide, h.userStore)).Methods(http.MethodPost)
	// router.HandleFunc("/ride", auth.WithJWTAuth(h.handleUpdateRide, h.userStore)).Methods(http.MethodPut)
}

func (h *Handler) handleGetRides(w http.ResponseWriter, r *http.Request) {
	Rides, err := h.store.GetRides()

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, Rides)
}

func (h *Handler) handleGetRide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	Ride, err := h.store.GetRideById(id)

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, Ride)
}

// func (h *Handler) handleUpdateRide(w http.ResponseWriter, r *http.Request) {
// 	var payload types.Ride

// 	if err := utils.ParseJSON(r, &payload); err != nil {
// 		utils.WriterError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	if err := utils.Validate.Struct(payload); err != nil {
// 		error := err.(validator.ValidationErrors)
// 		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", error))
// 		return
// 	}

// 	Ride := types.Ride{
// 		IdRide:   payload.IdRide,
// 		DsRide:   payload.DsRide,
// 		VlRide:   payload.VlRide,
// 		QtPerson: payload.QtPerson,
// 		Payments: convertToRidePayments(payload.Payments),
// 	}

// 	if err := h.store.UpdateRide(Ride); err != nil {
// 		utils.WriterError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, Ride)
// }

func (h *Handler) handleCreateRide(w http.ResponseWriter, r *http.Request) {
	// get the JSON payload
	var payload types.CreateRidePayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", error))

		return
	}

	err := h.store.CreateRide(types.Ride{
		DsRide:   payload.DsRide,
		VlRide:   payload.VlRide,
		DtInit:   payload.DtInit,
		DtFinish:   payload.DtFinish,
		FgCountWeekend: payload.FgCountWeekend,
		Payments: convertToRidePayments(payload.Payments),
	})

	if err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func convertToRidePayments(createPayments []types.RidePayment) []types.RidePayment {
	ridePayments := make([]types.RidePayment, len(createPayments))

	for i, createPayment := range createPayments {
		ridePayments[i] = types.RidePayment{
			DsPerson: createPayment.DsPerson,
		}
	}
	return ridePayments
}
