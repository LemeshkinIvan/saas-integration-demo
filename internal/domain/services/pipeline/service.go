package pipeline

import (
	"context"
	amo_tokens "daos_core/internal/data/repositories/amo_tokens"
	pipeline_adapter "daos_core/internal/external/adapters/pipeline"
)

type Service interface {
	Get(ctx context.Context, accountId string) ([]byte, error)
}

type impl struct {
	AmoTokensRepo   amo_tokens.Repository
	PipelineAdapter pipeline_adapter.Adapter
}

func NewImpl(r amo_tokens.Repository, a pipeline_adapter.Adapter) Service {
	return &impl{AmoTokensRepo: r, PipelineAdapter: a}
}

func (s *impl) Get(ctx context.Context, accountId string) ([]byte, error) {
	token, err := s.AmoTokensRepo.GetAccessToken(ctx, accountId)
	if err != nil {
		return nil, err
	}

	data, err := s.PipelineAdapter.GetPipeline(token)
	if err != nil {
		return nil, err
	}

	return data, nil
}
