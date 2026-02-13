package lark_dal

import (
	"context"
	"errors"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/logs"
	"github.com/BetaGoRobot/go_utils/reflecting"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.uber.org/zap"
)

func GetChatName(ctx context.Context, chatID string) (chatName string) {
	ctx, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	defer span.End()

	resp, err := client.Im.V1.Chat.Get(ctx, larkim.NewGetChatReqBuilder().ChatId(chatID).Build())
	if err != nil {
		return
	}
	if !resp.Success() {
		err = errors.New(resp.Error())
		return
	}
	if resp != nil && resp.Data != nil && resp.Data.Name != nil {
		chatName = *resp.Data.Name
	}
	return
}

func GetChatIDFromMsgID(ctx context.Context, msgID string) (chatID string, err error) {
	ctx, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	defer span.End()
	defer func() { span.RecordError(err) }()

	resp := GetMsgFullByID(ctx, msgID)
	if !resp.Success() {
		err = errors.New(resp.Error())
		return
	}
	chatID = *resp.Data.Items[0].ChatId
	return
}

func GetMsgFullByID(ctx context.Context, msgID string) *larkim.GetMessageResp {
	resp, err := client.Im.V1.Message.Get(ctx, larkim.NewGetMessageReqBuilder().MessageId(msgID).Build())
	if err != nil {
		logs.L().Ctx(ctx).Error("GetMsgByID", zap.Error(err))
	}
	if !resp.Success() {
		logs.L().Ctx(ctx).Error("GetMsgByID", zap.String("error", resp.Error()))
	}
	return resp
}
