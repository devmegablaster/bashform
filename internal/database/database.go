package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func New(ctx context.Context, wg *sync.WaitGroup, cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		if err := sqlDB.Close(); err != nil {
			slog.Error("❌ Unable to close database connection", slog.String("error", err.Error()))
			return
		}

		slog.Info("🔌 Database connection closed")
	}()

	conn.Logger = logger.Discard
	conn.TranslateError = true

	db := &Database{
		DB: conn,
	}

	if err := db.createExtension(); err != nil {
		return nil, err
	}

	if err := db.autoMigrate(); err != nil {
		return nil, err
	}

	slog.Info("✅ Connected to database")

	return db, nil
}

func (d *Database) autoMigrate() error {
	return d.DB.AutoMigrate(&models.User{}, &models.Form{}, &models.Question{}, &models.Option{}, &models.Response{}, &models.Answer{})
}

func (d *Database) createExtension() error {
	return d.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
}
