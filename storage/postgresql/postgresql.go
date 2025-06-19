package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Важно!
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBconfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DBName  string
	SSLMode string
}
type Storage struct {
	DB *gorm.DB
}

var cfgDB DBconfig

func New(cfg DBconfig) *Storage {

	storagePath := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s port=%s",
		cfg.Host,
		cfg.User,
		cfg.DBName,
		cfg.Pass,
		cfg.SSLMode,
		cfg.Port)

	cfgDB = cfg

	gormDB, err := gorm.Open(postgres.Open(storagePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	sqlDB, _ := gormDB.DB()

	err = runMigrations(sqlDB)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{DB: gormDB}
}

func (s *Storage) Close() {
	sqlDB, err := s.DB.DB()
	if err != nil {
		log.Fatal("Failed to close connection to DB", err)
	}

	_ = sqlDB.Close()
}

func runMigrations(db *sql.DB) error {
	// Проверяем рабочую директорию
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Не удалось получить рабочую директорию:")

	}
	log.Println("Рабочая директория:", wd)

	driver, err := postgresMigrate.WithInstance(db, &postgresMigrate.Config{
		DatabaseName: cfgDB.DBName,
		SchemaName:   "public",
	})
	if err != nil {
		return fmt.Errorf("Failed to get driver postgresMigrate: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://schema", "postgres", driver)
	if err != nil {
		return fmt.Errorf("Failed to create migration instance: %w", err)
	}

	err = m.Down()
	if err != nil {
		return fmt.Errorf("Faile to down migrations: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("Failed to get version of migrations: %w", err)
	}

	log.Printf("Database migrations version: %d, dirty: %v", version, dirty)

	return nil
}
