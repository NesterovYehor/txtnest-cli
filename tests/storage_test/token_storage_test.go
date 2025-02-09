package storagetest

import (
	"database/sql"
	"testing"
	"time"

	"github.com/NesterovYehor/txtnest-cli/internal/models"
	"github.com/NesterovYehor/txtnest-cli/internal/storage"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var testTokensData = models.TokenData{
	AccessToken:      "test-access-token",
	RefreshToken:     "test-refresh-token",
	ExpiresAt:        time.Now().Add(time.Hour),
	RefreshExpiresAt: time.Now().Add(24 * time.Hour),
}

func TestTokenStorage_Success(t *testing.T) {
	// Use an in-memory database for isolation
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	err = storage.InitTokenStorage(db)
	assert.NoError(t, err)

	ts, err := storage.GetTokenStorage()
	assert.NoError(t, err)

	err = ts.SaveTokens(&testTokensData)
	assert.NoError(t, err)

	res, err := ts.GetTokens()
	assert.NoError(t, err)
	assert.Equal(t, testTokensData.AccessToken, res.AccessToken)
	assert.Equal(t, testTokensData.RefreshToken, res.RefreshToken)
	// Because times can have slight differences, you might want to compare within a delta:
	assert.WithinDuration(t, testTokensData.ExpiresAt, res.ExpiresAt, time.Second)
	assert.WithinDuration(t, testTokensData.RefreshExpiresAt, res.RefreshExpiresAt, time.Second)
}
