package bill

import (
	"net/http"

	"github.com/gfmanica/splitz-backend/service/auth"
	"github.com/gfmanica/splitz-backend/types"
	"github.com/gfmanica/splitz-backend/utils"
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
}

func (h *Handler) handleGetBills(w http.ResponseWriter, r *http.Request) {
	bills, err := h.store.GetBills()

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, bills)
}
