package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"dental_clinic/internal/modules/inventory/dto"
	"dental_clinic/internal/modules/inventory/models"
	"dental_clinic/internal/modules/inventory/services"

	"github.com/jackc/pgx/v5"

	"github.com/gorilla/mux"
)

type InventoryHandler struct {
	service *services.InventoryService
}

func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// CreateProduct godoc
// @Summary Create product
// @Description Creates a new inventory product
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ProductRequest true "Product data"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Router /api/products [post]
func (h *InventoryHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	product, err := h.service.CreateProduct(req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toProductResponse(*product))
}

// GetProducts godoc
// @Summary Get products
// @Description Returns all inventory products
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {object} map[string]string
// @Router /api/products [get]
func (h *InventoryHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetProducts()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toProductResponseList(products))
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Returns one inventory product by UUID
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
func (h *InventoryHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	product, err := h.service.GetProductByID(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toProductResponse(*product))
}

// UpdateProduct godoc
// @Summary Update product
// @Description Updates an inventory product by UUID
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.ProductRequest true "Product data"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Router /api/products/{id} [put]
func (h *InventoryHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	product, err := h.service.UpdateProduct(mux.Vars(r)["id"], req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toProductResponse(*product))
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Deletes an inventory product by UUID
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ActionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [delete]
func (h *InventoryHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteProduct(mux.Vars(r)["id"])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "product not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, dto.ActionResponse{Success: "1", Message: "product deleted"})
}

// AddStock godoc
// @Summary Add stock
// @Description Adds product stock to a clinic address inventory and records a restocked transaction
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Clinic address ID"
// @Param request body dto.InventoryQuantityRequest true "Stock data"
// @Success 200 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinic-addresses/{id}/inventory [post]
func (h *InventoryHandler) AddStock(w http.ResponseWriter, r *http.Request) {
	var req dto.InventoryQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	item, err := h.service.AddStock(r.Context(), mux.Vars(r)["id"], req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toInventoryResponse(*item))
}

// GetInventory godoc
// @Summary Get clinic address inventory
// @Description Returns inventory for a clinic address
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic address ID"
// @Success 200 {array} dto.InventoryResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinic-addresses/{id}/inventory [get]
func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	inventory, err := h.service.GetInventory(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toInventoryResponseList(inventory))
}

// UpdateInventory godoc
// @Summary Update inventory quantity
// @Description Sets inventory quantity for a clinic address inventory item and records a manual adjustment transaction
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Clinic address ID"
// @Param inventoryId path string true "Inventory item ID"
// @Param request body dto.UpdateInventoryRequest true "Inventory quantity"
// @Success 200 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinic-addresses/{id}/inventory/{inventoryId} [put]
func (h *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	vars := mux.Vars(r)
	item, err := h.service.UpdateInventory(r.Context(), vars["id"], vars["inventoryId"], req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toInventoryResponse(*item))
}

// AttachMaterial godoc
// @Summary Attach material to clinic service
// @Description Attaches required product material to a clinic service
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Clinic service ID"
// @Param request body dto.AttachMaterialRequest true "Service material data"
// @Success 201 {object} dto.ServiceMaterialResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinic-services/{id}/materials [post]
func (h *InventoryHandler) AttachMaterial(w http.ResponseWriter, r *http.Request) {
	var req dto.AttachMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	material, err := h.service.AttachMaterial(mux.Vars(r)["id"], req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toServiceMaterialResponse(*material))
}

// GetServiceMaterials godoc
// @Summary Get clinic service materials
// @Description Returns required product materials attached to a clinic service
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic service ID"
// @Success 200 {array} dto.ServiceMaterialResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinic-services/{id}/materials [get]
func (h *InventoryHandler) GetServiceMaterials(w http.ResponseWriter, r *http.Request) {
	materials, err := h.service.GetServiceMaterials(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toServiceMaterialResponseList(materials))
}

// GetTransactions godoc
// @Summary Get inventory transactions
// @Description Returns inventory transactions for a clinic address. Optional transaction_type values: restocked, used, manual_adjustment.
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param id path string true "Clinic address ID"
// @Param transaction_type query string false "Transaction type"
// @Success 200 {array} dto.InventoryTransactionResponse
// @Failure 400 {object} map[string]string
// @Router /api/clinic-addresses/{id}/inventory-transactions [get]
func (h *InventoryHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.service.GetTransactions(mux.Vars(r)["id"], r.URL.Query().Get("transaction_type"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toTransactionResponseList(transactions))
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func toProductResponse(product models.Product) dto.ProductResponse {
	return dto.ProductResponse{
		Id:        product.Id.String(),
		Name:      product.Name,
		Unit:      product.Unit,
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toProductResponseList(products []models.Product) []dto.ProductResponse {
	result := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		result = append(result, toProductResponse(product))
	}
	return result
}

func toInventoryResponse(item models.AddressInventory) dto.InventoryResponse {
	return dto.InventoryResponse{
		Id:              item.Id.String(),
		ClinicAddressId: item.ClinicAddressId.String(),
		ProductId:       item.ProductId.String(),
		ProductName:     item.ProductName,
		ProductUnit:     item.ProductUnit,
		Quantity:        item.Quantity,
		UpdatedAt:       item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toInventoryResponseList(inventory []models.AddressInventory) []dto.InventoryResponse {
	result := make([]dto.InventoryResponse, 0, len(inventory))
	for _, item := range inventory {
		result = append(result, toInventoryResponse(item))
	}
	return result
}

func toServiceMaterialResponse(material models.ServiceMaterial) dto.ServiceMaterialResponse {
	return dto.ServiceMaterialResponse{
		Id:               material.Id.String(),
		ClinicServiceId:  material.ClinicServiceId.String(),
		ProductId:        material.ProductId.String(),
		ProductName:      material.ProductName,
		ProductUnit:      material.ProductUnit,
		QuantityRequired: material.QuantityRequired,
	}
}

func toServiceMaterialResponseList(materials []models.ServiceMaterial) []dto.ServiceMaterialResponse {
	result := make([]dto.ServiceMaterialResponse, 0, len(materials))
	for _, material := range materials {
		result = append(result, toServiceMaterialResponse(material))
	}
	return result
}

func toTransactionResponse(transaction models.InventoryTransaction) dto.InventoryTransactionResponse {
	response := dto.InventoryTransactionResponse{
		Id:              transaction.Id.String(),
		ClinicAddressId: transaction.ClinicAddressId.String(),
		ProductId:       transaction.ProductId.String(),
		ProductName:     transaction.ProductName,
		Quantity:        transaction.Quantity,
		TransactionType: transaction.TransactionType,
		CreatedAt:       transaction.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if transaction.AppointmentId.String() != "00000000-0000-0000-0000-000000000000" {
		response.AppointmentId = transaction.AppointmentId.String()
	}
	return response
}

func toTransactionResponseList(transactions []models.InventoryTransaction) []dto.InventoryTransactionResponse {
	result := make([]dto.InventoryTransactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		result = append(result, toTransactionResponse(transaction))
	}
	return result
}
