package postgres

import (
	"context"
	"database/sql"

	"github.com/X3nonxe/flashsale-engine/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
        INSERT INTO orders (user_id, product_id, quantity, status, created_at)
        VALUES ($1, $2, $3, 'PENDING', now())
        RETURNING id
    `
	return r.db.QueryRowContext(ctx, query,
		order.UserID, order.ProductID, order.Quantity,
	).Scan(&order.ID)
}
