package database

import (
	"context"
	"log"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool                       *pgxpool.Pool
	UserRepository             *repository.UserRepository
	ApiKeyRepository           *repository.ApiKeyRepository
	AdvancedRepository         *repository.AdvancedRepository
	WebhookAnalyticsRepository *repository.WebhookAnalyticsRepository
}

func NewPostgresConnection(dbURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL")
	return pool, nil
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		Pool:                       pool,
		UserRepository:             repository.NewUserRepository(pool),
		ApiKeyRepository:           repository.NewApiKeyRepository(pool),
		AdvancedRepository:         repository.NewAdvancedRepository(pool),
		WebhookAnalyticsRepository: repository.NewWebhookAnalyticsRepository(pool),
	}
}

// User Repository delegation methods
func (db *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return db.UserRepository.GetByEmail(ctx, email)
}

func (db *Database) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return db.UserRepository.GetByID(ctx, userID)
}

// Advanced Repository delegation methods - Password Reset
func (db *Database) CreatePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	return db.AdvancedRepository.CreatePasswordResetToken(ctx, token)
}

func (db *Database) GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	return db.AdvancedRepository.GetPasswordResetTokenByHash(ctx, tokenHash)
}

func (db *Database) MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string, usedAt time.Time) error {
	return db.AdvancedRepository.MarkPasswordResetTokenUsed(ctx, tokenHash, usedAt)
}

func (db *Database) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	return db.AdvancedRepository.UpdateUserPassword(ctx, userID, newPassword)
}

// Advanced Repository delegation methods - Audit Logs
func (db *Database) CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error {
	return db.AdvancedRepository.CreateAuditLog(ctx, auditLog)
}

func (db *Database) GetAuditLogsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.AuditLog, int64, error) {
	return db.AdvancedRepository.GetAuditLogsByUserID(ctx, userID, limit, offset)
}

func (db *Database) GetAuditLogsByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]models.AuditLog, error) {
	return db.AdvancedRepository.GetAuditLogsByResource(ctx, resourceType, resourceID, limit, offset)
}

func (db *Database) GetAuditActionSummary(ctx context.Context, userID uuid.UUID, days int) (map[string]int, error) {
	return db.AdvancedRepository.GetAuditActionSummary(ctx, userID, days)
}

// Webhook Analytics Repository delegation methods - Webhooks
func (db *Database) CreateWebhook(ctx context.Context, webhook *models.Webhook) error {
	return db.WebhookAnalyticsRepository.CreateWebhook(ctx, webhook)
}

func (db *Database) GetWebhook(ctx context.Context, webhookID uuid.UUID) (*models.Webhook, error) {
	return db.WebhookAnalyticsRepository.GetWebhook(ctx, webhookID)
}

func (db *Database) UpdateWebhook(ctx context.Context, webhook *models.Webhook) error {
	return db.WebhookAnalyticsRepository.UpdateWebhook(ctx, webhook)
}

func (db *Database) DeleteWebhook(ctx context.Context, webhookID uuid.UUID) error {
	return db.WebhookAnalyticsRepository.DeleteWebhook(ctx, webhookID)
}

func (db *Database) GetWebhooksByUserID(ctx context.Context, userID uuid.UUID) ([]models.Webhook, error) {
	return db.WebhookAnalyticsRepository.GetWebhooksByUserID(ctx, userID)
}

// Webhook Analytics Repository delegation methods - Webhook Events
func (db *Database) CreateWebhookEvent(ctx context.Context, event *models.WebhookEvent) error {
	return db.WebhookAnalyticsRepository.CreateWebhookEvent(ctx, event)
}

func (db *Database) UpdateWebhookEvent(ctx context.Context, event *models.WebhookEvent) error {
	return db.WebhookAnalyticsRepository.UpdateWebhookEvent(ctx, event)
}

func (db *Database) GetFailedWebhookEvents(ctx context.Context) ([]models.WebhookEvent, error) {
	return db.WebhookAnalyticsRepository.GetFailedWebhookEvents(ctx)
}

// Webhook Analytics Repository delegation methods - Analytics
func (db *Database) CreateAPIAnalytics(ctx context.Context, analytics *models.APIAnalytics) error {
	return db.WebhookAnalyticsRepository.CreateAPIAnalytics(ctx, analytics)
}

func (db *Database) GetAPIAnalyticsByKeyID(ctx context.Context, apiKeyID uuid.UUID, startDate time.Time) ([]models.APIAnalytics, error) {
	return db.WebhookAnalyticsRepository.GetAPIAnalyticsByKeyID(ctx, apiKeyID, startDate)
}

func (db *Database) GetAPIAnalyticsByUserID(ctx context.Context, userID uuid.UUID, startDate time.Time) ([]models.APIAnalytics, error) {
	return db.WebhookAnalyticsRepository.GetAPIAnalyticsByUserID(ctx, userID, startDate)
}

// Webhook Analytics Repository delegation methods - Rate Limit Stats
func (db *Database) CreateRateLimitStats(ctx context.Context, stats *models.RateLimitStats) error {
	return db.WebhookAnalyticsRepository.CreateRateLimitStats(ctx, stats)
}

func (db *Database) GetRateLimitStats(ctx context.Context, apiKeyID uuid.UUID, date time.Time) (*models.RateLimitStats, error) {
	return db.WebhookAnalyticsRepository.GetRateLimitStats(ctx, apiKeyID, date)
}

func (db *Database) UpdateRateLimitStats(ctx context.Context, stats *models.RateLimitStats) error {
	return db.WebhookAnalyticsRepository.UpdateRateLimitStats(ctx, stats)
}
