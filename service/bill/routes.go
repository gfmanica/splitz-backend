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

func NewHandler(store types.BillStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/bill", auth.WithJWTAuth(h.handleGetBills, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/bill/{id}", auth.WithJWTAuth(h.handleGetBill, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/bill", auth.WithJWTAuth(h.handleCreateBill, h.userStore)).Methods(http.MethodPost)
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
	})

	if err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
