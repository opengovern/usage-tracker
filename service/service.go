package service

import (
	"context"

	shared_entities "github.com/opengovern/og-util/pkg/api/shared-entities"
	"github.com/opengovern/usage-tracker/config"
	"github.com/opengovern/usage-tracker/db/model"
	"github.com/opengovern/usage-tracker/db/repo"
	"go.uber.org/zap"
)

type InformationService struct {
	cfg    config.InformationConfig
	logger *zap.Logger

	csmpUsageRepo repo.CspmUsageRepo
}

func NewInformationService(cfg config.InformationConfig, logger *zap.Logger, cspmUsageRepo repo.CspmUsageRepo) *InformationService {
	return &InformationService{
		cfg:           cfg,
		logger:        logger.Named("information-service"),
		csmpUsageRepo: cspmUsageRepo,
	}
}

func (s *InformationService) RecordUsage(ctx context.Context, req shared_entities.CspmUsageRequest) error {

	m := model.CspmUsage{
		InstallId:            req.InstallId,
		GatherTimestamp:      req.GatherTimestamp,
		Hostname:             req.Hostname,
		NumberOfUsers:        req.NumberOfUsers,
		IntegrationTypeCount: req.IntegrationTypeCount,
	}

	if err := s.csmpUsageRepo.Create(ctx, &m); err != nil {
		s.logger.Error("failed to create cspm usage", zap.Error(err))
		return err
	}
	return nil
}
