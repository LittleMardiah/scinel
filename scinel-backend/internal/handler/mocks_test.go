package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/LittleMardiah/scinel/internal/model"
)

// ---- UserRepository mock ----

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, email, passwordHash, fullName, role string) (*model.User, error) {
	args := m.Called(ctx, email, passwordHash, fullName, role)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error {
	args := m.Called(ctx, id, role)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) CountUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// ---- InventoryRepository mock ----

type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) GetBySKU(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
	args := m.Called(ctx, warehouseID, skuID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Inventory, int64, error) {
	args := m.Called(ctx, warehouseID, limit, offset)
	return args.Get(0).([]model.Inventory), args.Get(1).(int64), args.Error(2)
}

func (m *MockInventoryRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Inventory, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Inventory), args.Get(1).(int64), args.Error(2)
}

func (m *MockInventoryRepository) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	args := m.Called(ctx, inventory)
	return args.Get(0).(*model.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) UpdateQuantity(ctx context.Context, warehouseID, skuID uuid.UUID, delta int32) (*model.Inventory, error) {
	args := m.Called(ctx, warehouseID, skuID, delta)
	return args.Get(0).(*model.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) LockForUpdate(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
	args := m.Called(ctx, warehouseID, skuID)
	return args.Get(0).(*model.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) GetBySKUCached(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
	args := m.Called(ctx, warehouseID, skuID)
	return args.Get(0).(*model.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) InvalidateCache(ctx context.Context, warehouseID, skuID uuid.UUID) error {
	args := m.Called(ctx, warehouseID, skuID)
	return args.Error(0)
}

func (m *MockInventoryRepository) CreateSnapshot(ctx context.Context, snapshot *model.InventorySnapshot) error {
	args := m.Called(ctx, snapshot)
	return args.Error(0)
}

// ---- TransactionRepository mock ----

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, txn *model.Transaction) (*model.Transaction, error) {
	args := m.Called(ctx, txn)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, txnID uuid.UUID) (*model.Transaction, error) {
	args := m.Called(ctx, txnID)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Query(ctx context.Context, filter model.TransactionFilter, limit, offset int) ([]model.Transaction, int64, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]model.Transaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionRepository) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error) {
	args := m.Called(ctx, warehouseID, limit, offset)
	return args.Get(0).([]model.Transaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionRepository) GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error) {
	args := m.Called(ctx, skuID, limit, offset)
	return args.Get(0).([]model.Transaction), args.Get(1).(int64), args.Error(2)
}

// ---- AnomalyRepository mock ----

type MockAnomalyRepository struct {
	mock.Mock
}

func (m *MockAnomalyRepository) Create(ctx context.Context, anomaly *model.Anomaly) (*model.Anomaly, error) {
	args := m.Called(ctx, anomaly)
	return args.Get(0).(*model.Anomaly), args.Error(1)
}

func (m *MockAnomalyRepository) GetByID(ctx context.Context, anomalyID uuid.UUID) (*model.Anomaly, error) {
	args := m.Called(ctx, anomalyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Anomaly), args.Error(1)
}

func (m *MockAnomalyRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Anomaly), args.Get(1).(int64), args.Error(2)
}

func (m *MockAnomalyRepository) GetOpen(ctx context.Context, limit, offset int) ([]model.Anomaly, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Anomaly), args.Get(1).(int64), args.Error(2)
}

func (m *MockAnomalyRepository) UpdateStatus(ctx context.Context, anomalyID uuid.UUID, status model.AnomalyStatus, notes string, resolvedBy uuid.UUID) error {
	args := m.Called(ctx, anomalyID, status, notes, resolvedBy)
	return args.Error(0)
}

func (m *MockAnomalyRepository) GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Anomaly, int64, error) {
	args := m.Called(ctx, skuID, limit, offset)
	return args.Get(0).([]model.Anomaly), args.Get(1).(int64), args.Error(2)
}

// ---- APIKeyRepository mock ----

type MockAPIKeyRepository struct {
	mock.Mock
}

func (m *MockAPIKeyRepository) Create(ctx context.Context, key *model.APIKey) (*model.APIKey, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*model.APIKey), args.Error(1)
}

func (m *MockAPIKeyRepository) GetByID(ctx context.Context, keyID, userID uuid.UUID) (*model.APIKey, error) {
	args := m.Called(ctx, keyID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APIKey), args.Error(1)
}

func (m *MockAPIKeyRepository) GetByPrefix(ctx context.Context, prefix string) (*model.APIKey, error) {
	args := m.Called(ctx, prefix)
	return args.Get(0).(*model.APIKey), args.Error(1)
}

func (m *MockAPIKeyRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.APIKey), args.Error(1)
}

func (m *MockAPIKeyRepository) Update(ctx context.Context, key *model.APIKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockAPIKeyRepository) Delete(ctx context.Context, keyID, userID uuid.UUID) error {
	args := m.Called(ctx, keyID, userID)
	return args.Error(0)
}

// ---- WebhookRepository mock ----

type MockWebhookRepository struct {
	mock.Mock
}

func (m *MockWebhookRepository) CreateSubscription(ctx context.Context, sub *model.WebhookSubscription) (*model.WebhookSubscription, error) {
	args := m.Called(ctx, sub)
	return args.Get(0).(*model.WebhookSubscription), args.Error(1)
}

func (m *MockWebhookRepository) GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.WebhookSubscription, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.WebhookSubscription), args.Error(1)
}

func (m *MockWebhookRepository) GetSubscriptionsByEventType(ctx context.Context, eventType string) ([]model.WebhookSubscription, error) {
	args := m.Called(ctx, eventType)
	return args.Get(0).([]model.WebhookSubscription), args.Error(1)
}

func (m *MockWebhookRepository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*model.WebhookSubscription, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.WebhookSubscription), args.Error(1)
}

func (m *MockWebhookRepository) UpdateSubscription(ctx context.Context, sub *model.WebhookSubscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockWebhookRepository) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	args := m.Called(ctx, subscriptionID)
	return args.Error(0)
}

func (m *MockWebhookRepository) CreateDelivery(ctx context.Context, delivery *model.WebhookDelivery) (*model.WebhookDelivery, error) {
	args := m.Called(ctx, delivery)
	return args.Get(0).(*model.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookRepository) GetPendingDeliveries(ctx context.Context, limit int) ([]model.WebhookDelivery, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]model.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookRepository) UpdateDeliveryStatus(ctx context.Context, deliveryID uuid.UUID, status model.DeliveryStatus, httpStatus *int, errMsg *string) error {
	args := m.Called(ctx, deliveryID, status, httpStatus, errMsg)
	return args.Error(0)
}

func (m *MockWebhookRepository) GetDeliveryHistory(ctx context.Context, subscriptionID uuid.UUID, limit, offset int) ([]model.WebhookDelivery, int64, error) {
	args := m.Called(ctx, subscriptionID, limit, offset)
	return args.Get(0).([]model.WebhookDelivery), args.Get(1).(int64), args.Error(2)
}
