package repository

import (
	"gorm.io/gorm"
	"shop/internal/models"
)

type Order interface {
	DeleteById(id int) error
	Save(orderUpd models.OrderRequest) (models.Orders, error)
	Select(filter models.OrderFilter) ([]models.Orders, error)
}
type Repository struct {
	Order
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{Order: NewOrderRepository(db)}
}
