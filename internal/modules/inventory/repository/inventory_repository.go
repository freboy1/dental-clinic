package repository

import (
	"context"
	"time"

	"dental_clinic/internal/modules/inventory/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository interface {
	CreateProduct(product *models.Product) (*models.Product, error)
	GetProducts() ([]models.Product, error)
	GetProductByID(id uuid.UUID) (*models.Product, error)
	UpdateProduct(product *models.Product) (*models.Product, error)
	DeleteProduct(id uuid.UUID) error

	GetInventoryByAddress(clinicAddressId uuid.UUID) ([]models.AddressInventory, error)
	GetInventoryByID(id uuid.UUID) (*models.AddressInventory, error)
	GetInventoryByAddressAndProduct(clinicAddressId, productId uuid.UUID, tx pgx.Tx) (*models.AddressInventory, error)
	CreateInventoryTx(inventory *models.AddressInventory, tx pgx.Tx) (*models.AddressInventory, error)
	UpdateInventoryQuantityTx(id uuid.UUID, quantity float64, tx pgx.Tx) (*models.AddressInventory, error)
	CreateTransactionTx(transaction *models.InventoryTransaction, tx pgx.Tx) error
	GetTransactions(clinicAddressId uuid.UUID, transactionType string) ([]models.InventoryTransaction, error)

	CreateServiceMaterial(material *models.ServiceMaterial) (*models.ServiceMaterial, error)
	GetServiceMaterials(clinicServiceId uuid.UUID) ([]models.ServiceMaterial, error)
}

type inventoryRepo struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) InventoryRepository {
	return &inventoryRepo{db: db}
}

func (r *inventoryRepo) CreateProduct(product *models.Product) (*models.Product, error) {
	query := `INSERT INTO products (id, name, unit, created_at) VALUES ($1, $2, $3, $4) RETURNING id, name, unit, created_at`
	err := r.db.QueryRow(context.Background(), query, product.Id, product.Name, product.Unit, product.CreatedAt).
		Scan(&product.Id, &product.Name, &product.Unit, &product.CreatedAt)
	return product, err
}

func (r *inventoryRepo) GetProducts() ([]models.Product, error) {
	query := `SELECT id, name, unit, created_at FROM products ORDER BY name`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Unit, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *inventoryRepo) GetProductByID(id uuid.UUID) (*models.Product, error) {
	query := `SELECT id, name, unit, created_at FROM products WHERE id = $1`
	product := &models.Product{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&product.Id, &product.Name, &product.Unit, &product.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *inventoryRepo) UpdateProduct(product *models.Product) (*models.Product, error) {
	query := `UPDATE products SET name = $2, unit = $3 WHERE id = $1 RETURNING id, name, unit, created_at`
	err := r.db.QueryRow(context.Background(), query, product.Id, product.Name, product.Unit).
		Scan(&product.Id, &product.Name, &product.Unit, &product.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *inventoryRepo) DeleteProduct(id uuid.UUID) error {
	result, err := r.db.Exec(context.Background(), `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *inventoryRepo) GetInventoryByAddress(clinicAddressId uuid.UUID) ([]models.AddressInventory, error) {
	query := `
		SELECT ai.id, ai.clinic_address_id, ai.product_id, p.name, p.unit, ai.quantity, ai.updated_at
		FROM address_inventory ai
		JOIN products p ON p.id = ai.product_id
		WHERE ai.clinic_address_id = $1
		ORDER BY p.name
	`
	rows, err := r.db.Query(context.Background(), query, clinicAddressId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := make([]models.AddressInventory, 0)
	for rows.Next() {
		var item models.AddressInventory
		if err := rows.Scan(&item.Id, &item.ClinicAddressId, &item.ProductId, &item.ProductName, &item.ProductUnit, &item.Quantity, &item.UpdatedAt); err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}
	return inventory, rows.Err()
}

func (r *inventoryRepo) GetInventoryByID(id uuid.UUID) (*models.AddressInventory, error) {
	query := `
		SELECT ai.id, ai.clinic_address_id, ai.product_id, p.name, p.unit, ai.quantity, ai.updated_at
		FROM address_inventory ai
		JOIN products p ON p.id = ai.product_id
		WHERE ai.id = $1
	`
	item := &models.AddressInventory{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&item.Id, &item.ClinicAddressId, &item.ProductId, &item.ProductName, &item.ProductUnit, &item.Quantity, &item.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *inventoryRepo) GetInventoryByAddressAndProduct(clinicAddressId, productId uuid.UUID, tx pgx.Tx) (*models.AddressInventory, error) {
	query := `
		SELECT id, clinic_address_id, product_id, quantity, updated_at
		FROM address_inventory
		WHERE clinic_address_id = $1 AND product_id = $2
		FOR UPDATE
	`
	item := &models.AddressInventory{}
	err := tx.QueryRow(context.Background(), query, clinicAddressId, productId).Scan(&item.Id, &item.ClinicAddressId, &item.ProductId, &item.Quantity, &item.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *inventoryRepo) CreateInventoryTx(inventory *models.AddressInventory, tx pgx.Tx) (*models.AddressInventory, error) {
	query := `
		INSERT INTO address_inventory (id, clinic_address_id, product_id, quantity, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, clinic_address_id, product_id, quantity, updated_at
	`
	err := tx.QueryRow(context.Background(), query, inventory.Id, inventory.ClinicAddressId, inventory.ProductId, inventory.Quantity, inventory.UpdatedAt).
		Scan(&inventory.Id, &inventory.ClinicAddressId, &inventory.ProductId, &inventory.Quantity, &inventory.UpdatedAt)
	return inventory, err
}

func (r *inventoryRepo) UpdateInventoryQuantityTx(id uuid.UUID, quantity float64, tx pgx.Tx) (*models.AddressInventory, error) {
	query := `
		UPDATE address_inventory
		SET quantity = $2, updated_at = $3
		WHERE id = $1
		RETURNING id, clinic_address_id, product_id, quantity, updated_at
	`
	item := &models.AddressInventory{}
	err := tx.QueryRow(context.Background(), query, id, quantity, time.Now()).
		Scan(&item.Id, &item.ClinicAddressId, &item.ProductId, &item.Quantity, &item.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *inventoryRepo) CreateTransactionTx(transaction *models.InventoryTransaction, tx pgx.Tx) error {
	query := `
		INSERT INTO inventory_transactions (id, clinic_address_id, product_id, quantity, transaction_type, appointment_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	var appointmentId interface{}
	if transaction.AppointmentId != uuid.Nil {
		appointmentId = transaction.AppointmentId
	}
	_, err := tx.Exec(context.Background(), query, transaction.Id, transaction.ClinicAddressId, transaction.ProductId, transaction.Quantity, transaction.TransactionType, appointmentId, transaction.CreatedAt)
	return err
}

func (r *inventoryRepo) GetTransactions(clinicAddressId uuid.UUID, transactionType string) ([]models.InventoryTransaction, error) {
	query := `
		SELECT it.id, it.clinic_address_id, it.product_id, p.name, it.quantity, it.transaction_type, COALESCE(it.appointment_id, '00000000-0000-0000-0000-000000000000'::uuid), it.created_at
		FROM inventory_transactions it
		JOIN products p ON p.id = it.product_id
		WHERE it.clinic_address_id = $1
			AND ($2 = '' OR it.transaction_type = $2)
		ORDER BY it.created_at DESC
	`
	rows, err := r.db.Query(context.Background(), query, clinicAddressId, transactionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.InventoryTransaction, 0)
	for rows.Next() {
		var transaction models.InventoryTransaction
		if err := rows.Scan(&transaction.Id, &transaction.ClinicAddressId, &transaction.ProductId, &transaction.ProductName, &transaction.Quantity, &transaction.TransactionType, &transaction.AppointmentId, &transaction.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, rows.Err()
}

func (r *inventoryRepo) CreateServiceMaterial(material *models.ServiceMaterial) (*models.ServiceMaterial, error) {
	query := `
		INSERT INTO service_materials (id, service_id, product_id, quantity_required)
		VALUES ($1, $2, $3, $4)
		RETURNING id, service_id, product_id, quantity_required
	`
	err := r.db.QueryRow(context.Background(), query, material.Id, material.ClinicServiceId, material.ProductId, material.QuantityRequired).
		Scan(&material.Id, &material.ClinicServiceId, &material.ProductId, &material.QuantityRequired)
	return material, err
}

func (r *inventoryRepo) GetServiceMaterials(clinicServiceId uuid.UUID) ([]models.ServiceMaterial, error) {
	query := `
		SELECT sm.id, sm.service_id, sm.product_id, p.name, p.unit, sm.quantity_required
		FROM service_materials sm
		JOIN products p ON p.id = sm.product_id
		WHERE sm.service_id = $1
		ORDER BY p.name
	`
	rows, err := r.db.Query(context.Background(), query, clinicServiceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	materials := make([]models.ServiceMaterial, 0)
	for rows.Next() {
		var material models.ServiceMaterial
		if err := rows.Scan(&material.Id, &material.ClinicServiceId, &material.ProductId, &material.ProductName, &material.ProductUnit, &material.QuantityRequired); err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, rows.Err()
}
