package main

import (
	"log"
	"os"

	"github.com/X3nonxe/flashsale-engine/internal/delivery/http"
	"github.com/X3nonxe/flashsale-engine/internal/repository/postgres"
	"github.com/X3nonxe/flashsale-engine/internal/repository/redis"
	"github.com/X3nonxe/flashsale-engine/internal/usecase"
	"github.com/X3nonxe/flashsale-engine/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	db, err := database.NewPostgresDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	rdb, err := database.NewRedisClient(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	productRepo := postgres.NewProductRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	redisStockRepo := redis.NewRedisStockRepository(rdb)

	_ = productRepo

	purchaseUC := usecase.NewPurchaseUsecase(redisStockRepo, orderRepo)

	purchaseHandler := http.NewPurchaseHandler(purchaseUC)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/flash-sale/purchase", purchaseHandler.Purchase)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8080"
	}
	log.Printf("Flash Sale Engine running on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
