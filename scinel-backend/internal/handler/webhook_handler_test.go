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

func newTestWebhookHandler() (*WebhookHandler, *MockWebhookRepository) {
	repo := new(MockWebhookRepository)
	return NewWebhookHandler(repo), repo
}

func TestWebhookHandler_CreateSubscription(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}
	subID := uuid.New()

	t.Run("success https", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		saved := &model.WebhookSubscription{
			ID:        subID,
			UserID:    userID,
			EventType: "inventory.changed",
			WebhookURL: "https://example.com/hook",
			SecretKey: "whsec_test",
			IsActive:  true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		repo.On("CreateSubscription", mock.Anything, mock.Anything).Return(saved, nil).Run(func(args mock.Arguments) {
			sub := args.Get(1).(*model.WebhookSubscription)
			assert.Equal(t, "inventory.changed", sub.EventType)
			assert.True(t, sub.IsActive)
			assert.Equal(t, 5, sub.MaxRetries)  // default
			assert.Equal(t, 60, sub.RetryInterval) // default
		})

		body := `{"event_type":"inventory.changed","webhook_url":"https://example.com/hook"}`
		req := requestWithUser(http.MethodPost, "/api/v1/webhooks/subscribe", body, user)
		rr := httptest.NewRecorder()
		h.CreateSubscription(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		resp := decodeResponse(t, rr)
		assert.Equal(t, "Subscription created", resp.Message)
		repo.AssertExpectations(t)
	})

	t.Run("localhost url allowed", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		saved := &model.WebhookSubscription{
			ID: subID, UserID: userID, EventType: "anomaly.detected",
			WebhookURL: "http://localhost:8080/hook", SecretKey: "whsec_x",
			IsActive: true, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
		}
		repo.On("CreateSubscription", mock.Anything, mock.Anything).Return(saved, nil)

		body := `{"event_type":"anomaly.detected","webhook_url":"http://localhost:8080/hook"}`
		req := requestWithUser(http.MethodPost, "/api/v1/webhooks/subscribe", body, user)
		rr := httptest.NewRecorder()
		h.CreateSubscription(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		body := `{"event_type":"inventory.changed","webhook_url":"https://example.com/hook"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/subscribe", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.CreateSubscription(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid event type", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		body := `{"event_type":"bogus.event","webhook_url":"https://example.com/hook"}`
		req := requestWithUser(http.MethodPost, "/api/v1/webhooks/subscribe", body, user)
		rr := httptest.NewRecorder()
		h.CreateSubscription(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid url scheme", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		body := `{"event_type":"inventory.changed","webhook_url":"ftp://example.com/hook"}`
		req := requestWithUser(http.MethodPost, "/api/v1/webhooks/subscribe", body, user)
		rr := httptest.NewRecorder()
		h.CreateSubscription(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing event type", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		body := `{"event_type":"","webhook_url":"https://example.com/hook"}`
		req := requestWithUser(http.MethodPost, "/api/v1/webhooks/subscribe", body, user)
		rr := httptest.NewRecorder()
		h.CreateSubscription(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("create error", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		repo.On("CreateSubscription", mock.Anything, mock.Anything).Return((*model.WebhookSubscription)(nil), errors.New("db error"))

		body := `{"event_type":"inventory.changed","webhook_url":"https://example.com/hook"}`
		req := requestWithUser(http.MethodPost, "/api/v1/webhooks/subscribe", body, user)
		rr := httptest.NewRecorder()
		h.CreateSubscription(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestWebhookHandler_ListSubscriptions(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}

	t.Run("success", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		subs := []model.WebhookSubscription{
			{ID: uuid.New(), UserID: userID, EventType: "inventory.changed", WebhookURL: "https://example.com/hook", IsActive: true, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		}
		repo.On("GetSubscriptions", mock.Anything, userID).Return(subs, nil)

		req := requestWithUser(http.MethodGet, "/api/v1/webhooks/subscriptions", "", user)
		rr := httptest.NewRecorder()
		h.ListSubscriptions(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		items := data["subscriptions"].([]interface{})
		assert.Equal(t, 1, len(items))
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/webhooks/subscriptions", nil)
		rr := httptest.NewRecorder()
		h.ListSubscriptions(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("error", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		repo.On("GetSubscriptions", mock.Anything, userID).Return([]model.WebhookSubscription{}, errors.New("db error"))

		req := requestWithUser(http.MethodGet, "/api/v1/webhooks/subscriptions", "", user)
		rr := httptest.NewRecorder()
		h.ListSubscriptions(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestWebhookHandler_DeleteSubscription(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}
	subID := uuid.New()
	otherSubID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		subs := []model.WebhookSubscription{{ID: subID, UserID: userID}}
		repo.On("GetSubscriptions", mock.Anything, userID).Return(subs, nil)
		repo.On("DeleteSubscription", mock.Anything, subID).Return(nil)

		req := requestWithUser(http.MethodDelete, "/api/v1/webhooks/subscriptions/"+subID.String(), "", user)
		rr := httptest.NewRecorder()
		h.DeleteSubscription(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/webhooks/subscriptions/"+subID.String(), nil)
		rr := httptest.NewRecorder()
		h.DeleteSubscription(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("not owner forbidden", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		subs := []model.WebhookSubscription{{ID: otherSubID, UserID: userID}}
		repo.On("GetSubscriptions", mock.Anything, userID).Return(subs, nil)

		req := requestWithUser(http.MethodDelete, "/api/v1/webhooks/subscriptions/"+subID.String(), "", user)
		rr := httptest.NewRecorder()
		h.DeleteSubscription(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		req := requestWithUser(http.MethodDelete, "/api/v1/webhooks/subscriptions/not-a-uuid", "", user)
		rr := httptest.NewRecorder()
		h.DeleteSubscription(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("delete error", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		subs := []model.WebhookSubscription{{ID: subID, UserID: userID}}
		repo.On("GetSubscriptions", mock.Anything, userID).Return(subs, nil)
		repo.On("DeleteSubscription", mock.Anything, subID).Return(errors.New("db error"))

		req := requestWithUser(http.MethodDelete, "/api/v1/webhooks/subscriptions/"+subID.String(), "", user)
		rr := httptest.NewRecorder()
		h.DeleteSubscription(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestWebhookHandler_ListDeliveries(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}
	subID := uuid.New()

	t.Run("empty without subscription filter", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		req := requestWithUser(http.MethodGet, "/api/v1/webhooks/deliveries", "", user)
		rr := httptest.NewRecorder()
		h.ListDeliveries(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(0), data["total"])
	})

	t.Run("success with subscription filter", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		subs := []model.WebhookSubscription{{ID: subID, UserID: userID}}
		deliveries := []model.WebhookDelivery{
			{ID: uuid.New(), SubscriptionID: subID, EventID: "evt-1", Status: model.DeliverySuccess, AttemptCount: 1, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		}
		repo.On("GetSubscriptions", mock.Anything, userID).Return(subs, nil)
		repo.On("GetDeliveryHistory", mock.Anything, subID, 20, 0).Return(deliveries, int64(1), nil)

		req := requestWithUser(http.MethodGet, "/api/v1/webhooks/deliveries?subscription_id="+subID.String(), "", user)
		rr := httptest.NewRecorder()
		h.ListDeliveries(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		items := data["deliveries"].([]interface{})
		assert.Equal(t, 1, len(items))
		repo.AssertExpectations(t)
	})

	t.Run("not owner forbidden", func(t *testing.T) {
		h, repo := newTestWebhookHandler()
		repo.On("GetSubscriptions", mock.Anything, userID).Return([]model.WebhookSubscription{}, nil)

		req := requestWithUser(http.MethodGet, "/api/v1/webhooks/deliveries?subscription_id="+subID.String(), "", user)
		rr := httptest.NewRecorder()
		h.ListDeliveries(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("invalid subscription id", func(t *testing.T) {
		h, _ := newTestWebhookHandler()
		req := requestWithUser(http.MethodGet, "/api/v1/webhooks/deliveries?subscription_id=not-a-uuid", "", user)
		rr := httptest.NewRecorder()
		h.ListDeliveries(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestWebhookHandler_New(t *testing.T) {
	repo := new(MockWebhookRepository)
	h := NewWebhookHandler(repo)
	assert.NotNil(t, h)
	assert.Equal(t, repo, h.repo)
}
