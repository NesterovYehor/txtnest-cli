package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/NesterovYehor/txtnest-cli/internal/models"
)

type TokenStorage struct {
	DB *sql.DB
}

var (
	instance *TokenStorage
	once     sync.Once
	initErr  error
)

// InitTokenStorage initializes the singleton instance using a once.Do block.
func InitTokenStorage(db *sql.DB) error {
	once.Do(func() {
		createTableQuery := `
            CREATE TABLE IF NOT EXISTS tokens (
            id INTEGER PRIMARY KEY,
            access_token TEXT,
            refresh_token TEXT,
            expires_at DATETIME,
            refresh_expires_at DATETIME
        );
		`
		_, err := db.Exec(createTableQuery)
		if err != nil {
			initErr = fmt.Errorf("failed to create tokens table: %w", err)
		}
		instance = &TokenStorage{DB: db}
	})
	if initErr != nil {
		return initErr
	}
	return nil
}

func GetTokenStorage() (*TokenStorage, error) {
	if instance == nil {
		return nil, errors.New("token storage not initialized; call InitTokenStorage() first")
	}
	return instance, nil
}

func (ts *TokenStorage) SaveTokens(jwt *models.TokenData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := ts.DB.ExecContext(ctx, "DELETE FROM tokens")
	if err != nil {
		return fmt.Errorf("failed cleaning up token table: %v", err)
	}
	query := `
        INSERT INTO tokens (access_token, refresh_token, expires_at, refresh_expires_at)
        VALUES ($1, $2, $3, $4)
    `
	args := []any{
		jwt.AccessToken,
		jwt.RefreshToken,
		jwt.ExpiresAt,
		jwt.RefreshExpiresAt,
	}
	_, err = ts.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Failed to save tokens data:%v", err)
	}
	return nil
}

func (ts *TokenStorage) GetTokens() (*models.TokenData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `
        SELECT access_token, refresh_token, expires_at, refresh_expires_at 
        FROM tokens
    `

	var tokens models.TokenData
	err := ts.DB.QueryRowContext(ctx, query).Scan(
		&tokens.AccessToken,
		&tokens.RefreshToken,
		&tokens.ExpiresAt,
		&tokens.RefreshExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to get tokens data form db: %v", err)
	}
	return &tokens, nil
}
