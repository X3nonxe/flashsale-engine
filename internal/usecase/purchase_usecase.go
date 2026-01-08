package usecase

import (
	"context"
	"fmt"

	"github.com/X3nonxe/flashsale-engine/internal/domain"
)

type RedisRepository interface {
	DecrStock(ctx context.Context, key string, qty int) error
	IncrStock(ctx context.Context, key string, qty int) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
}

type PurchaseRequest struct {
	UserID    int64
	ProductID int64
	Quantity  int
}

type PurchaseUsecase struct {
	redisRepo RedisRepository
	orderRepo OrderRepository
}

func NewPurchaseUsecase(r RedisRepository, o OrderRepository) *PurchaseUsecase {
	return &PurchaseUsecase{
		redisRepo: r,
		orderRepo: o,
	}
}

func (uc *PurchaseUsecase) Purchase(ctx context.Context, req PurchaseRequest) error {
	stockKey := fmt.Sprintf("stock:%d", req.ProductID)

	err := uc.redisRepo.DecrStock(ctx, stockKey, req.Quantity)
	if err != nil {
		return fmt.Errorf("purchase failed: %w", err) 
	}

	order := &domain.Order{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    "PENDING",
	}

	err = uc.orderRepo.Create(ctx, order)
	if err != nil {
		_ = uc.redisRepo.IncrStock(ctx, stockKey, req.Quantity)
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}
