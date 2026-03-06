package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hasanraj3100/fridge-inventory/internal/api/handlers"
	"github.com/hasanraj3100/fridge-inventory/internal/api/routes"
	"github.com/hasanraj3100/fridge-inventory/internal/config"
	"github.com/hasanraj3100/fridge-inventory/internal/db"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
	"github.com/hasanraj3100/fridge-inventory/internal/service"
	"github.com/hasanraj3100/fridge-inventory/internal/utils"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println("config loaded successfully")

	dbCon, err := db.NewDatabaseConnection(cfg)
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		return
	}
	defer dbCon.Close()
	fmt.Println("database connected successfully")

	if err := db.Migrate(dbCon); err != nil {
		fmt.Println("failed to migrate database:", err)
		return
	}
	fmt.Println("database migrated successfully")

	// Utils
	jwtManager := utils.NewJWTManager(cfg.JWTSecretKey)
	passwordManager := utils.NewPasswordManager(12)

	// Repositories and Services
	txManager := repository.NewSQLTransactionProvider(dbCon)

	userRepo := repository.NewUserRepository(dbCon)
	userService := service.NewUserService(userRepo, passwordManager, jwtManager, txManager)

	fridgeRepo := repository.NewFridgeItemRepository(dbCon)
	fridgeItemService := service.NewFridgeItemService(fridgeRepo, txManager)

	itemUsageRepo := repository.NewItemUsageRepository(dbCon)
	itemUsageService := service.NewItemUsageService(itemUsageRepo, fridgeRepo, txManager)

	// Handlers
	authHandler := handlers.NewAuthHandler(userService)
	fridgeItemHandler := handlers.NewFridgeItemHandler(fridgeItemService)
	itemUsageHandler := handlers.NewItemUsageHandler(itemUsageService)

	// Router
	router := routes.NewRouter(authHandler, fridgeItemHandler, itemUsageHandler, jwtManager)
	handler := router.Setup()

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start : %v", err)
	}
}
