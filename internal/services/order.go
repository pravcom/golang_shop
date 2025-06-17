package services

import (
	"shop/internal/models"
	"shop/internal/repository"
)

type OrderService struct {
	repo repository.Order
}

func (s *OrderService) Select(filter models.OrderFilter) ([]models.Orders, error) {
	return s.repo.Select(filter)

}

func (s *OrderService) Save(orderUpd models.OrderRequest) (models.Orders, error) {
	return s.repo.Save(orderUpd)
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) DeleteById(id int) error {
	return s.repo.DeleteById(id)
}
