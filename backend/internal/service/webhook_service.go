package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

type WebhookService struct {
	db *database.Database
}

func NewWebhookService(db *database.Database) *WebhookService {
	return &WebhookService{db: db}
}

// CreateWebhook creates a new webhook
func (s *WebhookService) CreateWebhook(ctx context.Context, userID uuid.UUID, url string, events []string) (*models.Webhook, error) {
	webhook := &models.Webhook{
		ID:        uuid.New(),
		UserID:    userID,
		URL:       url,
		Events:    events,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.CreateWebhook(ctx, webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}

// UpdateWebhook updates webhook configuration
func (s *WebhookService) UpdateWebhook(ctx context.Context, webhookID uuid.UUID, url string, events []string, active bool) (*models.Webhook, error) {
	webhook, err := s.db.GetWebhook(ctx, webhookID)
	if err != nil {
		return nil, err
	}

	webhook.URL = url
	webhook.Events = events
	webhook.Active = active
	webhook.UpdatedAt = time.Now()

	if err := s.db.UpdateWebhook(ctx, webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}

// DeleteWebhook deletes a webhook
func (s *WebhookService) DeleteWebhook(ctx context.Context, webhookID uuid.UUID) error {
	return s.db.DeleteWebhook(ctx, webhookID)
}

// GetUserWebhooks retrieves all webhooks for a user
func (s *WebhookService) GetUserWebhooks(ctx context.Context, userID uuid.UUID) ([]models.Webhook, error) {
	return s.db.GetWebhooksByUserID(ctx, userID)
}

// TriggerWebhook sends event to webhook
func (s *WebhookService) TriggerWebhook(ctx context.Context, webhookID uuid.UUID, eventType string, payload map[string]interface{}) error {
	webhook, err := s.db.GetWebhook(ctx, webhookID)
	if err != nil {
		return err
	}

	if !webhook.Active {
		return nil
	}

	// Check if webhook is subscribed to this event
	subscribed := false
	for _, event := range webhook.Events {
		if event == eventType || event == "*" {
			subscribed = true
			break
		}
	}

	if !subscribed {
		return nil
	}

	// Create webhook event record
	webhookEvent := &models.WebhookEvent{
		ID:          uuid.New(),
		WebhookID:   webhookID,
		EventType:   eventType,
		Payload:     payload,
		Status:      "pending",
		Attempts:    0,
		MaxAttempts: 5,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.CreateWebhookEvent(ctx, webhookEvent); err != nil {
		return err
	}

	// Attempt to send webhook
	return s.sendWebhookEvent(ctx, webhook, webhookEvent)
}

// sendWebhookEvent attempts to send webhook event
func (s *WebhookService) sendWebhookEvent(ctx context.Context, webhook *models.Webhook, event *models.WebhookEvent) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.URL, bytes.NewBuffer(payload))
	if err != nil {
		return s.recordWebhookEventFailure(ctx, event, http.StatusInternalServerError, "")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-MasterGo-Event", event.EventType)
	req.Header.Set("X-MasterGo-Timestamp", time.Now().Format(time.RFC3339))

	// Add custom headers
	for k, v := range webhook.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return s.recordWebhookEventFailure(ctx, event, 0, err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return s.recordWebhookEventFailure(ctx, event, resp.StatusCode, string(body))
	}

	event.Status = "sent"
	event.ResponseCode = &resp.StatusCode
	event.ResponseBody = string(body)
	event.UpdatedAt = time.Now()
	event.Attempts++

	return s.db.UpdateWebhookEvent(ctx, event)
}

// recordWebhookEventFailure records webhook event failure
func (s *WebhookService) recordWebhookEventFailure(ctx context.Context, event *models.WebhookEvent, statusCode int, errorBody string) error {
	event.Attempts++
	event.Status = "failed"
	if statusCode != 0 {
		event.ResponseCode = &statusCode
	}
	event.ResponseBody = errorBody
	event.UpdatedAt = time.Now()

	// Schedule retry if attempts remain
	if event.Attempts < event.MaxAttempts {
		event.Status = "retrying"
		nextRetry := time.Now().Add(time.Minute * time.Duration(event.Attempts*2)) // Exponential backoff
		event.NextRetryAt = &nextRetry
	}

	return s.db.UpdateWebhookEvent(ctx, event)
}

// RetryFailedWebhooks retries failed webhook events
func (s *WebhookService) RetryFailedWebhooks(ctx context.Context) error {
	events, err := s.db.GetFailedWebhookEvents(ctx)
	if err != nil {
		return err
	}

	for _, event := range events {
		webhook, err := s.db.GetWebhook(ctx, event.WebhookID)
		if err != nil {
			continue
		}

		s.sendWebhookEvent(ctx, webhook, &event)
	}

	return nil
}
