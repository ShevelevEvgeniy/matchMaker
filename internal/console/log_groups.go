package console

import (
	"context"

	"go.uber.org/zap"
	"matchMaker/internal/dto"
)

type LogGroups struct {
	log *zap.Logger
}

func NewLogGroups(log *zap.Logger) *LogGroups {
	return &LogGroups{
		log: log,
	}
}

func (l *LogGroups) PrintGroupInfo(ctx context.Context, groups []dto.Group) {
	// TODO: implement me
}
