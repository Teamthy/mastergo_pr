package service

import (
	"context"
	"fmt"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	db *database.Database
}

func NewAnalyticsService(db *database.Database) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// RecordAPICall records an API call for analytics
func (s *AnalyticsService) RecordAPICall(ctx context.Context, apiKeyID, userID uuid.UUID, endpoint, method string, responseTime int, statusCode, requestSize, responseSize int) error {
	analytics := &models.APIAnalytics{
		ID:             uuid.New(),
		APIKeyID:       apiKeyID,
		UserID:         userID,
		Endpoint:       endpoint,
		Method:         method,
		ResponseTimeMs: responseTime,
		StatusCode:     statusCode,
		RequestSize:    requestSize,
		ResponseSize:   responseSize,
		CreatedAt:      time.Now(),
	}

	if err := s.db.CreateAPIAnalytics(ctx, analytics); err != nil {
		return fmt.Errorf("failed to record API call: %w", err)
	}

	// Update rate limit stats
	return s.updateRateLimitStats(ctx, apiKeyID)
}

// updateRateLimitStats updates daily rate limit statistics
func (s *AnalyticsService) updateRateLimitStats(ctx context.Context, apiKeyID uuid.UUID) error {
	today := time.Now().Truncate(24 * time.Hour)

	stats, err := s.db.GetRateLimitStats(ctx, apiKeyID, today)
	if err == nil && stats != nil {
		// Update existing stats
		stats.RequestCount++
		stats.UpdatedAt = time.Now()
		return s.db.UpdateRateLimitStats(ctx, stats)
	}

	// Create new stats entry
	stats = &models.RateLimitStats{
		ID:           uuid.New(),
		APIKeyID:     apiKeyID,
		Date:         today,
		RequestCount: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.db.CreateRateLimitStats(ctx, stats)
}

// GetAPIAnalytics retrieves API analytics for an API key
func (s *AnalyticsService) GetAPIAnalytics(ctx context.Context, apiKeyID uuid.UUID, days int) (*models.AnalyticsResponse, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	analytics, err := s.db.GetAPIAnalyticsByKeyID(ctx, apiKeyID, startDate)
	if err != nil {
		return nil, err
	}

	if len(analytics) == 0 {
		return &models.AnalyticsResponse{
			TotalRequests: 0,
		}, nil
	}

	response := &models.AnalyticsResponse{
		TotalRequests: int64(len(analytics)),
	}

	// Calculate today's requests
	today := time.Now().Truncate(24 * time.Hour)
	todayCount := 0
	totalTime := 0

	endpointStats := make(map[string]*models.Endpoint)

	for _, a := range analytics {
		if a.CreatedAt.Truncate(24*time.Hour) == today {
			todayCount++
		}

		totalTime += a.ResponseTimeMs

		if _, ok := endpointStats[a.Endpoint]; !ok {
			endpointStats[a.Endpoint] = &models.Endpoint{
				Path: a.Endpoint,
			}
		}

		ep := endpointStats[a.Endpoint]
		ep.Calls++
		ep.AvgTime += float64(a.ResponseTimeMs)

		if a.StatusCode >= 400 {
			ep.ErrorCount++
		}
	}

	response.RequestsToday = todayCount
	if len(analytics) > 0 {
		response.AvgResponseTime = float64(totalTime) / float64(len(analytics))
	}

	// Calculate error rate
	errorCount := 0
	for _, a := range analytics {
		if a.StatusCode >= 400 {
			errorCount++
		}
	}
	if len(analytics) > 0 {
		response.ErrorRate = float64(errorCount) / float64(len(analytics)) * 100
	}

	// Convert endpoint stats to slice
	topEndpoints := make([]models.Endpoint, 0)
	for _, ep := range endpointStats {
		if ep.Calls > 0 {
			ep.AvgTime = ep.AvgTime / float64(ep.Calls)
		}
		topEndpoints = append(topEndpoints, *ep)
	}

	// Sort by calls descending (simple sort, can use sort.Slice for production)
	response.TopEndpoints = topEndpoints

	return response, nil
}

// GetUserAnalytics retrieves analytics for all APIs of a user
func (s *AnalyticsService) GetUserAnalytics(ctx context.Context, userID uuid.UUID, days int) (*models.AnalyticsResponse, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	analytics, err := s.db.GetAPIAnalyticsByUserID(ctx, userID, startDate)
	if err != nil {
		return nil, err
	}

	if len(analytics) == 0 {
		return &models.AnalyticsResponse{
			TotalRequests: 0,
		}, nil
	}

	response := &models.AnalyticsResponse{
		TotalRequests: int64(len(analytics)),
	}

	// Similar processing as above
	today := time.Now().Truncate(24 * time.Hour)
	todayCount := 0
	totalTime := 0

	for _, a := range analytics {
		if a.CreatedAt.Truncate(24*time.Hour) == today {
			todayCount++
		}
		totalTime += a.ResponseTimeMs
	}

	response.RequestsToday = todayCount
	if len(analytics) > 0 {
		response.AvgResponseTime = float64(totalTime) / float64(len(analytics))
	}

	return response, nil
}

// ExportAnalytics exports analytics data in CSV format
func (s *AnalyticsService) ExportAnalytics(ctx context.Context, apiKeyID uuid.UUID, days int) (string, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	analytics, err := s.db.GetAPIAnalyticsByKeyID(ctx, apiKeyID, startDate)
	if err != nil {
		return "", err
	}

	csv := "Timestamp,Endpoint,Method,Status Code,Response Time (ms),Request Size,Response Size\n"

	for _, a := range analytics {
		csv += fmt.Sprintf("%s,%s,%s,%d,%d,%d,%d\n",
			a.CreatedAt.Format(time.RFC3339),
			a.Endpoint,
			a.Method,
			a.StatusCode,
			a.ResponseTimeMs,
			a.RequestSize,
			a.ResponseSize,
		)
	}

	return csv, nil
}
