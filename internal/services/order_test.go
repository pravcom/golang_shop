package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"shop/internal/models"
)

func TestValidateOrder(t *testing.T) {
	cases := []struct {
		name      string
		order     models.Orders
		expectErr error
	}{{
		name: "Order is 0",
		order: models.Orders{
			Id: uint64Ptr(0),
		},
		expectErr: ErrInvalidOrderId,
	},
		{
			name: "SourceLocationID is 0",
			order: models.Orders{
				SourceLocationId: uint64Ptr(0),
			},
			expectErr: ErrInvalidSourceLocationId,
		},
		{
			name: "DestinationLocationId is 0",
			order: models.Orders{
				DestinationLocationId: uint64Ptr(0),
			},
			expectErr: ErrInvalidDestinationLocationId,
		},
		{
			name: "Item.Id is 0",
			order: models.Orders{
				Items: orderItemsPtr([]models.OrderItems{
					{
						Id: uint64Ptr(0),
					},
				}),
			},
			expectErr: ErrInvalidItemId,
		},
		{
			name: "Item.RootId is 0",
			order: models.Orders{
				Items: orderItemsPtr([]models.OrderItems{
					{
						RootId: uint64Ptr(0),
					},
				}),
			},
			expectErr: ErrInvalidItemRootId,
		},
		{
			name: "Item.Index is 0",
			order: models.Orders{
				Items: orderItemsPtr([]models.OrderItems{
					{
						ItemIndex: intPtr(0),
					},
				}),
			},
			expectErr: ErrInvalidIndex,
		},
		{
			name: "Item.ProductId is 0",
			order: models.Orders{
				Items: orderItemsPtr([]models.OrderItems{
					{
						ProductId: uint64Ptr(0),
					},
				}),
			},
			expectErr: ErrInvalidItemProductId,
		},

		{
			name: "Item.VolumeValue is below zero",
			order: models.Orders{
				Items: orderItemsPtr([]models.OrderItems{
					{
						VolumeValue: float64Ptr(-1),
					},
				}),
			},
			expectErr: ErrInvalidItemVolumeValue,
		},
		{
			name: "Item.WeightValue is below zero",
			order: models.Orders{
				Items: orderItemsPtr([]models.OrderItems{
					{
						WeightValue: float64Ptr(-1),
					},
				}),
			},
			expectErr: ErrInvalidItemWeightValue,
		},
		{
			name: "Success test",
			order: models.Orders{
				Id:                    uint64Ptr(1),
				SourceLocationId:      uint64Ptr(1),
				DestinationLocationId: uint64Ptr(2),
				Items: orderItemsPtr([]models.OrderItems{
					{
						Id:          uint64Ptr(1),
						RootId:      uint64Ptr(1),
						ProductId:   uint64Ptr(1),
						VolumeValue: float64Ptr(23),
						WeightValue: float64Ptr(21),
					},
				}),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateOrder(c.order)

			require.ErrorIs(t, err, c.expectErr)

		})
	}

}

func uint64Ptr(val uint64) *uint64 {
	return &val
}

func orderItemsPtr(items []models.OrderItems) *[]models.OrderItems {
	return &items
}

func intPtr(val int) *int {
	return &val
}

func float64Ptr(val float64) *float64 {
	return &val
}

func TestSelect(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	//repo := repoMock.NewMockOrder(ctl)
	//TODO test
	//service := services.NewService(repo)
}
