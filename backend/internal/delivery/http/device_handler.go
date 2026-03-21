package http

import (
	"net/http"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

type DeviceHandler struct {
	deviceUC *usecase.DeviceUseCase
}

func NewDeviceHandler(deviceUC *usecase.DeviceUseCase) *DeviceHandler {
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
		writeError(w, http.StatusInternalServerError, err.Error())
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mapSlice(devices, toDeviceResponse))
}
