package larkchunking

import (
	"context"

	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/otel"
	"github.com/BetaGoRobot/BetaGo-Redefine/pkg/xchunk"
	"github.com/BetaGoRobot/go_utils/reflecting"
)

var M *xchunk.Management

func Init() {
	M = xchunk.NewManagement()
	ctx, span := otel.T().Start(context.Background(), reflecting.GetCurrentFunc())
	defer span.End()
	M.StartBackgroundCleaner(ctx)
}
