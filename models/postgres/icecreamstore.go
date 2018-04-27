package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/sudarshan-reddy/benjerry/db"
	"github.com/sudarshan-reddy/benjerry/models"
)

type iceCreamStore struct {
	*db.DB
}

//NewIceCreamStore returns a new instance of IceCreamStore that
//is coupled to postgresql
func NewIceCreamStore(db *db.DB) models.IceCreamStore {
	return &iceCreamStore{db}
}

func (i *iceCreamStore) StoreContext(ctx context.Context, iceCreamInput models.IceCream) error {
	query := `
	INSERT INTO ice_cream (name,
    image_open, 
    image_closed,
    story,
    description,
    sourcing_values,
    ingredients,
    allergy_info, 
    dietary_certification,
    product_id)
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    ON CONFLICT DO NOTHING
	`

	db, err := i.GetContextDB(ctx)
	if err != nil {
		return fmt.Errorf("error preparing context: %s", err)
	}

	_, err = db.ExecContext(ctx, query, iceCreamInput.Name, iceCreamInput.ImageOpen,
		iceCreamInput.ImageClosed, iceCreamInput.Story, iceCreamInput.Description,
		pq.Array(iceCreamInput.SourcingValues), pq.Array(iceCreamInput.Ingredients),
		iceCreamInput.AllergyInfo, iceCreamInput.DietaryCertification,
		iceCreamInput.ProductID)

	return err
}

func (i *iceCreamStore) Get(name string) (*models.IceCream, error) {
	query := `
	SELECT name,
    image_open, 
    image_closed,
    story,
    description,
    sourcing_values,
    ingredients,
    allergy_info, 
    dietary_certification,
    product_id 
    FROM ice_cream 
    WHERE name = $1
    `

	var iceCream models.IceCream
	err := i.QueryRow(query, name).Scan(&iceCream.Name, &iceCream.ImageOpen,
		&iceCream.ImageClosed, &iceCream.Story, &iceCream.Description,
		pq.Array(&iceCream.SourcingValues), pq.Array(&iceCream.Ingredients), &iceCream.AllergyInfo,
		&iceCream.DietaryCertification, &iceCream.ProductID)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &iceCream, nil
}

func (i *iceCreamStore) GetAll() ([]models.IceCream, error) {
	return nil, errors.New("not implemented")
}

func (i *iceCreamStore) Update(iceCreamInput models.IceCream) error {
	query := `
		UPDATE ice_cream
		SET  
		image_open = COALESCE(NULLIF($2,''), image_open), 
    	image_closed = COALESCE(NULLIF($3, ''), image_closed),
    	story = COALESCE(NULLIF($4,''), story),
    	description = COALESCE(NULLIF($5,''), description),
    	sourcing_values = COALESCE($6, sourcing_values),
    	ingredients = COALESCE($7, ingredients),
    	allergy_info = COALESCE($8, allergy_info), 
    	dietary_certification = COALESCE(NULLIF($9,''), dietary_certification),
    	product_id = COALESCE(NULLIF($10,''), product_id)
    	WHERE name = $1
    `

	_, err := i.Exec(query, iceCreamInput.Name, iceCreamInput.ImageOpen,
		iceCreamInput.ImageClosed, iceCreamInput.Story, iceCreamInput.Description,
		pq.Array(iceCreamInput.SourcingValues), pq.Array(iceCreamInput.Ingredients),
		iceCreamInput.AllergyInfo, iceCreamInput.DietaryCertification,
		iceCreamInput.ProductID)

	return err
}

func (i *iceCreamStore) Delete(name string) error {
	query := `
	DELETE FROM ice_cream
	WHERE name = $1
	`

	_, err := i.Exec(query, name)
	return err
}
