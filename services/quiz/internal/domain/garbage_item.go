package domain

import "context"

/*
***********
Entity
***********
*/
type GarbageItem struct {
	ID                GarbageItemID
	MunicipalityID    MunicipalityID
	GarbageCategoryID GarbageCategoryID
	AreaName          string
	ItemName          string
	Notes             string
	Remarks           string
	BulkGarbageFee    int
}

/*
***********
Value Object
***********
*/
type GarbageItemID int

/*
***********
Repository
***********
*/
type GarbageItemRepository interface {
	GetByMunicipalityID(ctx context.Context, municipalityID int) ([]GarbageItem, error)
	GetByID(ctx context.Context, id int) (*GarbageItem, error)
	GetByIDWithCategory(ctx context.Context, id int) (*GarbageItem, *GarbageCategory, error)
}
