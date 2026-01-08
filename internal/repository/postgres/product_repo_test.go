package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.TODO()

	query := `
SELECT p.id, p.name, p.price, ps.quantity
FROM products p
JOIN product_stocks ps ON p.id = ps.product_id
WHERE p.id = $1
`

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"}).
			AddRow(int64(1), "iPhone 15", int64(15000000), 10)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnRows(rows)

		product, qty, err := repo.GetProductByID(ctx, 1)

		require.NoError(t, err)
		require.NotNil(t, product)

		assert.Equal(t, int64(1), product.ID)
		assert.Equal(t, "iPhone 15", product.Name)
		assert.Equal(t, float64(15000000), product.Price)
		assert.Equal(t, 10, qty)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		product, qty, err := repo.GetProductByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, 0, qty)
	})
}

func TestUpdateStock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.TODO()

	query := `
UPDATE product_stocks
SET quantity = quantity - $1, last_updated = now()
WHERE product_id = $2 AND quantity >= $1
`

	t.Run("SuccessUpdate", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(1, 101).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		err := repo.UpdateStock(ctx, 101, 1)

		assert.NoError(t, err)
	})

	t.Run("FailedInsufficientStock", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(5, 101).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		err := repo.UpdateStock(ctx, 101, 5)

		assert.Error(t, err)
		assert.Equal(t, "stock update failed: insufficient stock or invalid product", err.Error())
	})
}
