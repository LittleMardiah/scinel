package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/middleware"
    "github.com/LittleMardiah/scinel/internal/model"
    "github.com/LittleMardiah/scinel/internal/repository"
    "github.com/LittleMardiah/scinel/internal/util"
)

// InventoryHandler handles inventory-related endpoints
type InventoryHandler struct {
    invRepo repository.InventoryRepository
    txnRepo repository.TransactionRepository
}

func NewInventoryHandler(invRepo repository.InventoryRepository, txnRepo repository.TransactionRepository) *InventoryHandler {
    return &InventoryHandler{invRepo: invRepo, txnRepo: txnRepo}
}

// ==================== REQUEST/RESPONSE TYPES ====================

type createTransactionRequest struct {
    WarehouseID string          `json:"warehouse_id"`
    SKUID       string          `json:"sku_id"`
    Type        string          `json:"type"`
    Quantity    int32           `json:"quantity"`
    Reason      string          `json:"reason"`
    Details     json.RawMessage `json:"details,omitempty"`
}

type transactionResponse struct {
    ID               string    `json:"id"`
    WarehouseID      string    `json:"warehouse_id"`
    SKUID            string    `json:"sku_id"`
    SKUCode          string    `json:"sku_code"`
    Type             string    `json:"type"`
    Quantity         int32     `json:"quantity"`
    PreviousQuantity int32     `json:"previous_quantity"`
    NewQuantity      int32     `json:"new_quantity"`
    Reason           string    `json:"reason"`
    CreatedAt        time.Time `json:"created_at"`
}

type inventoryResponse struct {
    ID             string     `json:"id"`
    WarehouseID    string     `json:"warehouse_id"`
    SKUID          string     `json:"sku_id"`
    SKUCode        string     `json:"sku_code"`
    Quantity       int32      `json:"quantity"`
    ReorderPoint   int32      `json:"reorder_point"`
    MaxStockLevel  int32      `json:"max_stock_level"`
    Status         string     `json:"status"`
    LastRestockDate *time.Time `json:"last_restock_date,omitempty"`
    UpdatedAt      time.Time  `json:"updated_at"`
}

// ==================== ENDPOINT 1: POST /api/v1/transactions ====================

func (h *InventoryHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", requestID)
        return
    }

    var req createTransactionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error(), requestID)
        return
    }

    if req.WarehouseID == "" || req.SKUID == "" || req.Type == "" || req.Quantity <= 0 {
        util.WriteError(w, http.StatusBadRequest, "Missing required fields", "warehouse_id, sku_id, type, quantity (positive) required", requestID)
        return
    }
    if req.Type != "in" && req.Type != "out" {
        util.WriteError(w, http.StatusBadRequest, "Invalid type", "type must be 'in' or 'out'", requestID)
        return
    }
    if req.Reason == "" {
        util.WriteError(w, http.StatusBadRequest, "Reason is required", "Please provide a reason", requestID)
        return
    }

    warehouseUUID, err := uuid.Parse(req.WarehouseID)
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid warehouse_id", err.Error(), requestID)
        return
    }
    skuUUID, err := uuid.Parse(req.SKUID)
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid sku_id", err.Error(), requestID)
        return
    }

    // RBAC: Manager only own warehouse (placeholder)
    if userInfo.Role == "manager" {
        // TODO: check warehouse ownership
    }

    currentInv, err := h.invRepo.GetBySKU(r.Context(), warehouseUUID, skuUUID)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to get current stock", err.Error(), requestID)
        return
    }
    if currentInv == nil {
        util.WriteError(w, http.StatusNotFound, "Inventory record not found", "", requestID)
        return
    }

    var delta int32
    if req.Type == "in" {
        delta = req.Quantity
    } else {
        if currentInv.Quantity < req.Quantity {
            util.WriteError(w, http.StatusConflict, "Insufficient stock", "Requested out quantity exceeds available stock", requestID)
            return
        }
        delta = -req.Quantity
    }

    updatedInv, err := h.invRepo.UpdateQuantity(r.Context(), warehouseUUID, skuUUID, delta)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to update stock", err.Error(), requestID)
        return
    }

    createdBy, _ := uuid.Parse(userInfo.UserID)
    txn := &model.Transaction{
        WarehouseID:      warehouseUUID,
        SKUID:            skuUUID,
        SKUCode:          currentInv.SKUCode,
        Type:             model.TransactionType(req.Type),
        Quantity:         req.Quantity,
        Reason:           req.Reason,
        PreviousQuantity: currentInv.Quantity,
        NewQuantity:      updatedInv.Quantity,
        CreatedBy:        createdBy,
        Details:          req.Details,
    }
    txn.CreatedAt = time.Now().UTC()

    savedTxn, err := h.txnRepo.Create(r.Context(), txn)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to record transaction", err.Error(), requestID)
        return
    }

    go func() {
        _ = h.invRepo.InvalidateCache(context.Background(), warehouseUUID, skuUUID)
    }()

    resp := transactionResponse{
        ID:               savedTxn.ID.String(),
        WarehouseID:      savedTxn.WarehouseID.String(),
        SKUID:            savedTxn.SKUID.String(),
        SKUCode:          savedTxn.SKUCode,
        Type:             string(savedTxn.Type),
        Quantity:         savedTxn.Quantity,
        PreviousQuantity: savedTxn.PreviousQuantity,
        NewQuantity:      savedTxn.NewQuantity,
        Reason:           savedTxn.Reason,
        CreatedAt:        savedTxn.CreatedAt,
    }
    util.WriteResponse(w, http.StatusCreated, "Transaction created", resp, requestID)
}

// ==================== ENDPOINT 2: GET /api/v1/transactions ====================

func (h *InventoryHandler) QueryTransactions(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()
    limit, offset := util.ParsePagination(r)

    filter := model.TransactionFilter{}
    if warehouseID := r.URL.Query().Get("warehouse_id"); warehouseID != "" {
        if id, err := uuid.Parse(warehouseID); err == nil {
            filter.WarehouseID = &id
        }
    }
    if skuID := r.URL.Query().Get("sku_id"); skuID != "" {
        if id, err := uuid.Parse(skuID); err == nil {
            filter.SKUID = &id
        }
    }
    if t := r.URL.Query().Get("type"); t != "" {
        typ := model.TransactionType(t)
        filter.Type = &typ
    }
    if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
        if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
            filter.DateFrom = &t
        }
    }
    if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
        if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
            filter.DateTo = &t
        }
    }

    if limit > 100 {
        limit = 100
    }

    txns, total, err := h.txnRepo.Query(r.Context(), filter, limit, offset)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to query transactions", err.Error(), requestID)
        return
    }

    items := make([]transactionResponse, len(txns))
    for i, t := range txns {
        items[i] = transactionResponse{
            ID:               t.ID.String(),
            WarehouseID:      t.WarehouseID.String(),
            SKUID:            t.SKUID.String(),
            SKUCode:          t.SKUCode,
            Type:             string(t.Type),
            Quantity:         t.Quantity,
            PreviousQuantity: t.PreviousQuantity,
            NewQuantity:      t.NewQuantity,
            Reason:           t.Reason,
            CreatedAt:        t.CreatedAt,
        }
    }
    data := map[string]interface{}{
        "transactions": items,
        "total":        total,
        "limit":        limit,
        "offset":       offset,
    }
    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== ENDPOINT 3: GET /api/v1/inventory ====================

func (h *InventoryHandler) QueryInventory(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()
    limit, offset := util.ParsePagination(r)

    warehouseIDStr := r.URL.Query().Get("warehouse_id")
    skuIDStr := r.URL.Query().Get("sku_id")
    statusStr := r.URL.Query().Get("status")

    var inventory []model.Inventory
    var total int64

    if warehouseIDStr != "" {
        warehouseID, _ := uuid.Parse(warehouseIDStr)
        inv, t, err := h.invRepo.GetByWarehouse(r.Context(), warehouseID, limit, offset)
        if err != nil {
            util.WriteError(w, http.StatusInternalServerError, "Failed to query inventory", err.Error(), requestID)
            return
        }
        inventory = inv
        total = t
    } else if skuIDStr != "" {
        inv, _, err := h.invRepo.GetAll(r.Context(), 1000, 0)
        if err != nil {
            util.WriteError(w, http.StatusInternalServerError, "Failed to query inventory", err.Error(), requestID)
            return
        }
        skuID, _ := uuid.Parse(skuIDStr)
        filtered := []model.Inventory{}
        for _, item := range inv {
            if item.SKUID == skuID {
                filtered = append(filtered, item)
            }
        }
        inventory = filtered
        total = int64(len(filtered))
    } else {
        inv, t, err := h.invRepo.GetAll(r.Context(), limit, offset)
        if err != nil {
            util.WriteError(w, http.StatusInternalServerError, "Failed to query inventory", err.Error(), requestID)
            return
        }
        inventory = inv
        total = t
    }

    if statusStr != "" {
        filtered := []model.Inventory{}
        for _, item := range inventory {
            if string(item.GetStatus()) == statusStr {
                filtered = append(filtered, item)
            }
        }
        inventory = filtered
        total = int64(len(filtered))
    }

    items := make([]inventoryResponse, len(inventory))
    for i, inv := range inventory {
        items[i] = inventoryResponse{
            ID:             inv.ID.String(),
            WarehouseID:    inv.WarehouseID.String(),
            SKUID:          inv.SKUID.String(),
            SKUCode:        inv.SKUCode,
            Quantity:       inv.Quantity,
            ReorderPoint:   inv.ReorderPoint,
            MaxStockLevel:  inv.MaxStockLevel,
            Status:         string(inv.GetStatus()),
            LastRestockDate: inv.LastRestockDate,
            UpdatedAt:      inv.UpdatedAt,
        }
    }
    data := map[string]interface{}{
        "inventory": items,
        "total":     total,
        "limit":     limit,
        "offset":    offset,
    }
    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== ENDPOINT 4: GET /api/v1/inventory/{warehouse_id} ====================

func (h *InventoryHandler) GetWarehouseInventory(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/inventory/")
    if path == "" {
        util.WriteError(w, http.StatusBadRequest, "Warehouse ID required", "", requestID)
        return
    }
    warehouseID, err := uuid.Parse(path)
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "Invalid warehouse_id", err.Error(), requestID)
        return
    }

    userInfo, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserInfo)
    if !ok {
        util.WriteError(w, http.StatusUnauthorized, "Unauthorized", "", requestID)
        return
    }
    if userInfo.Role == "manager" {
        // TODO: check warehouse ownership
    }

    limit, offset := util.ParsePagination(r)
    if limit > 100 {
        limit = 100
    }

    inventory, total, err := h.invRepo.GetByWarehouse(r.Context(), warehouseID, limit, offset)
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, "Failed to get warehouse inventory", err.Error(), requestID)
        return
    }

    items := make([]inventoryResponse, len(inventory))
    for i, inv := range inventory {
        items[i] = inventoryResponse{
            ID:             inv.ID.String(),
            WarehouseID:    inv.WarehouseID.String(),
            SKUID:          inv.SKUID.String(),
            SKUCode:        inv.SKUCode,
            Quantity:       inv.Quantity,
            ReorderPoint:   inv.ReorderPoint,
            MaxStockLevel:  inv.MaxStockLevel,
            Status:         string(inv.GetStatus()),
            LastRestockDate: inv.LastRestockDate,
            UpdatedAt:      inv.UpdatedAt,
        }
    }
    data := map[string]interface{}{
        "inventory": items,
        "total":     total,
        "limit":     limit,
        "offset":    offset,
    }
    util.WriteResponse(w, http.StatusOK, "OK", data, requestID)
}

// ==================== ENDPOINT 5: POST /api/v1/inventory/import ====================

func (h *InventoryHandler) ImportInventory(w http.ResponseWriter, r *http.Request) {
    requestID := util.GenerateRequestID()
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        util.WriteError(w, http.StatusBadRequest, "Failed to parse form", err.Error(), requestID)
        return
    }
    file, _, err := r.FormFile("file")
    if err != nil {
        util.WriteError(w, http.StatusBadRequest, "File is required", err.Error(), requestID)
        return
    }
    defer file.Close()

    data := map[string]interface{}{
        "imported": 0,
        "errors":   0,
        "message":  "Import functionality not yet implemented",
    }
    util.WriteResponse(w, http.StatusOK, "Import processed (placeholder)", data, requestID)
}