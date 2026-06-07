package services

import (
	"context"
	"errors"
	"time"

	"dental_clinic/internal/modules/inventory/dto"
	"dental_clinic/internal/modules/inventory/models"
	"dental_clinic/internal/modules/inventory/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService struct {
	repo repository.InventoryRepository
	db   *pgxpool.Pool
}

func NewInventoryService(repo repository.InventoryRepository, db *pgxpool.Pool) *InventoryService {
	return &InventoryService{repo: repo, db: db}
}

func (s *InventoryService) CreateProduct(req dto.ProductRequest) (*models.Product, error) {
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.Unit == "" {
		return nil, errors.New("product unit is required")
	}

	return s.repo.CreateProduct(&models.Product{
		Id:        uuid.New(),
		Name:      req.Name,
		Unit:      req.Unit,
		CreatedAt: time.Now(),
	})
}

func (s *InventoryService) GetProducts() ([]models.Product, error) {
	return s.repo.GetProducts()
}

func (s *InventoryService) GetProductByID(id string) (*models.Product, error) {
	productId, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product id")
	}
	product, err := s.repo.GetProductByID(productId)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *InventoryService) UpdateProduct(id string, req dto.ProductRequest) (*models.Product, error) {
	productId, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product id")
	}
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.Unit == "" {
		return nil, errors.New("product unit is required")
	}

	product, err := s.repo.UpdateProduct(&models.Product{
		Id:   productId,
		Name: req.Name,
		Unit: req.Unit,
	})
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *InventoryService) DeleteProduct(id string) error {
	productId, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product id")
	}
	return s.repo.DeleteProduct(productId)
}

func (s *InventoryService) AddStock(ctx context.Context, clinicAddressId string, req dto.InventoryQuantityRequest) (*models.AddressInventory, error) {
	addressId, err := uuid.Parse(clinicAddressId)
	if err != nil {
		return nil, errors.New("invalid clinic address id")
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, errors.New("invalid product id")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	if _, err := s.GetProductByID(req.ProductId); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	item, err := s.repo.GetInventoryByAddressAndProduct(addressId, productId, tx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if item == nil {
		item, err = s.repo.CreateInventoryTx(&models.AddressInventory{
			Id:              uuid.New(),
			ClinicAddressId: addressId,
			ProductId:       productId,
			Quantity:        req.Quantity,
			UpdatedAt:       now,
		}, tx)
	} else {
		item, err = s.repo.UpdateInventoryQuantityTx(item.Id, item.Quantity+req.Quantity, tx)
	}
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateTransactionTx(&models.InventoryTransaction{
		Id:              uuid.New(),
		ClinicAddressId: addressId,
		ProductId:       productId,
		Quantity:        req.Quantity,
		TransactionType: "restocked",
		CreatedAt:       now,
	}, tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.repo.GetInventoryByID(item.Id)
}

func (s *InventoryService) GetInventory(clinicAddressId string) ([]models.AddressInventory, error) {
	addressId, err := uuid.Parse(clinicAddressId)
	if err != nil {
		return nil, errors.New("invalid clinic address id")
	}
	return s.repo.GetInventoryByAddress(addressId)
}

func (s *InventoryService) UpdateInventory(ctx context.Context, clinicAddressId, inventoryId string, req dto.UpdateInventoryRequest) (*models.AddressInventory, error) {
	addressId, err := uuid.Parse(clinicAddressId)
	if err != nil {
		return nil, errors.New("invalid clinic address id")
	}
	itemId, err := uuid.Parse(inventoryId)
	if err != nil {
		return nil, errors.New("invalid inventory id")
	}
	if req.Quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	current, err := s.repo.GetInventoryByID(itemId)
	if err != nil {
		return nil, err
	}
	if current == nil || current.ClinicAddressId != addressId {
		return nil, errors.New("inventory not found")
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	updated, err := s.repo.UpdateInventoryQuantityTx(itemId, req.Quantity, tx)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("inventory not found")
	}

	if err := s.repo.CreateTransactionTx(&models.InventoryTransaction{
		Id:              uuid.New(),
		ClinicAddressId: addressId,
		ProductId:       current.ProductId,
		Quantity:        req.Quantity - current.Quantity,
		TransactionType: "manual_adjustment",
		CreatedAt:       time.Now(),
	}, tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.repo.GetInventoryByID(itemId)
}

func (s *InventoryService) AttachMaterial(clinicServiceId string, req dto.AttachMaterialRequest) (*models.ServiceMaterial, error) {
	serviceId, err := uuid.Parse(clinicServiceId)
	if err != nil {
		return nil, errors.New("invalid clinic service id")
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, errors.New("invalid product id")
	}
	if req.QuantityRequired <= 0 {
		return nil, errors.New("quantity_required must be greater than 0")
	}
	if _, err := s.GetProductByID(req.ProductId); err != nil {
		return nil, err
	}

	material, err := s.repo.CreateServiceMaterial(&models.ServiceMaterial{
		Id:               uuid.New(),
		ClinicServiceId:  serviceId,
		ProductId:        productId,
		QuantityRequired: req.QuantityRequired,
	})
	if err != nil {
		return nil, err
	}
	product, err := s.repo.GetProductByID(productId)
	if err != nil {
		return nil, err
	}
	if product != nil {
		material.ProductName = product.Name
		material.ProductUnit = product.Unit
	}
	return material, nil
}

func (s *InventoryService) GetServiceMaterials(clinicServiceId string) ([]models.ServiceMaterial, error) {
	serviceId, err := uuid.Parse(clinicServiceId)
	if err != nil {
		return nil, errors.New("invalid clinic service id")
	}
	return s.repo.GetServiceMaterials(serviceId)
}

func (s *InventoryService) GetTransactions(clinicAddressId, transactionType string) ([]models.InventoryTransaction, error) {
	addressId, err := uuid.Parse(clinicAddressId)
	if err != nil {
		return nil, errors.New("invalid clinic address id")
	}
	switch transactionType {
	case "", "restocked", "used", "manual_adjustment":
	default:
		return nil, errors.New("invalid transaction_type")
	}
	return s.repo.GetTransactions(addressId, transactionType)
}
