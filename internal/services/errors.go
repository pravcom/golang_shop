package services

import "errors"

var (
	ErrInvalidOrderId               = errors.New("Invalid Order Id")
	ErrInvalidSourceLocationId      = errors.New("Invalid Order SourceLocationId")
	ErrInvalidDestinationLocationId = errors.New("Invalid Order DestinationLocationId")
	ErrInvalidIndex                 = errors.New("Invalid Item Index")
	ErrInvalidItemId                = errors.New("Invalid Item Id")
	ErrInvalidItemRootId            = errors.New("Invalid Item RootId")
	ErrInvalidItemProductId         = errors.New("Invalid Item ProductId")
	ErrInvalidItemVolumeValue       = errors.New("Invalid Item VolumeValue")
	ErrInvalidItemWeightValue       = errors.New("Invalid Item WeightValue")
)
