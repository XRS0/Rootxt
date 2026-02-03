package handlers

import (
	"encoding/json"
	"net/http"

	appUser "github.com/rootix/portfolio/internal/application/user"
	"github.com/rootix/portfolio/internal/infrastructure/auth"
)

type AdminAuthHandler struct {
	LoginUseCase appUser.LoginUseCase
	JWT          auth.JWTManager
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	userEntity, err := h.LoginUseCase.Execute(r.Context(), appUser.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := h.JWT.Generate(userEntity.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}
