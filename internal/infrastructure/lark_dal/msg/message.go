package msg

import (
	"context"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/logs"
	"github.com/BetaGoRobot/go_utils/reflecting"
	"github.com/bytedance/sonic"
	"github.com/kevinmatthe/zaplog"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// PreGetTextMsg 获取消息内容
//
//	@param ctx
//	@param event
//	@return string
func PreGetTextMsg(ctx context.Context, event *larkim.P2MessageReceiveV1) string {
	ctx, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	defer span.End()
	return GetContentFromTextMsg(*event.Event.Message.Content)
}

func GetContentFromTextMsg(s string) string {
	msgMap := make(map[string]interface{})
	err := sonic.UnmarshalString(s, &msgMap)
	if err != nil {
		logs.L().Error("repeatMessage", zaplog.Error(err))
		return ""
	}
	if text, ok := msgMap["text"]; ok {
		s = text.(string)
	}
	return s
}
