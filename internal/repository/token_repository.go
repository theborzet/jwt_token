package repository

import (
	"database/sql"
	"time"
)

func (r *ApiRepository) SaveRefreshTokenHash(userID, hash string, expiresAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at", userID, hash, expiresAt)
	if err != nil {
		r.logger.Printf("Error saving refresh token hash for user %s: %v", userID, err)
		return err
	}
	return nil
}

func (r *ApiRepository) GetRefreshTokenHash(userID string) (string, time.Time, error) {
	var hash string
	var expiresAt time.Time
	if err := r.db.QueryRow("SELECT token_hash, expires_at FROM refresh_tokens WHERE user_id = $1", userID).Scan(&hash, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return "", time.Time{}, nil
		}
		r.logger.Printf("Error getting refresh token hash for user %s: %v", userID, err)
		return "", time.Time{}, err
	}
	return hash, expiresAt, nil
}
