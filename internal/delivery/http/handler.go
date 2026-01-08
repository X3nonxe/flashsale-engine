package http

import (
	"net/http"

	"github.com/X3nonxe/flashsale-engine/internal/usecase"
	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	usecase *usecase.PurchaseUsecase
}

func NewPurchaseHandler(uc *usecase.PurchaseUsecase) *PurchaseHandler {
	return &PurchaseHandler{usecase: uc}
}

type purchaseRequestDTO struct {
	UserID    int64 `json:"user_id" binding:"required"`
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

func (h *PurchaseHandler) Purchase(c *gin.Context) {
	var req purchaseRequestDTO
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ucReq := usecase.PurchaseRequest{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	err := h.usecase.Purchase(c.Request.Context(), ucReq)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": "fail",
			"message": "Purchase failed. Likely out of stock.",
			"debug_error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"message": "Order placed successfully!",
	})
}