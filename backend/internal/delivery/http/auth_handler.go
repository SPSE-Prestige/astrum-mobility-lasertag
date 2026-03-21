package http

import (
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type AuthHandler struct {
	authUC *usecase.AuthUseCase
}

func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// Login godoc
//
//	@Summary	Admin login
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		LoginRequest	true	"Credentials"
//	@Success	200		{object}	LoginResponse
//	@Failure	400		{object}	ErrorResponse
//	@Failure	401		{object}	ErrorResponse
//	@Router		/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password required")
		return
	}

	session, err := h.authUC.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	})
}

// Logout godoc
//
//	@Summary	Logout (invalidate token)
//	@Tags		auth
//	@Security	BearerAuth
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Router		/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) > 7 {
		token = token[7:] // strip "Bearer "
	}
	_ = h.authUC.Logout(r.Context(), token)
	w.WriteHeader(http.StatusNoContent)
}
