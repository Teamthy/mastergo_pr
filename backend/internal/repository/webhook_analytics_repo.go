package repository

import (
	"context"
	"encoding/json"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type WebhookAnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewWebhookAnalyticsRepository(db *pgxpool.Pool) *WebhookAnalyticsRepository {
	return &WebhookAnalyticsRepository{db: db}
}

// Webhook Repository Methods
func (r *WebhookAnalyticsRepository) CreateWebhook(ctx context.Context, webhook *models.Webhook) error {
	query := `
		INSERT INTO webhooks (id, user_id, url, events, headers, active, retry_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	headersJSON, _ := json.Marshal(webhook.Headers)

	_, err := r.db.Exec(ctx, query,
		webhook.ID, webhook.UserID, webhook.URL, pq.Array(webhook.Events),
		headersJSON, webhook.Active, webhook.RetryCount, webhook.CreatedAt, webhook.UpdatedAt)

	return err
}

func (r *WebhookAnalyticsRepository) GetWebhook(ctx context.Context, webhookID uuid.UUID) (*models.Webhook, error) {
	query := `SELECT id, user_id, url, events, headers, active, retry_count, last_triggered_at, created_at, updated_at FROM webhooks WHERE id = $1`

	row := r.db.QueryRow(ctx, query, webhookID)

	var webhook models.Webhook
	var headersJSON []byte
	var lastTriggered *time.Time

	err := row.Scan(&webhook.ID, &webhook.UserID, &webhook.URL, pq.Array(&webhook.Events),
		&headersJSON, &webhook.Active, &webhook.RetryCount, &lastTriggered, &webhook.CreatedAt, &webhook.UpdatedAt)

	if err != nil {
		return nil, err
	}

	webhook.LastTriggered = lastTriggered
	json.Unmarshal(headersJSON, &webhook.Headers)

	return &webhook, nil
}

func (r *WebhookAnalyticsRepository) UpdateWebhook(ctx context.Context, webhook *models.Webhook) error {
	query := `
		UPDATE webhooks
		SET url = $1, events = $2, headers = $3, active = $4, retry_count = $5, last_triggered_at = $6, updated_at = $7
		WHERE id = $8
	`

	headersJSON, _ := json.Marshal(webhook.Headers)

	_, err := r.db.Exec(ctx, query,
		webhook.URL, pq.Array(webhook.Events), headersJSON, webhook.Active, webhook.RetryCount,
		webhook.LastTriggered, webhook.UpdatedAt, webhook.ID)

	return err
}

func (r *WebhookAnalyticsRepository) DeleteWebhook(ctx context.Context, webhookID uuid.UUID) error {
	query := `DELETE FROM webhooks WHERE id = $1`
	_, err := r.db.Exec(ctx, query, webhookID)
	return err
}

func (r *WebhookAnalyticsRepository) GetWebhooksByUserID(ctx context.Context, userID uuid.UUID) ([]models.Webhook, error) {
	query := `SELECT id, user_id, url, events, headers, active, retry_count, last_triggered_at, created_at, updated_at FROM webhooks WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []models.Webhook
	for rows.Next() {
		var webhook models.Webhook
		var headersJSON []byte
		var lastTriggered *time.Time

		err := rows.Scan(&webhook.ID, &webhook.UserID, &webhook.URL, pq.Array(&webhook.Events),
			&headersJSON, &webhook.Active, &webhook.RetryCount, &lastTriggered, &webhook.CreatedAt, &webhook.UpdatedAt)

		if err != nil {
			return nil, err
		}

		webhook.LastTriggered = lastTriggered
		json.Unmarshal(headersJSON, &webhook.Headers)

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// Webhook Event Repository Methods
func (r *WebhookAnalyticsRepository) CreateWebhookEvent(ctx context.Context, event *models.WebhookEvent) error {
	query := `
		INSERT INTO webhook_events (id, webhook_id, event_type, payload, status, attempts, max_attempts, response_code, response_body, next_retry_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	payloadJSON, _ := json.Marshal(event.Payload)

	_, err := r.db.Exec(ctx, query,
		event.ID, event.WebhookID, event.EventType, payloadJSON, event.Status,
		event.Attempts, event.MaxAttempts, event.ResponseCode, event.ResponseBody,
		event.NextRetryAt, event.CreatedAt, event.UpdatedAt)

	return err
}

func (r *WebhookAnalyticsRepository) UpdateWebhookEvent(ctx context.Context, event *models.WebhookEvent) error {
	query := `
		UPDATE webhook_events
		SET status = $1, attempts = $2, response_code = $3, response_body = $4, next_retry_at = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.Exec(ctx, query,
		event.Status, event.Attempts, event.ResponseCode, event.ResponseBody, event.NextRetryAt, event.UpdatedAt, event.ID)

	return err
}

func (r *WebhookAnalyticsRepository) GetFailedWebhookEvents(ctx context.Context) ([]models.WebhookEvent, error) {
	query := `
		SELECT id, webhook_id, event_type, payload, status, attempts, max_attempts, response_code, response_body, next_retry_at, created_at, updated_at
		FROM webhook_events
		WHERE status IN ('retrying', 'failed') AND (next_retry_at IS NULL OR next_retry_at <= NOW())
		LIMIT 100
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.WebhookEvent
	for rows.Next() {
		var event models.WebhookEvent
		var payloadJSON []byte

		err := rows.Scan(&event.ID, &event.WebhookID, &event.EventType, &payloadJSON,
			&event.Status, &event.Attempts, &event.MaxAttempts, &event.ResponseCode,
			&event.ResponseBody, &event.NextRetryAt, &event.CreatedAt, &event.UpdatedAt)

		if err != nil {
			return nil, err
		}

		json.Unmarshal(payloadJSON, &event.Payload)
		events = append(events, event)
	}

	return events, nil
}

// API Analytics Repository Methods
func (r *WebhookAnalyticsRepository) CreateAPIAnalytics(ctx context.Context, analytics *models.APIAnalytics) error {
	query := `
		INSERT INTO api_analytics (id, api_key_id, user_id, endpoint, method, response_time_ms, status_code, request_size, response_size, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(ctx, query,
		analytics.ID, analytics.APIKeyID, analytics.UserID, analytics.Endpoint, analytics.Method,
		analytics.ResponseTimeMs, analytics.StatusCode, analytics.RequestSize, analytics.ResponseSize, analytics.CreatedAt)

	return err
}

func (r *WebhookAnalyticsRepository) GetAPIAnalyticsByKeyID(ctx context.Context, apiKeyID uuid.UUID, startDate time.Time) ([]models.APIAnalytics, error) {
	query := `
		SELECT id, api_key_id, user_id, endpoint, method, response_time_ms, status_code, request_size, response_size, created_at
		FROM api_analytics
		WHERE api_key_id = $1 AND created_at >= $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, apiKeyID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []models.APIAnalytics
	for rows.Next() {
		var a models.APIAnalytics

		err := rows.Scan(&a.ID, &a.APIKeyID, &a.UserID, &a.Endpoint, &a.Method,
			&a.ResponseTimeMs, &a.StatusCode, &a.RequestSize, &a.ResponseSize, &a.CreatedAt)

		if err != nil {
			return nil, err
		}

		analytics = append(analytics, a)
	}

	return analytics, nil
}

func (r *WebhookAnalyticsRepository) GetAPIAnalyticsByUserID(ctx context.Context, userID uuid.UUID, startDate time.Time) ([]models.APIAnalytics, error) {
	query := `
		SELECT id, api_key_id, user_id, endpoint, method, response_time_ms, status_code, request_size, response_size, created_at
		FROM api_analytics
		WHERE user_id = $1 AND created_at >= $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []models.APIAnalytics
	for rows.Next() {
		var a models.APIAnalytics

		err := rows.Scan(&a.ID, &a.APIKeyID, &a.UserID, &a.Endpoint, &a.Method,
			&a.ResponseTimeMs, &a.StatusCode, &a.RequestSize, &a.ResponseSize, &a.CreatedAt)

		if err != nil {
			return nil, err
		}

		analytics = append(analytics, a)
	}

	return analytics, nil
}

// Rate Limit Stats Repository Methods
func (r *WebhookAnalyticsRepository) CreateRateLimitStats(ctx context.Context, stats *models.RateLimitStats) error {
	query := `
		INSERT INTO rate_limit_stats (id, api_key_id, date, request_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		stats.ID, stats.APIKeyID, stats.Date, stats.RequestCount, stats.CreatedAt, stats.UpdatedAt)

	return err
}

func (r *WebhookAnalyticsRepository) GetRateLimitStats(ctx context.Context, apiKeyID uuid.UUID, date time.Time) (*models.RateLimitStats, error) {
	query := `SELECT id, api_key_id, date, request_count, created_at, updated_at FROM rate_limit_stats WHERE api_key_id = $1 AND date = $2`

	row := r.db.QueryRow(ctx, query, apiKeyID, date)

	var stats models.RateLimitStats
	err := row.Scan(&stats.ID, &stats.APIKeyID, &stats.Date, &stats.RequestCount, &stats.CreatedAt, &stats.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *WebhookAnalyticsRepository) UpdateRateLimitStats(ctx context.Context, stats *models.RateLimitStats) error {
	query := `UPDATE rate_limit_stats SET request_count = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, stats.RequestCount, stats.UpdatedAt, stats.ID)
	return err
}
