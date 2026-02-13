package larkmsg

import (
	"context"
	"errors"
	"fmt"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/lark_dal/larkmsg/larktpl"
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/go_utils/reflecting"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.opentelemetry.io/otel/attribute"
)

// PatchCard to be filled PatchCard
//
//	@param ctx context.Context
//	@param cardContent *templates.TemplateCardContent
//	@param msgID string
//	@return err error
//	@author kevinmatthe
//	@update 2025-06-05 13:23:46
func PatchCard(ctx context.Context, cardContent *larktpl.TemplateCardContent, msgID string) (err error) {
	_, span := otel.T().Start(ctx, reflecting.GetCurrentFunc())
	span.SetAttributes(attribute.Key("msgID").String(msgID))
	for k, v := range cardContent.Data.TemplateVariable {
		span.SetAttributes(attribute.Key(k).String(fmt.Sprintf("%v", v)))
	}
	defer span.End()
	defer func() { span.RecordError(err) }()
	resp, err := lark_dal.Client().Im.V1.Message.Patch(
		ctx, larkim.NewPatchMessageReqBuilder().
			MessageId(msgID).
			Body(
				larkim.NewPatchMessageReqBodyBuilder().
					Content(cardContent.String()).
					Build(),
			).
			Build(),
	)
	if err != nil {
		return
	}
	if !resp.Success() {
		return errors.New(resp.Error())
	}
	return
}
