package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/db/model"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/db/query"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal/larkmsg"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal/larkmsg/larktpl"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/logs"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/utils"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/xhandler"

	"github.com/BetaGoRobot/go_utils/reflecting"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

// WordAddHandler to be filled
//
//	@param ctx context.Context
//	@param data *larkim.P2MessageReceiveV1
//	@param args ...string
//	@return error
//	@author heyuhengmatt
//	@update 2024-08-06 08:27:09
func WordAddHandler(ctx context.Context, data *larkim.P2MessageReceiveV1, metaData *xhandler.BaseMetaData, args ...string) (err error) {
	ctx, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	span.SetAttributes(attribute.Key("event").String(larkcore.Prettify(data)))
	defer span.End()
	defer func() { span.RecordError(err) }()

	if len(args) < 2 {
		return errors.ErrUnsupported
	}
	argMap, _ := parseArgs(args...)
	logs.L().Ctx(ctx).Info("args", zap.Any("args", argMap))

	word, ok := argMap["word"]
	if !ok {
		return errors.New("word is required")
	}
	rate, ok := argMap["rate"]
	if !ok {
		return errors.New("rate is required")
	}

	ChatID := *data.Event.Message.ChatId
	return query.Q.RepeatWordsRateCustom.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&model.RepeatWordsRateCustom{
		GuildID: ChatID,
		Word:    word,
		Rate:    int64(utils.MustAtoI(rate)),
	})
}

// WordGetHandler to be filled
//
//	@param ctx context.Context
//	@param data *larkim.P2MessageReceiveV1
//	@param args ...string
//	@return error
//	@author heyuhengmatt
//	@update 2024-08-06 08:27:07
func WordGetHandler(ctx context.Context, data *larkim.P2MessageReceiveV1, metaData *xhandler.BaseMetaData, args ...string) (err error) {
	ctx, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	span.SetAttributes(attribute.Key("event").String(larkcore.Prettify(data)))
	defer span.End()
	defer func() { span.RecordError(err) }()
	argMap, _ := parseArgs(args...)
	logs.L().Ctx(ctx).Info("args", zap.Any("args", argMap))
	ChatID := *data.Event.Message.ChatId

	lines := make([]map[string]string, 0)
	ins := query.Q.RepeatWordsRateCustom
	resListCustom, err := ins.WithContext(ctx).
		Where(ins.GuildID.Eq(ChatID)).
		Find()
	if err != nil {
		return err
	}
	for _, res := range resListCustom {
		if res.GuildID == ChatID {
			lines = append(lines, map[string]string{
				"title1": "Custom",
				"title2": res.Word,
				"title3": strconv.Itoa(int(res.Rate)),
			})
		}
	}
	ins2 := query.Q.RepeatWordsRate
	resListGlobal, err := ins2.WithContext(ctx).
		Find()
	if err != nil {
		return err
	}
	for _, res := range resListGlobal {
		lines = append(lines, map[string]string{
			"title1": "Global",
			"title2": res.Word,
			"title3": strconv.Itoa(int(res.Rate)),
		})
	}
	cardContent := larktpl.NewCardContent(
		ctx,
		larktpl.ThreeColSheetTemplate,
	).
		AddVariable("title1", "Scope").
		AddVariable("title2", "Keyword").
		AddVariable("title3", "Rate").
		AddVariable("table_raw_array_1", lines)

	err = larkmsg.ReplyCard(ctx, cardContent, *data.Event.Message.MessageId, "_wordGet", false)
	if err != nil {
		return err
	}
	return nil
}
