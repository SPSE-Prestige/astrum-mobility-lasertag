package http

import (
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type AdminHandler struct {
	adminUC *usecase.AdminUseCase
}

func NewAdminHandler(adminUC *usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUC: adminUC}
}

// Login handles admin authentication.
// @Summary      Admin login
// @Description  Authenticate an admin user and receive a bearer token.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        body body LoginRequest true "Credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /api/admin/login [post]
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password required")
		return
	}

	token, err := h.adminUC.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
