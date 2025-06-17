package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math"

	"gorm.io/gorm"
)

type Multilang struct {
	Ru string `json:"ru,omitempty"`
	En string `json:"en,omitempty"`
}

type MeasureUnits struct {
	Id   uint64    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	Code string    `json:"code" gorm:"type:text;unique;not null"`
	Name Multilang `json:"name" gorm:"type:jsonb"`
}

type Products struct {
	Id   uint64    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	Name Multilang `json:"name" gorm:"type:jsonb; column:name"`
}

type Locations struct {
	Id      uint64    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	Name    Multilang `json:"name" gorm:"type:jsonb; column:name"`
	Address string    `json:"address" gorm:"type:text; column:address"`
}

type Orders struct {
	Id                         uint64        `json:"id" gorm:"column:id;primary_key;auto_increment"`
	Comment                    Multilang     `json:"comment" gorm:"type:jsonb;column:comment"`
	SourceLocationId           uint64        `gorm:"column:source_location_id"`
	SourceLocation             *Locations    `json:"source_location,omitempty" gorm:"foreignKey:SourceLocationId"`
	DestinationLocationId      uint64        `gorm:"column:destination_location_id"`
	DestinationLocation        *Locations    `json:"destination_location" gorm:"foreignKey:DestinationLocationId"`
	TotalItemsCount            int           `json:"total_items_count" gorm:"-"`
	TotalWeightNumeric         float64       `json:"total_weight_numeric" gorm:"-"`
	TotalWeightMeasureUnitCode string        `json:"total_weight_measure_unit_code gorm:column:total_weight_measure_unit_code"`
	TotalWeightMeasureUnit     *MeasureUnits `json:"totalWeightMeasureUnit,omitempty" gorm:"foreignKey:TotalWeightMeasureUnitCode;references:Code"`
	TotalVolume                float64       `json:"total_volume" gorm:"-"`
	TotalVolumeMeasureUnitCode string        `json:"total_volume_measure_unit_code gorm:column:total_volume_measure_unit_code"`
	TotalVolumeMeasureUnit     *MeasureUnits `json:"totalVolumeMeasureUnit,omitempty" gorm:"foreignKey:TotalVolumeMeasureUnitCode;references:Code"`
	Items                      []*OrderItems `json:"items" gorm:"foreignKey:RootId"`
}

type OrderItems struct {
	Id                    uint64        `json:"id" gorm:"column:id;primary_key;auto_increment"`
	RootId                uint64        `gorm:"column:root_id"`
	Root                  *Orders       `json:"root" gorm:"foreignKey:RootId;references:id"`
	ProductId             uint64        `json:"product_id" gorm:"column:product_id"`
	Product               *Products     `json:"product" gorm:"foreignKey:ProductId"`
	ItemIndex             int           `json:"item_index" gorm:"column:item_index"`
	WeightValue           float64       `json:"weight_value" gorm:"column:weight_value"`
	WeightMeasureUnitCode string        `gorm:"column:weight_measure_unit_code"`
	WeightMeasureUnit     *MeasureUnits `json:"weight_measure_unit_code" gorm:"foreignKey:WeightMeasureUnitCode;references:Code"`
	VolumeValue           float64       `json:"volume_value" gorm:"column:volume_value"`
	VolumeMeasureUnitCode string        `gorm:"column:volume_measure_unit_code"`
	VolumeMeasureUnit     *MeasureUnits `json:"volume_measure_unit" gorm:"foreignKey:VolumeMeasureUnitCode;references:Code"`
}

// OrderUpsertRequest - структура запроса
type OrderRequest struct {
	Id                         *uint64               `json:"id,omitempty"` // nil для создания
	Comment                    Multilang             `json:"comment"`
	SourceLocation             LocationUpsertRequest `json:"source_location"`
	DestinationLocation        LocationUpsertRequest `json:"destination_location"`
	Items                      []OrderItemRequest    `json:"items"`
	TotalWeightMeasureUnitCode string                `json:"total_weight_measure_unit_code"`
	TotalVolumeMeasureUnitCode string                `json:"total_volume_measure_unit_code"`
}

// LocationRequest - запрос для Location
type LocationUpsertRequest struct {
	Id      *uint64   `json:"id,omitempty"`
	Name    Multilang `json:"name"`
	Address string    `json:"address"`
}

type ProductsRequest struct {
	Id   *uint64    `json:"id"`
	Name *Multilang `json:"name"`
}

// OrderItemRequest - запрос для OrderItem
type OrderItemRequest struct {
	Id                    *uint64          `json:"id,omitempty"`
	Product               *ProductsRequest `json:"product"`
	ItemIndex             int              `json:"item_index"`
	WeightValue           float64          `json:"weight_value"`
	WeightMeasureUnitCode string           `json:"weight_measure_unit_code"`
	VolumeValue           float64          `json:"volume_value"`
	VolumeMeasureUnitCode string           `json:"volume_measure_unit_code"`
}

type OrderFilter struct {
	Id         *uint64 `form:"id"`
	IdOperator *string `form:"id_operator"`

	Comment         *string `form:"comment"`
	CommentOperator *string `form:"comment_operator"` // "eq", "contains", "starts", "ends"

	SourceLocationID      *uint64 `form:"source_location_id"`
	DestinationLocationID *uint64 `form:"destination_location_id"`

	TotalWeight         *float64 `form:"total_weight"`
	TotalWeightOperator *string  `form:"total_weight_operator"`

	TotalVolume         *float64 `form:"total_volume"`
	TotalVolumeOperator *string  `form:"total_volume_operator"`

	//CreatedAfter *string `form:"created_after"` // Дата в формате "2006-01-02"
	//HasItems  *bool   `form:"has_items"`
	SortBy    *string `form:"sort_by"`    // "id", "created_at", etc.
	SortOrder *string `form:"sort_order"` // "asc" или "desc"
	Limit     *int    `form:"limit"`
	//Offset    *int    `form:"offset"`
	Lang *string `form:"lang"`
}

func (m *Multilang) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, m)
}

func (m Multilang) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (o *Orders) AfterFind(tx *gorm.DB) (err error) {
	o.TotalItemsCount = CalculateTotalItemsCount(o.Items)
	o.TotalVolume = CalculateTotalVolume(o.Items)
	o.TotalWeightNumeric = CalculateTotalWeight(o.Items)

	return nil
}

func CalculateTotalWeight(items []*OrderItems) float64 {
	total := 0.0
	for _, item := range items {
		total += item.WeightValue
	}
	return total
}

func CalculateTotalItemsCount(items []*OrderItems) int {
	return len(items)
}

func CalculateTotalVolume(items []*OrderItems) float64 {
	total := 0.0
	for _, item := range items {
		total += item.VolumeValue
	}

	if total != 0.0 {
		total = math.Round(total*100) / 100
	}

	return total
}
