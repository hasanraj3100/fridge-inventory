package main

import (
	"fmt"

	"github.com/hasanraj3100/fridge-inventory/internal/config"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
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
}
