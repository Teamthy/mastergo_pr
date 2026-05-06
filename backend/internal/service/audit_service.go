package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

type AuditLogService struct {
	db *database.Database
}

func NewAuditLogService(db *database.Database) *AuditLogService {
	return &AuditLogService{db: db}
}

// LogAction logs a user action
func (s *AuditLogService) LogAction(ctx context.Context, userID *uuid.UUID, action, resourceType, resourceID string, oldValues, newValues interface{}, ipAddress, userAgent string) error {
	oldValuesJSON, _ := json.Marshal(oldValues)
	newValuesJSON, _ := json.Marshal(newValues)

	oldMap := make(map[string]interface{})
	newMap := make(map[string]interface{})

	json.Unmarshal(oldValuesJSON, &oldMap)
	json.Unmarshal(newValuesJSON, &newMap)

	auditLog := &models.AuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		OldValues:    oldMap,
		NewValues:    newMap,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Status:       "success",
		CreatedAt:    time.Now(),
	}

	return s.db.CreateAuditLog(ctx, auditLog)
}

// LogFailedAction logs a failed action
func (s *AuditLogService) LogFailedAction(ctx context.Context, userID *uuid.UUID, action, resourceType, resourceID, errorMessage string, ipAddress, userAgent string) error {
	auditLog := &models.AuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Status:       "failure",
		ErrorMessage: errorMessage,
		CreatedAt:    time.Now(),
	}

	return s.db.CreateAuditLog(ctx, auditLog)
}

// GetAuditLogs retrieves audit logs for a user
func (s *AuditLogService) GetAuditLogs(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.AuditLog, int64, error) {
	logs, total, err := s.db.GetAuditLogsByUserID(ctx, userID, limit, offset)
	return logs, total, err
}

// GetAuditLogsForResource retrieves audit logs for a specific resource
func (s *AuditLogService) GetAuditLogsForResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]models.AuditLog, error) {
	return s.db.GetAuditLogsByResource(ctx, resourceType, resourceID, limit, offset)
}

// CriticalAction logs a critical action for immediate notification
func (s *AuditLogService) CriticalAction(ctx context.Context, userID *uuid.UUID, action string, ipAddress, userAgent string) error {
	auditLog := &models.AuditLog{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    "success",
		CreatedAt: time.Now(),
	}

	if err := s.db.CreateAuditLog(ctx, auditLog); err != nil {
		return err
	}

	// TODO: Send alert notification for critical actions
	fmt.Printf("CRITICAL ACTION: %s by user %v from %s\n", action, userID, ipAddress)

	return nil
}

// GetActionSummary retrieves summary of actions in time range
func (s *AuditLogService) GetActionSummary(ctx context.Context, userID uuid.UUID, days int) (map[string]int, error) {
	return s.db.GetAuditActionSummary(ctx, userID, days)
}
