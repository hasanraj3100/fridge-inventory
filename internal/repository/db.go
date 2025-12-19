package repository

import (
	"fmt"

	"github.com/hasanraj3100/fridge-inventory/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabaseConnection(configuration *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s",
		configuration.DB.User,
		configuration.DB.Password,
		configuration.DB.Host,
		configuration.DB.Port,
		configuration.DB.Name,
	)

	dbConnection, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return dbConnection, nil
}
