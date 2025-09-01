package postgres

import (
	"fmt"
	"log/slog"
	"scheduler/internal/config"
	"scheduler/internal/persistence"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db     *gorm.DB
	tx     *gorm.DB
	config *config.DatabaseConfig
}

var DB *Database

func Load(config *config.DatabaseConfig) (*Database, error) {
	if DB != nil {
		return DB, nil
	}

	err := createInitialDatabase(config)
	if err != nil {
		return nil, err
	}

	db, err := connectToDatabase(config)
	if err != nil {
		return nil, err
	}

	DB = &Database{
		db:     db,
		tx:     nil,
		config: config,
	}

	err = DB.migrations()
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func (db *Database) GetInstance() *gorm.DB {
	if db.tx != nil {
		return db.tx
	}

	return db.db
}

func (db *Database) BeginTransaction() {
	slog.Debug("transaction started...")
	db.tx = db.db.Begin()
}

func (db *Database) RollbackTransaction() error {
	if db.tx == nil {
		return fmt.Errorf("you need to initialize the transction first")
	}

	slog.Debug("transaction rollback...")
	_ = db.tx.Rollback()

	db.tx = nil

	return nil
}

func (db *Database) CommitTransaction() error {
	if db.tx == nil {
		return fmt.Errorf("you need to initialize the transction first")
	}

	slog.Debug("transaction commit...")
	_ = db.tx.Commit()

	db.tx = nil

	return nil
}

func (db *Database) migrations() error {
	slog.Info("running migrations...")

	return db.db.AutoMigrate(&persistence.UserModel{})
}

func connectToDatabase(config *config.DatabaseConfig) (*gorm.DB, error) {
	slog.Debug(fmt.Sprintf("connecting to database %s", config.DBName))
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("connected to database %s", config.DBName))

	return db, nil
}

func createInitialDatabase(conf *config.DatabaseConfig) error {
	slog.Debug("creating initial database...")
	newConfig := &config.DatabaseConfig{
		DBName:   "postgres",
		Host:     conf.Host,
		Port:     conf.Port,
		User:     conf.User,
		Password: conf.Password,
	}

	db, err := connectToDatabase(newConfig)
	if err != nil {
		return err
	}

	sql, err := db.DB()
	if err != nil {
		return err
	}

	defer func() {
		_ = sql.Close()
	}()

	var exists bool

	query := fmt.Sprintf(
		"SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s');",
		conf.DBName,
	)

	slog.Debug(query)

	err = sql.QueryRow(query).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		slog.Debug("initial database already created...")
		return nil
	}

	query = fmt.Sprintf("CREATE DATABASE %s", conf.DBName)

	_, err = sql.Exec(query)
	if err != nil {
		return err
	}

	slog.Debug("initial database created...")

	return nil
}
