package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/X3nonxe/flashsale-engine/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCKS DEFINITION ---
type MockRedisRepo struct {
    mock.Mock
}
func (m *MockRedisRepo) DecrStock(ctx context.Context, key string, qty int) error {
    args := m.Called(ctx, key, qty)
    return args.Error(0)
}
func (m *MockRedisRepo) IncrStock(ctx context.Context, key string, qty int) error {
    args := m.Called(ctx, key, qty)
    return args.Error(0)
}

type MockOrderRepo struct {
    mock.Mock
}
func (m *MockOrderRepo) Create(ctx context.Context, order *domain.Order) error {
    args := m.Called(ctx, order)
    return args.Error(0)
}

// --- TEST CASES ---
func TestPurchase(t *testing.T) {
    mockRedis := new(MockRedisRepo)
    mockOrder := new(MockOrderRepo)
    uc := NewPurchaseUsecase(mockRedis, mockOrder)
    ctx := context.TODO()

    req := PurchaseRequest{
        UserID:    100,
        ProductID: 1,
        Quantity:  1,
    }

    // Skenario 1: Sukses
    t.Run("Success", func(t *testing.T) {
        // Expect Redis Decr OK
        mockRedis.On("DecrStock", ctx, "stock:1", 1).Return(nil).Once()
        // Expect DB Create OK
        mockOrder.On("Create", ctx, mock.AnythingOfType("*domain.Order")).Return(nil).Once()

        err := uc.Purchase(ctx, req)
        assert.NoError(t, err)
    })

    // Skenario 2: Stok Habis (Redis Fail)
    t.Run("OutOfStock", func(t *testing.T) {
        mockRedis.On("DecrStock", ctx, "stock:1", 1).Return(errors.New("insufficient stock")).Once()
        
        err := uc.Purchase(ctx, req)
        assert.Error(t, err)
        mockOrder.AssertNotCalled(t, "Create") // Pastikan DB tidak dipanggil
    })

    // Skenario 3: DB Error (Harus Rollback Redis)
    t.Run("DatabaseFail_Rollback", func(t *testing.T) {
        // Redis OK
        mockRedis.On("DecrStock", ctx, "stock:1", 1).Return(nil).Once()
        // DB Fail
        mockOrder.On("Create", ctx, mock.AnythingOfType("*domain.Order")).Return(errors.New("db error")).Once()
        // Expect Rollback (Incr)
        mockRedis.On("IncrStock", ctx, "stock:1", 1).Return(nil).Once()

        err := uc.Purchase(ctx, req)
        assert.Error(t, err)
    })
}