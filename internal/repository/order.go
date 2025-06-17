package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"shop/internal/models"
)

type OrderRepository struct {
	db *gorm.DB
}

func (r *OrderRepository) Select(filter models.OrderFilter) ([]models.Orders, error) {
	var orders []models.Orders

	query := r.db.Model(&orders).
		Preload("SourceLocation").
		Preload("DestinationLocation").
		Preload("Items")

	if filter.Id != nil {
		op := getOperator(filter.IdOperator, "eq")

		switch op {
		case "eq":
			query = query.Where("id = ?", filter.Id)
		case "ne":
			query = query.Where("id != ?", *filter.Id)
		case "gt":
			query = query.Where("id > ?", *filter.Id)
		case "lt":
			query = query.Where("id < ?", *filter.Id)
		case "gte":
			query = query.Where("id >= ?", *filter.Id)
		case "lte":
			query = query.Where("id <= ?", *filter.Id)
		}
	}

	if filter.Comment != nil {
		op := getOperator(filter.CommentOperator, "contains")

		switch op {
		case "eq":
			query = query.Where(fmt.Sprintf("comment->>'%s' = ?", *filter.Lang), *filter.Comment)
		case "contains":
			query = query.Where(fmt.Sprintf("comment->>'%s' LIKE ?", *filter.Lang), "%"+*filter.Comment+"%")
		case "starts":
			query = query.Where(fmt.Sprintf("comment->>'%s' LIKE ?", *filter.Lang), *filter.Comment+"%")
		case "ends":
			query = query.Where(fmt.Sprintf("comment->>'%s' LIKE ?", *filter.Lang), "%"+*filter.Comment)
		}
	}

	if filter.SourceLocationID != nil {
		query = query.Where("source_location_id = ?", *filter.SourceLocationID)
	}

	if filter.DestinationLocationID != nil {
		query = query.Where("destination_location_id = ?", *filter.DestinationLocationID)
	}

	if filter.TotalWeight != nil {
		op := getOperator(filter.TotalWeightOperator, "gte")
		switch op {
		case "eq":
			query = query.Where("total_weight = ?", *filter.TotalWeight)
		case "ne":
			query = query.Where("total_weight != ?", *filter.TotalWeight)
		case "gt":
			query = query.Where("total_weight > ?", *filter.TotalWeight)
		case "lt":
			query = query.Where("total_weight < ?", *filter.TotalWeight)
		case "gte":
			query = query.Where("total_weight >= ?", *filter.TotalWeight)
		case "lte":
			query = query.Where("total_weight <= ?", *filter.TotalWeight)
		}
	}

	if filter.TotalVolume != nil {
		op := getOperator(filter.TotalVolumeOperator, "gte")
		switch op {
		case "eq":
			query = query.Where("total_volume = ?", *filter.TotalVolume)
		case "ne":
			query = query.Where("total_volume != ?", *filter.TotalVolume)
		case "gt":
			query = query.Where("total_volume > ?", *filter.TotalVolume)
		case "lt":
			query = query.Where("total_volume < ?", *filter.TotalVolume)
		case "gte":
			query = query.Where("total_volume >= ?", *filter.TotalVolume)
		case "lte":
			query = query.Where("total_volume <= ?", *filter.TotalVolume)
		}
	}

	if filter.SortBy != nil {
		order := "ASC"
		if filter.SortOrder != nil && strings.ToLower(*filter.SortOrder) == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", *filter.SortBy, order))
	}

	if filter.Limit != nil {
		query = query.Limit(*filter.Limit)
	}

	err := query.Find(&orders).Error
	if err != nil {
		return orders, err
	}

	if len(orders) == 0 {
		return orders, fmt.Errorf("no orders found")
	}

	return orders, nil

}

// Вспомогательная функция для получения оператора
func getOperator(op *string, defaultOp string) string {
	if op == nil {
		return defaultOp
	}
	return *op
}

func (r *OrderRepository) Save(orderUpd models.OrderRequest) (models.Orders, error) {
	transaction := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	var order models.Orders
	//обрабатываем Locations
	sourceLocation, err := r.UpdateLocations(transaction, &orderUpd.SourceLocation)
	if err != nil {
		transaction.Rollback()
		return order, err
	}

	destinationLocation, err := r.UpdateLocations(transaction, &orderUpd.DestinationLocation)
	if err != nil {
		transaction.Rollback()
		return order, err
	}

	order.Comment = orderUpd.Comment
	order.SourceLocationId = sourceLocation.Id
	order.DestinationLocationId = destinationLocation.Id
	order.TotalWeightMeasureUnitCode = orderUpd.TotalWeightMeasureUnitCode
	order.TotalVolumeMeasureUnitCode = orderUpd.TotalVolumeMeasureUnitCode

	if orderUpd.Id != nil {
		//Обновление
		order.Id = *orderUpd.Id
		err := transaction.Model(&order).Updates(&order).Error
		if err != nil {
			transaction.Rollback()
			return order, err
		}
		// Удаляем старые items
		if order.Items != nil {
			err = transaction.Where("root_id = ?", order.Id).Delete(&models.OrderItems{}).Error
			if err != nil {
				transaction.Rollback()
				return order, err
			}
		}
	} else {
		err := transaction.Create(&order).Error
		if err != nil {
			transaction.Rollback()
			return order, err
		}
	}

	// Обрабатываем Items
	_, err = r.UpdateItems(transaction, order.Id, orderUpd.Items)
	if err != nil {
		transaction.Rollback()
		return order, err
	}

	err = transaction.Commit().Error
	if err != nil {
		return order, err
	}

	fullOrder, err := r.getFullOrder(order.Id)
	if err != nil {
		return order, err
	}

	return fullOrder, nil
}

func (r *OrderRepository) getFullOrder(id uint64) (models.Orders, error) {
	var order models.Orders
	err := r.db.Preload("SourceLocation").
		Preload("DestinationLocation").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.WeightMeasureUnit").
		Preload("Items.VolumeMeasureUnit").
		First(&order, id).Error

	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *OrderRepository) UpdateItems(transaction *gorm.DB, orderId uint64, itemUpd []models.OrderItemRequest) ([]models.OrderItems, error) {
	var item models.OrderItems
	var items []models.OrderItems

	for _, itemUpd := range itemUpd {
		product, err := r.UpdateProducts(transaction, *itemUpd.Product)
		if err != nil {
			return items, err
		}

		item = models.OrderItems{
			ProductId:             product.Id,
			RootId:                orderId,
			ItemIndex:             itemUpd.ItemIndex,
			WeightValue:           itemUpd.WeightValue,
			WeightMeasureUnitCode: itemUpd.WeightMeasureUnitCode,
			VolumeValue:           itemUpd.VolumeValue,
			VolumeMeasureUnitCode: itemUpd.VolumeMeasureUnitCode,
		}

		if itemUpd.Id != nil {
			//Обновление текущей сущности
			item.Id = *itemUpd.Id

			err := transaction.Model(&item).Updates(&item).Error
			if err != nil {
				return items, err
			}
		} else {
			err := transaction.Create(&item).Error
			if err != nil {
				return items, err
			}
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *OrderRepository) UpdateProducts(transaction *gorm.DB, productUpd models.ProductsRequest) (*models.Products, error) {

	product := models.Products{
		Name: *productUpd.Name,
	}

	if productUpd.Id != nil {
		//Обновление текущей сущности
		product.Id = *productUpd.Id

		err := transaction.Model(&product).Updates(&product).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := transaction.Create(&product).Error
		if err != nil {
			return nil, err
		}
	}

	return &product, nil
}

func (r *OrderRepository) UpdateLocations(transaction *gorm.DB, locUpd *models.LocationUpsertRequest) (models.Locations, error) {
	locations := models.Locations{
		Name:    locUpd.Name,
		Address: locUpd.Address,
	}

	if locations.Address == "" && locations.Name.Ru == "" && locations.Name.En == "" {
		return locations, nil
	}

	if locUpd.Id != nil {
		//Обновление текущей сущности
		locations.Id = *locUpd.Id
		err := transaction.Model(&locations).Updates(&locations).Error
		if err != nil {
			return locations, err
		}
	} else {
		err := transaction.Create(&locations).Error
		if err != nil {
			return locations, err
		}
	}

	return locations, nil

}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) DeleteById(id int) error {

	err := r.db.Delete(&models.Orders{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
