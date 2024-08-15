package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveRefreshTokenHash(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		tokenID      string
		hash         string
		expiresAt    time.Time
		expectedErr  error
		mockBehavior func(mock sqlmock.Sqlmock, userID, tokenID, hash string)
	}{
		{
			name:      "Success",
			userID:    "1",
			tokenID:   "token-id-1",
			hash:      "test-hash",
			expiresAt: time.Now().Add(24 * time.Hour),
			mockBehavior: func(mock sqlmock.Sqlmock, userID, tokenID, hash string) {
				mock.ExpectExec("INSERT INTO refresh_tokens \\(user_id, token_id, token_hash, expires_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT \\(user_id\\) DO UPDATE SET token_id = EXCLUDED.token_id, token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at").
					WithArgs(userID, tokenID, hash, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name:      "Database Error",
			userID:    "1",
			tokenID:   "token-id-1",
			hash:      "test-hash",
			expiresAt: time.Now().Add(24 * time.Hour),
			mockBehavior: func(mock sqlmock.Sqlmock, userID, tokenID, hash string) {
				mock.ExpectExec("INSERT INTO refresh_tokens \\(user_id, token_id, token_hash, expires_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) ON CONFLICT \\(user_id\\) DO UPDATE SET token_id = EXCLUDED.token_id, token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at").
					WithArgs(userID, tokenID, hash, sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			logger := log.New(os.Stdout, "", log.LstdFlags)
			repo := NewApiRepository(db, logger)

			tt.mockBehavior(mock, tt.userID, tt.tokenID, tt.hash)

			err = repo.SaveRefreshTokenHash(tt.userID, tt.tokenID, tt.hash, tt.expiresAt)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRefreshTokenHash(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		expectedHash    string
		expectedTokenID string
		expectedTime    time.Time
		expectedErr     error
		mockBehavior    func(mock sqlmock.Sqlmock, userID string)
	}{
		{
			name:            "Success",
			userID:          "1",
			expectedHash:    "test-hash",
			expectedTokenID: "token-id-1",
			expectedTime:    time.Now().Add(24 * time.Hour),
			mockBehavior: func(mock sqlmock.Sqlmock, userID string) {
				rows := sqlmock.NewRows([]string{"token_hash", "token_id", "expires_at"}).
					AddRow("test-hash", "token-id-1", time.Now().Add(24*time.Hour))
				mock.ExpectQuery(`SELECT token_hash, token_id, expires_at FROM refresh_tokens WHERE user_id = \$1`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedErr: nil,
		},
		{
			name:   "No Rows",
			userID: "1",
			mockBehavior: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT token_hash, token_id, expires_at FROM refresh_tokens WHERE user_id = \$1`).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedHash:    "",
			expectedTokenID: "",
			expectedTime:    time.Time{},
			expectedErr:     nil,
		},
		{
			name:   "Database Error",
			userID: "1",
			mockBehavior: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT token_hash, token_id, expires_at FROM refresh_tokens WHERE user_id = \$1`).
					WithArgs(userID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedHash:    "",
			expectedTokenID: "",
			expectedTime:    time.Time{},
			expectedErr:     sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			logger := log.New(os.Stdout, "", log.LstdFlags)
			repo := NewApiRepository(db, logger)

			tt.mockBehavior(mock, tt.userID)

			hash, tokenID, expiresAt, err := repo.GetRefreshTokenHash(tt.userID)
			assert.Equal(t, tt.expectedHash, hash)
			assert.Equal(t, tt.expectedTokenID, tokenID)
			if !tt.expectedTime.IsZero() {
				assert.WithinDuration(t, tt.expectedTime, expiresAt, time.Minute)
			}
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
