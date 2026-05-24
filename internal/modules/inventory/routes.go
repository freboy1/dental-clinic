package inventory

import (
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/inventory/handlers"
	"dental_clinic/internal/modules/inventory/repository"
	"dental_clinic/internal/modules/inventory/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPrivateRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	_ = cfg
	repo := repository.NewInventoryRepository(db)
	service := services.NewInventoryService(repo, db)
	handler := handlers.NewInventoryHandler(service)

	r.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/products", handler.GetProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handler.GetProductByID).Methods("GET")
	r.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")

	r.HandleFunc("/clinic-addresses/{id}/inventory", handler.AddStock).Methods("POST")
	r.HandleFunc("/clinic-addresses/{id}/inventory", handler.GetInventory).Methods("GET")
	r.HandleFunc("/clinic-addresses/{id}/inventory/{inventoryId}", handler.UpdateInventory).Methods("PUT")
	r.HandleFunc("/clinic-addresses/{id}/inventory-transactions", handler.GetTransactions).Methods("GET")

	r.HandleFunc("/clinic-services/{id}/materials", handler.AttachMaterial).Methods("POST")
}
