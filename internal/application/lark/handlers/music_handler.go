package handlers

import (
	"context"
	"errors"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/config"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal/larkmsg"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal/larkmsg/larktpl"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/neteaseapi"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/utils"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/xhandler"
	"github.com/BetaGoRobot/go_utils/reflecting"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.opentelemetry.io/otel/attribute"
)

func MusicSearchHandler(ctx context.Context, data *larkim.P2MessageReceiveV1, metaData *xhandler.BaseMetaData, args ...string) (err error) {
	ctx, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	span.SetAttributes(attribute.Key("event").String(larkcore.Prettify(data)))
	defer span.End()
	defer func() { span.RecordError(err) }()

	argsMap, input := parseArgs(args...)
	searchType, ok := argsMap["type"]
	if !ok {
		// 兼容简易搜索
		searchType = "song"
	}

	keywords := []string{input}

	var cardContent *larktpl.TemplateCardContent
	if searchType == "album" {
		albumList, err := neteaseapi.NetEaseGCtx.SearchAlbumByKeyWord(ctx, keywords...)
		if err != nil {
			return err
		}

		cardContent, err = neteaseapi.BuildMusicListCard(ctx,
			albumList,
			neteaseapi.MusicItemTransAlbum,
			neteaseapi.CommentTypeAlbum,
			keywords...,
		)
		if err != nil {
			return err
		}
	} else if searchType == "artist" {
	} else if searchType == "playlist" {
	} else if searchType == "song" {
		musicList, err := neteaseapi.NetEaseGCtx.SearchMusicByKeyWord(ctx, keywords...)
		if err != nil {
			return err
		}
		cardContent, err = neteaseapi.BuildMusicListCard(ctx,
			musicList,
			neteaseapi.MusicItemNoTrans,
			neteaseapi.CommentTypeSong,
			keywords...,
		)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Unknown search type")
	}

	err = larkmsg.ReplyCard(ctx, cardContent, *data.Event.Message.MessageId, "_musicSearch", utils.GetIfInthread(ctx, metaData, config.Get().NeteaseMusicConfig.MusicCardInThread))
	if err != nil {
		return err
	}
	return
}

// func init() {
// 	params := tools.NewParameters("object").
// 		AddProperty("keywords", &tools.Property{
// 			Type:        "string",
// 			Description: "音乐搜索的关键词, 多个关键词之间用空格隔开",
// 		}).AddRequired("keywords")
// 	fcu := tools.NewFunctionCallUnit().
// 		Name("music_search").Desc("根据输入的关键词搜索相关的音乐并发送卡片").Params(params).Func(musicSearchWrap)
// 	tools.M().Add(fcu)
// }

// func musicSearchWrap(ctx context.Context, meta *tools.FunctionCallMeta, args string) (any, error) {
// 	s := struct {
// 		Keywords string `json:"keywords"`
// 	}{}
// 	err := utils.UnmarshallStringPre(args, &s)
// 	if err != nil {
// 		return nil, err
// 	}
// 	metaData := xhandler.NewBaseMetaDataWithChatIDUID(ctx, meta.ChatID, meta.UserID)
// 	return "执行成功", MusicSearchHandler(ctx, meta.LarkData, metaData, s.Keywords)
// }
