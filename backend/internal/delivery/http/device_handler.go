package http

import (
	"log/slog"
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

type DeviceHandler struct {
	deviceUC domain.DeviceUseCasePort
}

func NewDeviceHandler(deviceUC domain.DeviceUseCasePort) *DeviceHandler {
	return &DeviceHandler{deviceUC: deviceUC}
}

// ListAll godoc
//
//	@Summary	List all devices
//	@Tags		devices
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{array}		DeviceResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/devices [get]
func (h *DeviceHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceUC.ListAll(r.Context())
	if err != nil {
		slog.Error("list all devices failed", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list devices")
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(devices, toDeviceResponse))
}

// ListAvailable godoc
//
//	@Summary	List available (online) devices
//	@Tags		devices
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{array}		DeviceResponse
//	@Failure	401	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/devices/available [get]
func (h *DeviceHandler) ListAvailable(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceUC.ListAvailable(r.Context())
	if err != nil {
		slog.Error("list available devices failed", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list available devices")
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(devices, toDeviceResponse))
}
