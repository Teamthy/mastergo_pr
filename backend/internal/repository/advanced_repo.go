package repository

import (
	"context"
	"encoding/json"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdvancedRepository struct {
	db *pgxpool.Pool
}

func NewAdvancedRepository(db *pgxpool.Pool) *AdvancedRepository {
	return &AdvancedRepository{db: db}
}

// Password Reset Repository Methods
func (r *AdvancedRepository) CreatePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *AdvancedRepository) GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, used_at, created_at FROM password_reset_tokens WHERE token_hash = $1`

	row := r.db.QueryRow(ctx, query, tokenHash)

	var token models.PasswordResetToken
	err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AdvancedRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string, usedAt time.Time) error {
	query := `UPDATE password_reset_tokens SET used_at = $1 WHERE token_hash = $2`
	_, err := r.db.Exec(ctx, query, usedAt, tokenHash)
	return err
}

func (r *AdvancedRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, newPassword, time.Now(), userID)
	return err
}

// Audit Log Repository Methods
func (r *AdvancedRepository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	oldValuesJSON, _ := json.Marshal(log.OldValues)
	newValuesJSON, _ := json.Marshal(log.NewValues)

	query := `
		INSERT INTO audit_logs (id, user_id, action, resource_type, resource_id, old_values, new_values, ip_address, user_agent, status, error_message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Exec(ctx, query,
		log.ID, log.UserID, log.Action, log.ResourceType, log.ResourceID,
		oldValuesJSON, newValuesJSON, log.IPAddress, log.UserAgent, log.Status, log.ErrorMessage, log.CreatedAt)

	return err
}

func (r *AdvancedRepository) GetAuditLogsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.AuditLog, int64, error) {
	query := `SELECT id, user_id, action, resource_type, resource_id, old_values, new_values, ip_address, user_agent, status, error_message, created_at FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var oldValuesJSON, newValuesJSON []byte

		err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&oldValuesJSON, &newValuesJSON, &log.IPAddress, &log.UserAgent, &log.Status, &log.ErrorMessage, &log.CreatedAt)

		if err != nil {
			return nil, 0, err
		}

		json.Unmarshal(oldValuesJSON, &log.OldValues)
		json.Unmarshal(newValuesJSON, &log.NewValues)

		logs = append(logs, log)
	}

	// Get total count
	var count int64
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1`
	r.db.QueryRow(ctx, countQuery, userID).Scan(&count)

	return logs, count, nil
}

func (r *AdvancedRepository) GetAuditLogsByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]models.AuditLog, error) {
	query := `SELECT id, user_id, action, resource_type, resource_id, old_values, new_values, ip_address, user_agent, status, error_message, created_at FROM audit_logs WHERE resource_type = $1 AND resource_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(ctx, query, resourceType, resourceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var oldValuesJSON, newValuesJSON []byte

		err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&oldValuesJSON, &newValuesJSON, &log.IPAddress, &log.UserAgent, &log.Status, &log.ErrorMessage, &log.CreatedAt)

		if err != nil {
			return nil, err
		}

		json.Unmarshal(oldValuesJSON, &log.OldValues)
		json.Unmarshal(newValuesJSON, &log.NewValues)

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *AdvancedRepository) GetAuditActionSummary(ctx context.Context, userID uuid.UUID, days int) (map[string]int, error) {
	query := `
		SELECT action, COUNT(*) as count 
		FROM audit_logs 
		WHERE user_id = $1 AND created_at > NOW() - INTERVAL '1 day' * $2
		GROUP BY action
	`

	rows, err := r.db.Query(ctx, query, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[string]int)
	for rows.Next() {
		var action string
		var count int

		if err := rows.Scan(&action, &count); err != nil {
			return nil, err
		}

		summary[action] = count
	}

	return summary, nil
}
