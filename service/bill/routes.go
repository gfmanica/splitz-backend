package bill

import (
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
	store     types.BillStore
	userStore types.UserStore
}

func convertToBillPayments(createPayments []types.BillPayment) []types.BillPayment {
	billPayments := make([]types.BillPayment, len(createPayments))

	for i, createPayment := range createPayments {
		billPayments[i] = types.BillPayment{
			DsPerson:        createPayment.DsPerson,
			VlPayment:       createPayment.VlPayment,
			FgPayed:         createPayment.FgPayed,
			FgCustomPayment: createPayment.FgCustomPayment,
		}
	}
	return billPayments
}

func NewHandler(store types.BillStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/bill", auth.WithJWTAuth(h.handleGetBills, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/bill", auth.WithJWTAuth(h.handleCreateBill, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/bill", auth.WithJWTAuth(h.handleUpdateBill, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/bill/{id}", auth.WithJWTAuth(h.handleGetBill, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/bill/{id}", auth.WithJWTAuth(h.handleDeleteBill, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetBills(w http.ResponseWriter, r *http.Request) {
	bills, err := h.store.GetBills()

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, bills)
}

func (h *Handler) handleGetBill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	bill, err := h.store.GetBillById(id)

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, bill)
}

func (h *Handler) handleUpdateBill(w http.ResponseWriter, r *http.Request) {
	var payload types.Bill

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", error))
		return
	}

	bill := types.Bill{
		IdBill:   payload.IdBill,
		DsBill:   payload.DsBill,
		VlBill:   payload.VlBill,
		QtPerson: payload.QtPerson,
		Payments: convertToBillPayments(payload.Payments),
	}

	if err := h.store.UpdateBill(bill); err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, bill)
}

func (h *Handler) handleCreateBill(w http.ResponseWriter, r *http.Request) {
	// get the JSON payload
	var payload types.CreateBillPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", error))

		return
	}

	err := h.store.CreateBill(types.Bill{
		DsBill:   payload.DsBill,
		VlBill:   payload.VlBill,
		QtPerson: payload.QtPerson,
		Payments: convertToBillPayments(payload.Payments),
	})

	if err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}

func (h *Handler) handleDeleteBill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := h.store.DeleteBill(id)

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
