package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hasanraj3100/fridge-inventory/internal/api/handlers"
	"github.com/hasanraj3100/fridge-inventory/internal/api/routes"
	"github.com/hasanraj3100/fridge-inventory/internal/config"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
	"github.com/hasanraj3100/fridge-inventory/internal/service"
	"github.com/hasanraj3100/fridge-inventory/internal/utils"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println("config loaded successfully")

	db, err := repository.NewDatabaseConnection(cfg)
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		return
	}
	defer db.Close()
	fmt.Println("database connected successfully")

	if err := repository.Migrate(db); err != nil {
		fmt.Println("failed to migrate database:", err)
		return
	}
	fmt.Println("database migrated successfully")

	// Utils
	jwtManager := utils.NewJWTManager(cfg.JWTSecretKey)
	passwordManager := utils.NewPasswordManager(12)

	// Domain Related
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, passwordManager, jwtManager)

	// Handlers
	authHandler := handlers.NewAuthHandler(userService)

	// Router
	router := routes.NewRouter(authHandler)
	handler := router.Setup()

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start : %v", err)
	}
}
