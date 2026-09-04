package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/LittleMardiah/scinel/internal/database"
	"github.com/LittleMardiah/scinel/internal/middleware"
	"github.com/LittleMardiah/scinel/internal/model"
)

func newTestInventoryHandler(t *testing.T) (*InventoryHandler, *MockInventoryRepository, *MockTransactionRepository, *database.DB, sqlmock.Sqlmock) {
	t.Helper()
	invRepo := new(MockInventoryRepository)
	txnRepo := new(MockTransactionRepository)

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })
	db := &database.DB{DB: sqlDB}

	return NewInventoryHandler(invRepo, txnRepo, db), invRepo, txnRepo, db, mock
}

func requestWithUser(method, target, body string, user *middleware.UserInfo) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	if user != nil {
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
	}
	return req
}

func TestInventoryHandler_CreateTransaction(t *testing.T) {
	now := time.Now().UTC()
	whID := uuid.New()
	skuID := uuid.New()
	userUUID := uuid.New()

	user := &middleware.UserInfo{UserID: userUUID.String(), Email: "u@example.com", Role: "viewer"}

	currentInv := &model.Inventory{
		ID:          uuid.New(),
		WarehouseID: whID,
		SKUID:       skuID,
		SKUCode:     "SKU-001",
		Quantity:    50,
	}

	t.Run("create in transaction success", func(t *testing.T) {
		h, invRepo, txnRepo, _, _ := newTestInventoryHandler(t)

		updatedInv := *currentInv
		updatedInv.Quantity = 60

		savedTxn := &model.Transaction{
			ID:               uuid.New(),
			WarehouseID:      whID,
			SKUID:            skuID,
			SKUCode:          "SKU-001",
			Type:             model.TransactionIn,
			Quantity:         10,
			PreviousQuantity: 50,
			NewQuantity:      60,
			Reason:           "restock",
			CreatedAt:        now,
		}

		invRepo.On("GetBySKU", mock.Anything, whID, skuID).Return(currentInv, nil)
		invRepo.On("UpdateQuantity", mock.Anything, whID, skuID, int32(10)).Return(&updatedInv, nil)
		txnRepo.On("Create", mock.Anything, mock.Anything).Return(savedTxn, nil).Run(func(args mock.Arguments) {
			txn := args.Get(1).(*model.Transaction)
			assert.Equal(t, model.TransactionIn, txn.Type)
			assert.Equal(t, int32(10), txn.Quantity)
			assert.Equal(t, int32(50), txn.PreviousQuantity)
			assert.Equal(t, int32(60), txn.NewQuantity)
		})
		invRepo.On("InvalidateCache", mock.Anything, whID, skuID).Return(nil)

		body := `{"warehouse_id":"` + whID.String() + `","sku_id":"` + skuID.String() + `","type":"in","quantity":10,"reason":"restock"}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		resp := decodeResponse(t, rr)
		assert.Equal(t, "Transaction created", resp.Message)
		time.Sleep(20 * time.Millisecond)
		invRepo.AssertExpectations(t)
		txnRepo.AssertExpectations(t)
	})

	t.Run("unauthorized without user context", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		body := `{"warehouse_id":"` + whID.String() + `","sku_id":"` + skuID.String() + `","type":"in","quantity":10,"reason":"restock"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", "{invalid", user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		body := `{"warehouse_id":"","type":"","quantity":0}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid type", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		body := `{"warehouse_id":"` + whID.String() + `","sku_id":"` + skuID.String() + `","type":"sideways","quantity":10,"reason":"x"}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid warehouse id", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		body := `{"warehouse_id":"not-a-uuid","sku_id":"` + skuID.String() + `","type":"in","quantity":10,"reason":"x"}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("insufficient stock for out transaction", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		lowInv := *currentInv
		lowInv.Quantity = 5
		invRepo.On("GetBySKU", mock.Anything, whID, skuID).Return(&lowInv, nil)

		body := `{"warehouse_id":"` + whID.String() + `","sku_id":"` + skuID.String() + `","type":"out","quantity":10,"reason":"sale"}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("inventory not found", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetBySKU", mock.Anything, whID, skuID).Return((*model.Inventory)(nil), nil)

		body := `{"warehouse_id":"` + whID.String() + `","sku_id":"` + skuID.String() + `","type":"in","quantity":10,"reason":"x"}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("get current stock error", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetBySKU", mock.Anything, whID, skuID).Return((*model.Inventory)(nil), errors.New("db error"))

		body := `{"warehouse_id":"` + whID.String() + `","sku_id":"` + skuID.String() + `","type":"in","quantity":10,"reason":"x"}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("sku code lookup via db", func(t *testing.T) {
		h, invRepo, txnRepo, _, mockDB := newTestInventoryHandler(t)

		updatedInv := *currentInv
		updatedInv.Quantity = 60

		savedTxn := &model.Transaction{
			ID:               uuid.New(),
			WarehouseID:      whID,
			SKUID:            skuID,
			SKUCode:          "SKU-001",
			Type:             model.TransactionIn,
			Quantity:         10,
			PreviousQuantity: 50,
			NewQuantity:      60,
			Reason:           "restock",
			CreatedAt:        now,
		}

		mockDB.ExpectQuery("SELECT id FROM skus WHERE sku_code = \\$1").
			WithArgs("SKU-001").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(skuID.String()))

		invRepo.On("GetBySKU", mock.Anything, whID, skuID).Return(currentInv, nil)
		invRepo.On("UpdateQuantity", mock.Anything, whID, skuID, int32(10)).Return(&updatedInv, nil)
		txnRepo.On("Create", mock.Anything, mock.Anything).Return(savedTxn, nil)
		invRepo.On("InvalidateCache", mock.Anything, whID, skuID).Return(nil)

		body := `{"warehouse_id":"` + whID.String() + `","sku_code":"SKU-001","type":"in","quantity":10,"reason":"restock"}`
		req := requestWithUser(http.MethodPost, "/api/v1/transactions", body, user)
		rr := httptest.NewRecorder()
		h.CreateTransaction(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		time.Sleep(20 * time.Millisecond)
		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("db expectations not met: %v", err)
		}
		invRepo.AssertExpectations(t)
		txnRepo.AssertExpectations(t)
	})
}

func TestInventoryHandler_QueryTransactions(t *testing.T) {
	whID := uuid.New()
	skuID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h, _, txnRepo, _, _ := newTestInventoryHandler(t)
		txns := []model.Transaction{
			{
				ID:               uuid.New(),
				WarehouseID:      whID,
				SKUID:            skuID,
				SKUCode:          "SKU-001",
				Type:             model.TransactionIn,
				Quantity:         5,
				PreviousQuantity: 10,
				NewQuantity:      15,
				Reason:           "restock",
				CreatedAt:        time.Now().UTC(),
			},
		}
		txnRepo.On("Query", mock.Anything, mock.Anything, 20, 0).Return(txns, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
		rr := httptest.NewRecorder()
		h.QueryTransactions(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		resp := decodeResponse(t, rr)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(1), data["total"])
		txnRepo.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		h, _, txnRepo, _, _ := newTestInventoryHandler(t)
		txnRepo.On("Query", mock.Anything, mock.Anything, 20, 0).Return([]model.Transaction(nil), int64(0), errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
		rr := httptest.NewRecorder()
		h.QueryTransactions(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("with filters", func(t *testing.T) {
		h, _, txnRepo, _, _ := newTestInventoryHandler(t)
		txnRepo.On("Query", mock.Anything, mock.Anything, 100, 10).Return([]model.Transaction{}, int64(0), nil)

		url := "/api/v1/transactions?warehouse_id=" + whID.String() + "&sku_id=" + skuID.String() + "&type=in&limit=200&offset=10"
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		h.QueryTransactions(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		txnRepo.AssertExpectations(t)
	})
}

func TestInventoryHandler_QueryInventory(t *testing.T) {
	whID := uuid.New()
	skuID := uuid.New()

	inv := model.Inventory{
		ID:            uuid.New(),
		WarehouseID:   whID,
		SKUID:         skuID,
		SKUCode:       "SKU-001",
		Quantity:      100,
		ReorderPoint:  20,
		MaxStockLevel: 200,
		UpdatedAt:     time.Now().UTC(),
	}

	t.Run("by warehouse", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetByWarehouse", mock.Anything, whID, 20, 0).Return([]model.Inventory{inv}, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory?warehouse_id="+whID.String(), nil)
		rr := httptest.NewRecorder()
		h.QueryInventory(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		invRepo.AssertExpectations(t)
	})

	t.Run("by sku", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetAll", mock.Anything, 1000, 0).Return([]model.Inventory{inv}, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory?sku_id="+skuID.String(), nil)
		rr := httptest.NewRecorder()
		h.QueryInventory(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		invRepo.AssertExpectations(t)
	})

	t.Run("all inventory", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetAll", mock.Anything, 20, 0).Return([]model.Inventory{inv}, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory", nil)
		rr := httptest.NewRecorder()
		h.QueryInventory(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		invRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetAll", mock.Anything, 20, 0).Return([]model.Inventory{}, int64(0), errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory", nil)
		rr := httptest.NewRecorder()
		h.QueryInventory(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestInventoryHandler_GetWarehouseInventory(t *testing.T) {
	whID := uuid.New()
	user := &middleware.UserInfo{UserID: uuid.New().String(), Email: "u@example.com", Role: "viewer"}

	t.Run("success", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetByWarehouse", mock.Anything, whID, 20, 0).Return([]model.Inventory{}, int64(0), nil)

		req := requestWithUser(http.MethodGet, "/api/v1/inventory/"+whID.String(), "", user)
		rr := httptest.NewRecorder()
		h.GetWarehouseInventory(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		invRepo.AssertExpectations(t)
	})

	t.Run("missing warehouse id", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		req := requestWithUser(http.MethodGet, "/api/v1/inventory/", "", user)
		rr := httptest.NewRecorder()
		h.GetWarehouseInventory(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid warehouse id", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		req := requestWithUser(http.MethodGet, "/api/v1/inventory/not-a-uuid", "", user)
		rr := httptest.NewRecorder()
		h.GetWarehouseInventory(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h, _, _, _, _ := newTestInventoryHandler(t)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/"+whID.String(), nil)
		rr := httptest.NewRecorder()
		h.GetWarehouseInventory(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("query error", func(t *testing.T) {
		h, invRepo, _, _, _ := newTestInventoryHandler(t)
		invRepo.On("GetByWarehouse", mock.Anything, whID, 20, 0).Return([]model.Inventory{}, int64(0), errors.New("db error"))
		req := requestWithUser(http.MethodGet, "/api/v1/inventory/"+whID.String(), "", user)
		rr := httptest.NewRecorder()
		h.GetWarehouseInventory(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestInventoryHandler_ImportInventory(t *testing.T) {
	h, _, _, _, _ := newTestInventoryHandler(t)

	body := bytes.NewBufferString("--X\r\nContent-Disposition: form-data; name=\"file\"; filename=\"data.csv\"\r\nContent-Type: text/csv\r\n\r\nsku,quantity\r\nSKU-001,10\r\n--X--\r\n")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/import", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=X")
	rr := httptest.NewRecorder()
	h.ImportInventory(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	resp := decodeResponse(t, rr)
	assert.Equal(t, "Import processed (placeholder)", resp.Message)
}

func TestInventoryHandler_New(t *testing.T) {
	invRepo := new(MockInventoryRepository)
	txnRepo := new(MockTransactionRepository)
	sqlDB, _, _ := sqlmock.New()
	defer sqlDB.Close()
	db := &database.DB{DB: sqlDB}
	h := NewInventoryHandler(invRepo, txnRepo, db)
	assert.NotNil(t, h)
	assert.Equal(t, invRepo, h.invRepo)
	assert.Equal(t, txnRepo, h.txnRepo)
	assert.Equal(t, db, h.db)
}
