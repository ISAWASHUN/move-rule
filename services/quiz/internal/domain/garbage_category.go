package domain

import "context"

/*
***********
Entity
***********
*/
type GarbageCategory struct {
	ID   GarbageCategoryID
	Name string
}

/*
***********
Value Object
***********
*/
type GarbageCategoryID int

/*
***********
Repository
***********
*/
type GarbageCategoryRepository interface {
	GetAll(ctx context.Context) ([]GarbageCategory, error)
	GetByName(ctx context.Context, name string) (*GarbageCategory, error)
	GetByID(ctx context.Context, id int) (*GarbageCategory, error)
}
