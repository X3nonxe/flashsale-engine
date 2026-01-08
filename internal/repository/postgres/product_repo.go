package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/X3nonxe/flashsale-engine/internal/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProductByID(ctx context.Context, id int64) (*domain.Product, int, error) {
	query := `
		SELECT p.id, p.name, p.price, ps.quantity
		FROM products p
		JOIN product_stocks ps ON p.id = ps.product_id
		WHERE p.id = $1
	`

	var p domain.Product
	var stockQty int

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Price, &stockQty,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, errors.New("product not found")
		}
		return nil, 0, err
	}

	return &p, stockQty, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, productID int64, qty int) error {
	query := `
		UPDATE product_stocks 
		SET quantity = quantity - $1, last_updated = now()
		WHERE product_id = $2 AND quantity >= $1
	`

	res, err := r.db.ExecContext(ctx, query, qty, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("stock update failed: insufficient stock or invalid product")
	}

	return nil
}