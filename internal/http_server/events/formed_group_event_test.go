package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"matchMaker/internal/dto"
	"matchMaker/internal/storage/postgres/repository/models"
)

func TestFillingMetrics(t *testing.T) {
	event := NewFormedGroupEvent(zap.NewNop())

	group := dto.Group{
		Users: []models.User{
			{ID: 1, Skill: 5.0, Latency: 2.0, SearchStartTime: time.Now().Add(-10 * time.Minute)},
			{ID: 2, Skill: 8.0, Latency: 10.0, SearchStartTime: time.Now().Add(-20 * time.Minute)},
		},
	}

	event.fillingMetrics(&group)

	assert.Equal(t, 5.0, group.MinSkill)
	assert.Equal(t, 8.0, group.MaxSkill)
	assert.Equal(t, 6.5, group.AvgSkill)
	assert.Equal(t, 2.0, group.MinLatency)
	assert.Equal(t, 10.0, group.MaxLatency)
	assert.Equal(t, 6.0, group.AvgLatency)

	const tolerance = time.Millisecond * 10

	assert.InDelta(t, 15*time.Minute, group.AvgTimeSpentInQueue, float64(tolerance))
}
