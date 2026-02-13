package larkmsg

import (
	"context"
	"errors"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/logs"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/utils"
	"github.com/BetaGoRobot/go_utils/reflecting"
	"github.com/kevinmatthe/zaplog"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func AddReaction(ctx context.Context, reactionType, msgID string) (reactionID string, err error) {
	_, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	span.SetAttributes(attribute.Key("msgID").String(msgID))
	defer span.End()

	req := larkim.NewCreateMessageReactionReqBuilder().Body(larkim.NewCreateMessageReactionReqBodyBuilder().ReactionType(larkim.NewEmojiBuilder().EmojiType(reactionType).Build()).Build()).MessageId(msgID).Build()
	resp, err := lark_dal.Client().Im.V1.MessageReaction.Create(ctx, req)
	if err != nil {
		logs.L().Ctx(ctx).Error("AddReaction", zaplog.Error(err))
		return "", err
	}
	if !resp.Success() {
		logs.L().Ctx(ctx).Error("AddReaction", zaplog.String("Error", resp.Error()))
		return "", errors.New(resp.Error())
	}
	utils.AddTrace2DB(ctx, msgID)
	return *resp.Data.ReactionId, err
}

func AddReactionAsync(ctx context.Context, reactionType, msgID string) (err error) {
	_, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	span.SetAttributes(attribute.Key("msgID").String(msgID))
	defer span.End()
	defer func() { span.RecordError(err) }()

	req := larkim.NewCreateMessageReactionReqBuilder().Body(larkim.NewCreateMessageReactionReqBodyBuilder().ReactionType(larkim.NewEmojiBuilder().EmojiType(reactionType).Build()).Build()).MessageId(msgID).Build()
	go func() {
		resp, err := lark_dal.Client().Im.V1.MessageReaction.Create(ctx, req)
		if err != nil {
			logs.L().Ctx(ctx).Error("AddReaction", zap.Error(err))
			return
		}
		if !resp.Success() {
			logs.L().Ctx(ctx).Error("AddReaction", zap.String("respError", resp.Error()))
			return
		}
		utils.AddTrace2DB(ctx, msgID)
	}()
	return nil
}

func RemoveReactionAsync(ctx context.Context, reactionID, msgID string) (err error) {
	_, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	span.SetAttributes(attribute.Key("msgID").String(msgID))
	defer span.End()
	defer func() { span.RecordError(err) }()
	req := larkim.NewDeleteMessageReactionReqBuilder().MessageId(msgID).ReactionId(reactionID).Build()
	go func() {
		resp, err := lark_dal.Client().Im.V1.MessageReaction.Delete(ctx, req)
		if err != nil {
			logs.L().Ctx(ctx).Error("RemoveReaction", zap.Error(err))
			return
		}
		if !resp.Success() {
			logs.L().Ctx(ctx).Error("RemoveReaction", zap.String("respError", resp.Error()))
			err = errors.New(resp.Error())
			return
		}
		utils.AddTrace2DB(ctx, msgID)
	}()
	return
}
