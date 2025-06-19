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

func (s *OrderService) Save(orderUpd models.Orders) (models.Orders, error) {
	err := validateOrder(orderUpd)
	if err != nil {
		return models.Orders{}, err
	}

	return s.repo.Save(orderUpd)
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) DeleteById(id int) error {
	return s.repo.DeleteById(id)
}

func validateOrder(order models.Orders) error {
	if order.Id != nil && *order.Id <= 0 {
		//return utils.New("Invalid Order Id")
		return ErrInvalidOrderId
	}

	if order.SourceLocationId != nil && *order.SourceLocationId <= 0 {
		return ErrInvalidSourceLocationId
	}

	if order.DestinationLocationId != nil && *order.DestinationLocationId <= 0 {
		return ErrInvalidDestinationLocationId
	}

	if order.Items != nil {
		for _, item := range *order.Items {
			if item.ItemIndex != nil && *item.ItemIndex <= 0 {
				return ErrInvalidIndex
			}

			if item.Id != nil && *item.Id < 0 {
				return ErrInvalidItemId
			}

			if item.RootId != nil && *item.RootId <= 0 {
				return ErrInvalidItemRootId
			}

			if item.ProductId != nil && *item.ProductId <= 0 {
				return ErrInvalidItemProductId
			}

			if item.VolumeValue != nil && *item.VolumeValue < 0 {
				return ErrInvalidItemVolumeValue
			}

			if item.WeightValue != nil && *item.WeightValue < 0 {
				return ErrInvalidItemWeightValue
			}
		}
	}

	return nil
}
