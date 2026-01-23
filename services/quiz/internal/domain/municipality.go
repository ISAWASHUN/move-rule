package domain

import "context"

/*
***********
Entity
***********
*/
type Municipality struct {
	ID   MunicipalityID
	Code int
	Name string
}

/*
***********
Value Object
***********
*/
type MunicipalityID int

/*
***********
Repository
***********
*/
type MunicipalityRepository interface {
	GetByName(ctx context.Context, name string) (*Municipality, error)
	GetByID(ctx context.Context, id int) (*Municipality, error)
	GetByCode(ctx context.Context, code int) (*Municipality, error)
}
