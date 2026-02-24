package handlers

import (
	"health-checker/internal/application/usecases"
	"health-checker/internal/infra/http/helpers"
	"net/http"

	"github.com/google/uuid"
)

type MonitorHandler struct {
	createMonitorUseCase usecases.CreateMonitorUseCase
}

func NewMonitorHandler(createMonitorUseCase usecases.CreateMonitorUseCase) *MonitorHandler {
	return &MonitorHandler{createMonitorUseCase: createMonitorUseCase}
}

type CreateMonitorRequest struct {
	Name               string            `json:"name" validate:"required,min=3"`
	URL                string            `json:"url" validate:"required,url"`
	Method             string            `json:"method" validate:"required,oneof=GET POST HEAD"`
	Headers            map[string]string `json:"headers" validate:"omitempty,dive,keys,ascii,lowercase"`
	Body               string            `json:"body" validate:"omitempty"`
	Interval           int               `json:"interval" validate:"required,min=1"`
	ExpectedStatusCode uint32            `json:"expected_status_code" validate:"omitempty,min=100,max=599"`
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

	helpers.WriteJSONResponse(w, http.StatusCreated, "Monitor created successfully")
}
