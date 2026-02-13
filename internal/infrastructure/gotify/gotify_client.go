package gotify

import (
	"context"
	"net/http"
	"net/url"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/config"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/logs"
	"github.com/go-openapi/runtime"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"go.uber.org/zap"
)

var (
	tokenParsed         runtime.ClientAuthInfoWriter
	DefaultGotifyClient *client.GotifyREST
)

func Init() {
	config := config.Get().GotifyConfig
	gotifyURLParsed, err := url.Parse(config.URL)
	if err != nil {
		panic("error parsing url for gotify" + err.Error())
	}
	DefaultGotifyClient = gotify.NewClient(gotifyURLParsed, &http.Client{})
	tokenParsed = auth.TokenAuth(config.ApplicationToken)
}

func SendMessage(ctx context.Context, title, msg string, priority int) {
	ctx, span := otel.T().Start(ctx, "SendMessage")
	defer span.End()
	logs.L().Ctx(ctx).Info("SendMessage...", zap.String("traceID", span.SpanContext().TraceID().String()))

	if title == "" {
		title = "BetaGo Notification"
	}
	title = "[" + config.Get().BaseInfo.RobotName + "]" + title
	params := message.NewCreateMessageParams()
	params.Body = &models.MessageExternal{
		Title:    title,
		Message:  msg,
		Priority: priority,
		Extras: map[string]interface{}{
			"client::display": map[string]string{"contentType": "text/markdown"},
		},
	}

	_, err := DefaultGotifyClient.Message.CreateMessage(params, tokenParsed)
	if err != nil {
		logs.L().Ctx(ctx).Error("Could not send message", zap.Error(err))
		return
	}
	logs.L().Ctx(ctx).Info("Gotify Message Sent!")
}
