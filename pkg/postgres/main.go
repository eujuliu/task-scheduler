package postgres

import (
	"fmt"
	"log/slog"
	"scheduler/internal/config"
	"scheduler/internal/entities"
	"scheduler/internal/persistence"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
)

type Database struct {
	db     *gorm.DB
	tx     *gorm.DB
	config *config.DatabaseConfig
	mu     sync.Mutex
}

func NewPostgres(config *config.DatabaseConfig) (*Database, error) {
	err := createInitialDatabase(config)
	if err != nil {
		return nil, err
	}

	db, err := connectToDatabase(config)
	if err != nil {
		return nil, err
	}

	err = db.Use(otelgorm.NewPlugin())
	if err != nil {
		return nil, err
	}

	DB := &Database{
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

func (db *Database) Get() *gorm.DB {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.tx != nil {
		return db.tx
	}

	return db.db
}

func (db *Database) BeginTransaction() {
	db.mu.Lock()
	defer db.mu.Unlock()
	slog.Debug("postgres transaction started...")
	db.tx = db.db.Begin()
}

func (db *Database) RollbackTransaction() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.tx == nil {
		return fmt.Errorf("you need to initialize the transaction first")
	}

	slog.Debug("postgres transaction rollback...")
	_ = db.tx.Rollback()

	db.tx = nil

	return nil
}

func (db *Database) CommitTransaction() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.tx == nil {
		return fmt.Errorf("you need to initialize the transaction first")
	}

	slog.Debug("postgres transaction commit...")
	_ = db.tx.Commit()

	db.tx = nil

	return nil
}

func (db *Database) SeedForTest() {
	user, _ := entities.NewUser("testuser", "testuser@email.com", "Password@123")

	user.AddCredits(100000000)

	m, _ := persistence.ToUserModel(user)
	_ = db.Get().Create(m).Error
}

func (db *Database) migrations() error {
	slog.Info("running migrations...")

	return db.db.AutoMigrate(
		&persistence.UserModel{},
		&persistence.PasswordRecoveryModel{},
		&persistence.TransactionModel{},
		&persistence.TaskModel{},
		&persistence.ErrorModel{},
	)
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
