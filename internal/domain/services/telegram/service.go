package telegram

import (
	"context"
	tg_models "daos_core/internal/data/models/api/tg"
	telegram_request "daos_core/internal/data/models/telegram/request"
	"daos_core/internal/data/repositories/bot"
	instance_repo_impl "daos_core/internal/data/repositories/instance"
	tg_repo_impl "daos_core/internal/data/repositories/telegram"
	telegram_adapter "daos_core/internal/external/adapters/telegram"
	"errors"
	"fmt"
	"strings"
)

// woah. its PK shit
const (
	InstanceStatusActive       int64 = 1
	InstanceStatusNotConnected int64 = 2
	InstanceStatusDeleted      int64 = 3
)

type Service interface {
	UpdateStatusByID(ctx context.Context, instanceID int64, statusID int64) (string, error)
	GetBotStatus(ctx context.Context, instanceID int64) (string, error)
	GetBotToken(ctx context.Context, userID int64, botID int64) (string, error)
	//UpdateBotToken(ctx context.Context, newToken string, instanceID int64) error

	GetBotInfo(ctx context.Context, botToken string) ([]byte, error)
	DeleteWebhook(ctx context.Context, requestData telegram_request.SendBotTokenDto) (*tg_models.TypicalResponseDto, error)
	SetWebhook(ctx context.Context, requestData telegram_request.SendBotTokenDto) (*tg_models.TypicalResponseDto, error)
	GetWebhookInfo(ctx context.Context, botToken string) ([]byte, error)
}

type impl struct {
	TelegramRepo tg_repo_impl.Repository
	InstanceRepo instance_repo_impl.Repository
	BotTokenRepo bot.Repository
	Adapter      telegram_adapter.Adapter
}

func NewService(
	t_r tg_repo_impl.Repository,
	i_r instance_repo_impl.Repository,
	b_r bot.Repository,
	t_a telegram_adapter.Adapter,
) Service {
	return &impl{
		TelegramRepo: t_r,
		InstanceRepo: i_r,
		BotTokenRepo: b_r,
		Adapter:      t_a,
	}
}

// func (serv *impl) UpdateBotToken(ctx context.Context, newToken string, instanceID int64) error {
// 	if err := serv.BotTokenRepo.Update(ctx, newToken, instanceID); err != nil {
// 		return fmt.Errorf("TelegramService: UpdateBotToken: %w", err)
// 	}

// 	return nil
// }

func (serv *impl) UpdateStatusByID(
	ctx context.Context,
	instanceID int64,
	statusID int64,
) (string, error) {
	err := serv.InstanceRepo.UpdateStatus(ctx, instanceID, statusID)
	if err != nil {
		return "", err
	}
	return "Status was updated", nil
}

func (serv *impl) GetBotStatus(ctx context.Context, instanceID int64) (string, error) {
	status, err := serv.TelegramRepo.GetStatusByInstance(ctx, instanceID)
	if err != nil {
		return "", err
	}

	if status == "" {
		return "", fmt.Errorf("status is empty")
	}

	return status, errors.New("either accountId or instanceId must be provided")
}

func (serv *impl) GetBotToken(ctx context.Context, userID int64, botID int64) (string, error) {
	token, err := serv.TelegramRepo.GetToken(ctx, userID, botID)
	if err != nil {
		return "", err
	}

	if token == "" {
		return "", fmt.Errorf("token is empty")
	}
	return token, nil
}

func (s *impl) GetBotInfo(ctx context.Context, botToken string) ([]byte, error) {
	if botToken == "" {
		return nil, ErrEmptyToken
	}

	data, err := s.Adapter.GetBotInfo(botToken)
	if err != nil {
		return nil, err
	}

	return data, nil
}

const (
	WebhookIsAlreadySet     = "Webhook is already set"
	WebhookIsAlreadyDeleted = "Webhook is already deleted"
)

func (s *impl) DeleteWebhook(
	ctx context.Context,
	dto telegram_request.SendBotTokenDto,
) (*tg_models.TypicalResponseDto, error) {
	if dto.BotToken == "" {
		return nil, ErrEmptyToken
	}

	answer, err := s.Adapter.DeleteWebhook(dto.BotToken)

	if answer.Ok {
		// error tg case
		if answer.Description == WebhookIsAlreadyDeleted {
			return nil, errors.New(strings.ToLower(WebhookIsAlreadyDeleted))
		} else {
			// if ok

			_, err = s.UpdateStatusByID(ctx, dto.InstanceId, InstanceStatusNotConnected)

			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (s *impl) SetWebhook(ctx context.Context, requestData telegram_request.SendBotTokenDto) (
	*tg_models.TypicalResponseDto, error) {
	if requestData.BotToken == "" {
		return nil, ErrEmptyToken
	}

	data, err := s.Adapter.SetWebhook(requestData.BotToken)
	if err != nil {
		return nil, err
	}

	if data.Ok {
		if data.Description == WebhookIsAlreadySet {
			return nil, ErrEmptyToken
		} else {
			// если все гуд
			_, err := s.UpdateStatusByID(ctx, requestData.InstanceId, InstanceStatusActive)
			if err != nil {
				return nil, err
			}

			// err = s.UpdateBotToken(ctx, requestData.BotToken, requestData.InstanceId)
			// if err != nil {
			// 	return nil, err
			// }
		}
	}

	return data, nil
}

func (s *impl) GetWebhookInfo(ctx context.Context, botToken string) ([]byte, error) {
	if botToken == "" {
		return nil, ErrEmptyToken
	}

	data, err := s.Adapter.GetWebhookInfo(botToken)
	if err != nil {
		return nil, err
	}

	return data, nil
}
