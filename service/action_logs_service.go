package service

import (
	"context"
	"ecommerce/entity"
	"ecommerce/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ActionLogService interface {
	GetLogs(ctx context.Context) ([]entity.ActionLog, error)
	GetLogByID(ctx context.Context, id string) (*entity.ActionLog, error)
	ReportSelling(ctx context.Context, reportType string, limit int) ([]map[string]interface{}, error)
	ReportStock(ctx context.Context, reportType string, limit int) ([]entity.Product, error)
	Log(ctx context.Context, actorType string, actorID *uuid.UUID, action, entityType string, entityID interface{}) error
}

type actionLogService struct {
	repo repository.ActionLogRepository
}

func NewActionLogService(repo repository.ActionLogRepository) ActionLogService {
	return &actionLogService{repo}
}

func (s *actionLogService) Log(
	ctx context.Context,
	actorType string,
	actorID *uuid.UUID, 
	action, entityType string,
	entityID interface{},
) error {
	var entityIDStr string
	switch v := entityID.(type) {
	case uuid.UUID:
		entityIDStr = v.String()
	case uint, int, int64:
		entityIDStr = fmt.Sprintf("%v", v)
	case string:
		entityIDStr = v
	default:
		entityIDStr = ""
	}

	log := &entity.ActionLog{
		ID:         uuid.New(),
		ActorType:  actorType,
		ActorID:    actorID, 
		Action:     action,
		EntityType: entityType,
		EntityID:   entityIDStr,
		CreatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, log)
}

func (s *actionLogService) GetLogs(ctx context.Context) ([]entity.ActionLog, error) {
	return s.repo.Getlogs(ctx)
}

func (s *actionLogService) GetLogByID(ctx context.Context, id string) (*entity.ActionLog, error) {
	log, err := s.repo.GetByLogID(ctx, id)
	if err != nil {
		return nil, errors.New("log not found")
	}
	return log, nil
}

func (s *actionLogService) ReportSelling(ctx context.Context, reportType string, limit int) ([]map[string]interface{}, error) {
	return s.repo.ReportSelling(ctx, reportType, limit)
}

func (s *actionLogService) ReportStock(ctx context.Context, reportType string, limit int) ([]entity.Product, error) {
	return s.repo.ReportStock(ctx, reportType, limit)
}
