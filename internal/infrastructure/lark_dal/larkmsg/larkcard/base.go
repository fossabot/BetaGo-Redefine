package larkcard

import (
	"context"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal/larkmsg/larktpl"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"

	"github.com/BetaGoRobot/go_utils/reflecting"
)

type CardBuilderBase struct {
	Title    string
	SubTitle string
	Content  string
}

func NewCardBuildHelper() *CardBuilderBase {
	return &CardBuilderBase{}
}

func (h *CardBuilderBase) SetTitle(title string) *CardBuilderBase {
	h.Title = title
	return h
}

func (h *CardBuilderBase) SetSubTitle(subTitle string) *CardBuilderBase {
	h.SubTitle = subTitle
	return h
}

func (h *CardBuilderBase) SetContent(text string) *CardBuilderBase {
	h.Content = text
	return h
}

func (h *CardBuilderBase) Build(ctx context.Context) *larktpl.TemplateCardContent {
	ctx, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	defer span.End()
	cardContent := larktpl.NewCardContent(ctx, larktpl.NormalCardReplyTemplate)
	return cardContent.
		AddVariable(
			"title", h.Title,
		).
		AddVariable(
			"subtitle", h.SubTitle,
		).
		AddVariable(
			"content", h.Content,
		)
}
