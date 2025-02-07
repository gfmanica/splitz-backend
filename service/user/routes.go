package user

import (
	"fmt"
	"net/http"

	"github.com/gfmanica/splitz-backend/config"
	"github.com/gfmanica/splitz-backend/service/auth"
	"github.com/gfmanica/splitz-backend/types"
	"github.com/gfmanica/splitz-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get the JSON payload
	var payload types.LoginUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)

		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", error))

		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("user with email %s or password  not found", payload.Email))

		return
	}

	fmt.Print(u.Email)
	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("user with email %s or password  not found", payload.Email))

		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)

		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get the JSON payload
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)

		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", error))

		return
	}

	// check if the user already exists
	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("user with email %s exists", payload.Email))

		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
	}

	// create the user
	err = h.store.CreateUser(types.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
