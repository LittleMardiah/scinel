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

func newTestAPIKeyHandler() (*APIKeyHandler, *MockAPIKeyRepository) {
	repo := new(MockAPIKeyRepository)
	return NewAPIKeyHandler(repo), repo
}

func TestAPIKeyHandler_CreateAPIKey(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}
	keyID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		saved := &model.APIKey{
			ID:          keyID,
			UserID:      userID,
			KeyPrefix:   "sk_test_abcdef",
			KeyHash:     "hash",
			Name:        "test key",
			Permissions: []string{"inventory:read"},
			IsActive:    true,
			CreatedAt:   time.Now().UTC(),
		}
		repo.On("Create", mock.Anything, mock.Anything).Return(saved, nil).Run(func(args mock.Arguments) {
			key := args.Get(1).(*model.APIKey)
			assert.Equal(t, userID, key.UserID)
			assert.Equal(t, "test key", key.Name)
			assert.Equal(t, []string{"inventory:read"}, key.Permissions)
			assert.True(t, key.IsActive)
		})

		body := `{"name":"test key","permissions":["inventory:read"],"expires_in_days":30}`
		req := requestWithUser(http.MethodPost, "/api/v1/api-keys", body, user)
		rr := httptest.NewRecorder()
		h.CreateAPIKey(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		resp := decodeResponse(t, rr)
		assert.Equal(t, "API key created", resp.Message)
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		body := `{"name":"key","permissions":["inventory:read"]}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/api-keys", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.CreateAPIKey(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("missing name", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		body := `{"name":"","permissions":["inventory:read"]}`
		req := requestWithUser(http.MethodPost, "/api/v1/api-keys", body, user)
		rr := httptest.NewRecorder()
		h.CreateAPIKey(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing permissions", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		body := `{"name":"key","permissions":[]}`
		req := requestWithUser(http.MethodPost, "/api/v1/api-keys", body, user)
		rr := httptest.NewRecorder()
		h.CreateAPIKey(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid permission", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		body := `{"name":"key","permissions":["invalid:perm"]}`
		req := requestWithUser(http.MethodPost, "/api/v1/api-keys", body, user)
		rr := httptest.NewRecorder()
		h.CreateAPIKey(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("create error", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		repo.On("Create", mock.Anything, mock.Anything).Return((*model.APIKey)(nil), errors.New("db error"))

		body := `{"name":"key","permissions":["inventory:read"]}`
		req := requestWithUser(http.MethodPost, "/api/v1/api-keys", body, user)
		rr := httptest.NewRecorder()
		h.CreateAPIKey(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestAPIKeyHandler_ListAPIKeys(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}

	t.Run("success", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		keys := []model.APIKey{
			{ID: uuid.New(), UserID: userID, KeyPrefix: "sk_test_x", Name: "k1", Permissions: []string{"inventory:read"}, IsActive: true, CreatedAt: time.Now().UTC()},
		}
		repo.On("GetByUser", mock.Anything, userID).Return(keys, nil)

		req := requestWithUser(http.MethodGet, "/api/v1/api-keys", "", user)
		rr := httptest.NewRecorder()
		h.ListAPIKeys(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		items := data["keys"].([]interface{})
		assert.Equal(t, 1, len(items))
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/api-keys", nil)
		rr := httptest.NewRecorder()
		h.ListAPIKeys(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("error", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		repo.On("GetByUser", mock.Anything, userID).Return([]model.APIKey{}, errors.New("db error"))

		req := requestWithUser(http.MethodGet, "/api/v1/api-keys", "", user)
		rr := httptest.NewRecorder()
		h.ListAPIKeys(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestAPIKeyHandler_RevokeAPIKey(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}
	keyID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		repo.On("Delete", mock.Anything, keyID, userID).Return(nil)

		req := requestWithUser(http.MethodDelete, "/api/v1/api-keys/"+keyID.String(), "", user)
		rr := httptest.NewRecorder()
		h.RevokeAPIKey(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/api-keys/"+keyID.String(), nil)
		rr := httptest.NewRecorder()
		h.RevokeAPIKey(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid key id", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		req := requestWithUser(http.MethodDelete, "/api/v1/api-keys/not-a-uuid", "", user)
		rr := httptest.NewRecorder()
		h.RevokeAPIKey(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("delete error", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		repo.On("Delete", mock.Anything, keyID, userID).Return(errors.New("db error"))

		req := requestWithUser(http.MethodDelete, "/api/v1/api-keys/"+keyID.String(), "", user)
		rr := httptest.NewRecorder()
		h.RevokeAPIKey(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestAPIKeyHandler_RotateAPIKey(t *testing.T) {
	userID := uuid.New()
	user := &middleware.UserInfo{UserID: userID.String(), Email: "u@example.com", Role: "admin"}
	keyID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		existing := &model.APIKey{
			ID:          keyID,
			UserID:      userID,
			KeyPrefix:   "sk_test_abcdef",
			KeyHash:     "oldhash",
			Name:        "k1",
			Permissions: []string{"inventory:read"},
			IsActive:    true,
			UpdatedAt:   time.Now().UTC(),
		}
		repo.On("GetByID", mock.Anything, keyID, userID).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			key := args.Get(1).(*model.APIKey)
			assert.NotEqual(t, "oldhash", key.KeyHash)
		})

		req := requestWithUser(http.MethodPost, "/api/v1/api-keys/"+keyID.String()+"/rotate", "", user)
		rr := httptest.NewRecorder()
		h.RotateAPIKey(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/api-keys/"+keyID.String()+"/rotate", nil)
		rr := httptest.NewRecorder()
		h.RotateAPIKey(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid path", func(t *testing.T) {
		h, _ := newTestAPIKeyHandler()
		req := requestWithUser(http.MethodPost, "/api/v1/api-keys/"+keyID.String()+"/wrong", "", user)
		rr := httptest.NewRecorder()
		h.RotateAPIKey(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("not found", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		repo.On("GetByID", mock.Anything, keyID, userID).Return((*model.APIKey)(nil), errors.New("not found"))

		req := requestWithUser(http.MethodPost, "/api/v1/api-keys/"+keyID.String()+"/rotate", "", user)
		rr := httptest.NewRecorder()
		h.RotateAPIKey(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("inactive key conflict", func(t *testing.T) {
		h, repo := newTestAPIKeyHandler()
		existing := &model.APIKey{ID: keyID, UserID: userID, KeyPrefix: "sk_test_x", IsActive: false}
		repo.On("GetByID", mock.Anything, keyID, userID).Return(existing, nil)

		req := requestWithUser(http.MethodPost, "/api/v1/api-keys/"+keyID.String()+"/rotate", "", user)
		rr := httptest.NewRecorder()
		h.RotateAPIKey(rr, req)
		assert.Equal(t, http.StatusConflict, rr.Code)
	})
}

func TestAPIKeyHandler_New(t *testing.T) {
	repo := new(MockAPIKeyRepository)
	h := NewAPIKeyHandler(repo)
	assert.NotNil(t, h)
	assert.Equal(t, repo, h.repo)
}
