package handlers

import (
	"health-checker/internal/application/usecases"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/infra/http/helpers"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type MonitorHandler struct {
	createMonitorUseCase usecases.CreateMonitorUseCase
	getMonitorsUseCase   usecases.GetMonitorsUseCase
	toggleMonitorUseCase usecases.ToggleMonitorUseCase
}

func NewMonitorHandler(
	createMonitorUseCase usecases.CreateMonitorUseCase,
	getMonitorsUseCase usecases.GetMonitorsUseCase,
	toggleMonitorUseCase usecases.ToggleMonitorUseCase,
) *MonitorHandler {
	return &MonitorHandler{
		createMonitorUseCase: createMonitorUseCase,
		getMonitorsUseCase:   getMonitorsUseCase,
		toggleMonitorUseCase: toggleMonitorUseCase,
	}
}

type CreateMonitorRequest struct {
	Name               string            `json:"name" validate:"required,min=3"`
	URL                string            `json:"url" validate:"required,url"`
	Method             string            `json:"method" validate:"required,oneof=GET POST HEAD"`
	Headers            map[string]string `json:"headers" validate:"omitempty,dive,keys,ascii,lowercase"`
	Body               string            `json:"body" validate:"omitempty"`
	Interval           int               `json:"interval" validate:"required,min=1"`
	ExpectedStatusCode int32             `json:"expected_status_code" validate:"omitempty,min=100,max=599"`
	Timeout            int               `json:"timeout" validate:"required,min=1"`
}

func (h *MonitorHandler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	req, err := helpers.DecodeAndValidateRequest[CreateMonitorRequest](w, r)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = h.createMonitorUseCase.Execute(r.Context(), usecases.CreateMonitorCommand{
		UserID:             userID,
		Name:               req.Name,
		URL:                req.URL,
		Method:             req.Method,
		Headers:            req.Headers,
		Body:               req.Body,
		Interval:           req.Interval,
		ExpectedStatusCode: req.ExpectedStatusCode,
		Timeout:            req.Timeout,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Monitor created successfully"})
}

func (h *MonitorHandler) GetMonitors(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit int32 = 10
	var offset int32 = 0

	if limitStr != "" {
		parsedLimit, parseErr := strconv.Atoi(limitStr)
		if parseErr != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid limit")
			return
		}
		limit = int32(parsedLimit)
	}

	if offsetStr != "" {
		parsedOffset, parseErr := strconv.Atoi(offsetStr)
		if parseErr != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid offset")
			return
		}
		offset = int32(parsedOffset)
	}

	cmd := usecases.GetMonitorsCommand{
		UserID: userID,
		Status: entities.MonitorStatus(status),
		Limit:  limit,
		Offset: offset,
	}

	monitorData, err := h.getMonitorsUseCase.Execute(r.Context(), cmd)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := helpers.NewPaginatedResponse(monitorData.Monitors, monitorData.Total, int(limit), int(offset))

	helpers.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *MonitorHandler) ToggleMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		helpers.WriteError(w, http.StatusBadRequest, "Monitor ID is required")
		return
	}

	err := h.toggleMonitorUseCase.Execute(r.Context(), uuid.MustParse(id))
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Monitor toggled successfully"})
}
