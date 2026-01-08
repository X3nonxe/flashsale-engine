package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestDecrStock(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisStockRepository(db)
	ctx := context.TODO()
	key := "stock:1"

	// Logika: Ambil stok -> Cek cukup/enggak -> Kurangi jika cukup
	expectedScript := `
		local stock = redis.call("GET", KEYS[1])
		if not stock then return 0 end
		stock = tonumber(stock)
		local qty = tonumber(ARGV[1])
		if stock >= qty then
			redis.call("DECRBY", KEYS[1], qty)
			return 1
		end
		return 0
	`

	t.Run("Success", func(t *testing.T) {
		qty := 1
		mock.ExpectEval(expectedScript, []string{key}, qty).
			SetVal(int64(1))

		err := repo.DecrStock(ctx, key, qty)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InsufficientStock", func(t *testing.T) {
		qty := 5
		mock.ExpectEval(expectedScript, []string{key}, qty).
			SetVal(int64(0))

		err := repo.DecrStock(ctx, key, qty)

		assert.Error(t, err)
		assert.Equal(t, "stock update failed: insufficient stock", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RedisError", func(t *testing.T) {
		qty := 1
		mock.ExpectEval(expectedScript, []string{key}, qty).
			SetErr(errors.New("connection refused"))

		err := repo.DecrStock(ctx, key, qty)

		assert.Error(t, err)
	})
}