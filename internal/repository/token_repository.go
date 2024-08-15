package repository

import (
	"database/sql"
	"time"
)

func (r *ApiRepository) SaveRefreshTokenHash(userID, tokenID, hash string, expiresAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token_id, token_hash, expires_at) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET token_id = EXCLUDED.token_id, token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at", userID, tokenID, hash, expiresAt)
	if err != nil {
		r.logger.Printf("Error saving refresh token hash for user %s: %v", userID, err)
		return err
	}
	return nil
}

func (r *ApiRepository) GetRefreshTokenHash(userID string) (string, string, time.Time, error) {
	var hash, tokenID string
	var expiresAt time.Time
	if err := r.db.QueryRow("SELECT token_hash, token_id, expires_at FROM refresh_tokens WHERE user_id = $1", userID).Scan(&hash, &tokenID, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return "", "", time.Time{}, nil
		}
		r.logger.Printf("Error getting refresh token hash for user %s: %v", userID, err)
		return "", "", time.Time{}, err
	}
	return hash, tokenID, expiresAt, nil
}
