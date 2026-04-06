package chat

import (
	"context"
	"daos_core/internal/data/repositories/account"
	account_repository "daos_core/internal/data/repositories/account"
	amo_auth "daos_core/internal/data/repositories/amo_tokens"
	"daos_core/internal/data/repositories/chat"
	dto "daos_core/internal/domain/dto/chat"
	chat_adapter "daos_core/internal/external/adapters/chat"
)

type Service interface {
	Create(ctx context.Context, request dto.CreateChatDTO) error
}

type impl struct {
	AccountRepository account.Repository
	ChatRepository    chat.Repository
	TokenRepository   amo_auth.Repository
	ChatAdapter       chat_adapter.Adapter
}

func NewChatService(
	a_r account_repository.Repository,
	a chat_adapter.Adapter,
	a_t amo_auth.Repository,
) Service {
	return &impl{
		AccountRepository: a_r,
		ChatAdapter:       a,
		TokenRepository:   a_t,
	}
}

func (s *impl) Create(ctx context.Context, request dto.CreateChatDTO) error {
	// scopeID, err := s.AccountRepository.GetScopeIDByInstance(ctx, request.UserID)
	// if err != nil {
	// 	return err
	// }

	// amoToken, err := s.TokenRepository.GetAccessToken(ctx, strconv.Itoa(request.UserID))
	// if err != nil {
	// 	return err
	// }

	// conversationID := uuid.NewString()

	// err = s.ChatAdapter.Create(chat_models.CreateInput{
	// 	Token:          amoToken,
	// 	ScopeID:        scopeID,
	// 	UserID:         request.UserID,
	// 	Referer:        request.Referer,
	// 	ConversationID: conversationID,
	// })

	// if err != nil {
	// 	return err
	// }

	// err = s.ChatRepository.Save(ctx, models.CreateInput{
	// 	ConversationID: conversationID,
	// 	InstansceID:    request.InstanceID,
	// 	UserID:         request.UserID,
	// })

	// if err != nil {
	// 	return err
	// }

	return nil
}
