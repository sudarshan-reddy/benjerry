package models

import (
	"context"
	"database/sql"

	"github.com/sudarshan-reddy/benjerry/db"
)

var (
	//ErrNoRows is a wrapper on sql.ErrNoRows
	ErrNoRows = sql.ErrNoRows
	//ErrRowAlreadyExists is thrown when we get pgerr.Code = "23505"
	ErrRowAlreadyExists = "row already exists"
)

//IceCream defines the model for IceCreamStore
type IceCream struct {
	Name                 string   `json:"name"`
	ImageOpen            string   `json:"image_open"`
	ImageClosed          string   `json:"image_closed"`
	Story                string   `json:"story"`
	Description          string   `json:"description"`
	SourcingValues       []string `json:"sourcing_values"`
	Ingredients          []string `json:"ingredients"`
	AllergyInfo          string   `json:"allergy_info"`
	DietaryCertification string   `json:"dietary_certification"`
	ProductID            string   `json:"product_id"`
}

//IceCreamStore specifies the operations to be performed
//for storing IceCream data
type IceCreamStore interface {
	db.TransactionalStore
	StoreContext(ctx context.Context, iceCreamInput IceCream) error
	Get(name string) (*IceCream, error)
	GetAll() ([]IceCream, error)
	Update(iceCreamInput IceCream) error
	Delete(name string) error
}
