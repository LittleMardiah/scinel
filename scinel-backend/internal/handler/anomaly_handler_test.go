package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/LittleMardiah/scinel/internal/middleware"
	"github.com/LittleMardiah/scinel/internal/model"
)

func newTestAnomalyHandler() (*AnomalyHandler, *MockAnomalyRepository) {
	repo := new(MockAnomalyRepository)
	return NewAnomalyHandler(repo), repo
}

func TestAnomalyHandler_ListAnomalies(t *testing.T) {
	whID := uuid.New()
	skuID := uuid.New()

	anomalies := []model.Anomaly{
		{
			ID:          uuid.New(),
			WarehouseID: whID,
			SKUID:       skuID,
			Type:        "spike",
			Severity:    model.SeverityHigh,
			Description: "unusual spike",
			ZScore:      4.5,
			Status:      model.StatusOpen,
			DetectedAt:  time.Now().UTC(),
		},
	}

	t.Run("success", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		repo.On("GetAll", mock.Anything, 20, 0).Return(anomalies, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies", nil)
		rr := httptest.NewRecorder()
		h.ListAnomalies(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		list := data["anomalies"].([]interface{})
		assert.Equal(t, 1, len(list))
		repo.AssertExpectations(t)
	})

	t.Run("with status filter", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		repo.On("GetAll", mock.Anything, 20, 0).Return(anomalies, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies?status=open", nil)
		rr := httptest.NewRecorder()
		h.ListAnomalies(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		repo.On("GetAll", mock.Anything, 20, 0).Return([]model.Anomaly{}, int64(0), errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies", nil)
		rr := httptest.NewRecorder()
		h.ListAnomalies(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestAnomalyHandler_GetAnomaly(t *testing.T) {
	anomalyID := uuid.New()

	anomaly := &model.Anomaly{
		ID:          anomalyID,
		WarehouseID: uuid.New(),
		SKUID:       uuid.New(),
		SKUCode:     "SKU-001",
		Type:        "spike",
		Severity:    model.SeverityMedium,
		Description: "unusual activity",
		ZScore:      3.5,
		Status:      model.StatusOpen,
		DetectedAt:  time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	t.Run("success", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		repo.On("GetByID", mock.Anything, anomalyID).Return(anomaly, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies/"+anomalyID.String(), nil)
		rr := httptest.NewRecorder()
		h.GetAnomaly(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("missing id", func(t *testing.T) {
		h, _ := newTestAnomalyHandler()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies/", nil)
		rr := httptest.NewRecorder()
		h.GetAnomaly(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		h, _ := newTestAnomalyHandler()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies/not-a-uuid", nil)
		rr := httptest.NewRecorder()
		h.GetAnomaly(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("not found", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		repo.On("GetByID", mock.Anything, anomalyID).Return((*model.Anomaly)(nil), errors.New("anomaly not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies/"+anomalyID.String(), nil)
		rr := httptest.NewRecorder()
		h.GetAnomaly(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestAnomalyHandler_UpdateAnomalyStatus(t *testing.T) {
	anomalyID := uuid.New()
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "analyst"}

	t.Run("success acknowledge", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		current := &model.Anomaly{ID: anomalyID, Status: model.StatusOpen}
		updated := &model.Anomaly{ID: anomalyID, Status: model.StatusAcknowledged, UpdatedAt: time.Now().UTC()}

		repo.On("GetByID", mock.Anything, anomalyID).Return(current, nil).Once()
		repo.On("UpdateStatus", mock.Anything, anomalyID, model.AnomalyStatus("acknowledged"), "", userID).Return(nil)
		repo.On("GetByID", mock.Anything, anomalyID).Return(updated, nil).Once()

		body := `{"status":"acknowledged"}`
		req := requestWithUser(http.MethodPut, "/api/v1/anomalies/"+anomalyID.String(), body, user)
		rr := httptest.NewRecorder()
		h.UpdateAnomalyStatus(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestAnomalyHandler()
		body := `{"status":"acknowledged"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/anomalies/"+anomalyID.String(), bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.UpdateAnomalyStatus(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid status", func(t *testing.T) {
		h, _ := newTestAnomalyHandler()
		body := `{"status":"invalid"}`
		req := requestWithUser(http.MethodPut, "/api/v1/anomalies/"+anomalyID.String(), body, user)
		rr := httptest.NewRecorder()
		h.UpdateAnomalyStatus(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("already resolved conflict", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		current := &model.Anomaly{ID: anomalyID, Status: model.StatusResolved}
		repo.On("GetByID", mock.Anything, anomalyID).Return(current, nil)

		body := `{"status":"acknowledged"}`
		req := requestWithUser(http.MethodPut, "/api/v1/anomalies/"+anomalyID.String(), body, user)
		rr := httptest.NewRecorder()
		h.UpdateAnomalyStatus(rr, req)
		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("open to resolved invalid transition", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		current := &model.Anomaly{ID: anomalyID, Status: model.StatusOpen}
		repo.On("GetByID", mock.Anything, anomalyID).Return(current, nil)

		body := `{"status":"resolved"}`
		req := requestWithUser(http.MethodPut, "/api/v1/anomalies/"+anomalyID.String(), body, user)
		rr := httptest.NewRecorder()
		h.UpdateAnomalyStatus(rr, req)
		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("not found", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		repo.On("GetByID", mock.Anything, anomalyID).Return((*model.Anomaly)(nil), errors.New("not found"))

		body := `{"status":"acknowledged"}`
		req := requestWithUser(http.MethodPut, "/api/v1/anomalies/"+anomalyID.String(), body, user)
		rr := httptest.NewRecorder()
		h.UpdateAnomalyStatus(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestAnomalyHandler_GetAnomalyHistory(t *testing.T) {
	anomalyID := uuid.New()
	now := time.Now().UTC()

	t.Run("open anomaly timeline", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		anomaly := &model.Anomaly{
			ID:          anomalyID,
			Description: "desc",
			Status:      model.StatusOpen,
			DetectedAt:  now,
		}
		repo.On("GetByID", mock.Anything, anomalyID).Return(anomaly, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies/"+anomalyID.String()+"/history", nil)
		rr := httptest.NewRecorder()
		h.GetAnomalyHistory(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		timeline := data["timeline"].([]interface{})
		assert.Equal(t, 1, len(timeline))
		repo.AssertExpectations(t)
	})

	t.Run("invalid path", func(t *testing.T) {
		h, _ := newTestAnomalyHandler()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies/"+anomalyID.String()+"/wrong", nil)
		rr := httptest.NewRecorder()
		h.GetAnomalyHistory(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("not found", func(t *testing.T) {
		h, repo := newTestAnomalyHandler()
		repo.On("GetByID", mock.Anything, anomalyID).Return((*model.Anomaly)(nil), errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/anomalies/"+anomalyID.String()+"/history", nil)
		rr := httptest.NewRecorder()
		h.GetAnomalyHistory(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
