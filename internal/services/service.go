package services

import (
	"shop/internal/models"
	"shop/internal/repository"
)

type Order interface {
	DeleteById(id int) error
	Save(orderUpd models.Orders) (models.Orders, error)
	Select(filter models.OrderFilter) ([]models.Orders, error)
}
type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {

	return &Service{NewOrderService(repos)}

}
