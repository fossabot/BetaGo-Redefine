package lark

import (
	"context"
	"strconv"
	"time"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/application/lark/messages"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/config"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal/larkmsg"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"

	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/logs"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/utils"

	"github.com/BetaGoRobot/go_utils/reflecting"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkapplication "github.com/larksuite/oapi-sdk-go/v3/service/application/v6"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func isOutDated(createTime string) bool {
	stamp, err := strconv.ParseInt(createTime, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Now().Sub(time.UnixMilli(stamp)) > time.Second*10
}

func MessageV2Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) (err error) {
	fn := reflecting.GetCurrentFunc()
	ctx, span := otel.T().Start(ctx, fn)
	defer larkmsg.RecoverMsg(ctx, *event.Event.Message.MessageId)
	span.SetAttributes(attribute.Key("event").String(larkcore.Prettify(event)))
	defer func() { span.RecordError(err) }()

	if isOutDated(*event.Event.Message.CreateTime) {
		return nil
	}
	if *event.Event.Sender.SenderId.OpenId == config.Get().LarkConfig.BotOpenID {
		return nil
	}
	logs.L().Ctx(ctx).Info("Inside the child span for complex handler", zap.String("event", larkcore.Prettify(event)))
	go func() {
		subCtx, span := otel.T().Start(context.Background(), fn+"_RealRun")
		defer span.End()
		span.SetAttributes(attribute.String("msgID", utils.AddrOrNil(event.Event.Message.MessageId)))
		messages.Handler.Clean().WithCtx(subCtx).WithData(event).Run()
	}()

	logs.L().Ctx(ctx).Info("Message event received", zap.String("event", larkcore.Prettify(event)))
	return nil
}

func MessageReactionHandler(ctx context.Context, event *larkim.P2MessageReactionCreatedV1) (err error) {
	return
}

func CardActionHandler(ctx context.Context, cardAction *callback.CardActionTriggerEvent) (resp *callback.CardActionTriggerResponse, err error) {
	return
}

func AuditV6Handler(ctx context.Context, event *larkapplication.P2ApplicationAppVersionAuditV6) (err error) {
	return
}
