package events

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"go.uber.org/zap"
	"matchMaker/internal/dto"
)

type FormedGroupEvent struct {
	log *zap.Logger
}

func NewFormedGroupEvent(log *zap.Logger) *FormedGroupEvent {
	return &FormedGroupEvent{
		log: log,
	}
}

func (e *FormedGroupEvent) Handle(_ context.Context, msg Message) {
	e.log.Info("formed group event received")

	group := msg.Value.(dto.Group)

	e.fillingMetrics(group)

	e.printGroupInfo(group)
}

func (e *FormedGroupEvent) printGroupInfo(group dto.Group) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Group #%d:\n", group.GroupID))
	builder.WriteString(fmt.Sprintf("Min Skill: %.2f, Max Skill: %.2f, Avg Skill: %.2f\n", group.MinSkill, group.MaxSkill, group.AvgSkill))
	builder.WriteString(fmt.Sprintf("Min Latency: %.2f, Max Latency: %.2f, Avg Latency: %.2f\n", group.MinLatency, group.MaxLatency, group.AvgLatency))
	builder.WriteString(fmt.Sprintf("Min Time Spent in Queue: %v, Max Time Spent in Queue: %v, Avg Time Spent in Queue: %v\n",
		group.MinTimeSpentInQueue, group.MaxTimeSpentInQueue, group.AvgTimeSpentInQueue))

	builder.WriteString("Users:\n")
	for _, user := range group.Users {
		builder.WriteString(fmt.Sprintf("- %s\n", user.Name))
	}

	e.log.Info(builder.String())
}

func (e *FormedGroupEvent) fillingMetrics(group dto.Group) {
	if len(group.Users) == 0 {
		return
	}

	var totalSkill float64
	var totalLatency float64
	var totalTimeSpentInQueue time.Duration

	group.MinSkill = math.MaxFloat64
	group.MaxSkill = -math.MaxFloat64
	group.MinLatency = math.MaxFloat64
	group.MaxLatency = -math.MaxFloat64
	group.MinTimeSpentInQueue = time.Duration(math.MaxInt64)
	group.MaxTimeSpentInQueue = 0

	for i := range group.Users {
		user := group.Users[i]

		group.MinSkill = math.Min(group.MinSkill, user.Skill)
		group.MaxSkill = math.Max(group.MaxSkill, user.Skill)
		totalSkill += user.Skill

		group.MinLatency = math.Min(group.MinLatency, user.Latency)
		group.MaxLatency = math.Max(group.MaxLatency, user.Latency)
		totalLatency += user.Latency

		timeSpentInQueue := time.Since(user.SearchStartTime)
		if timeSpentInQueue < group.MinTimeSpentInQueue {
			group.MinTimeSpentInQueue = timeSpentInQueue
		}
		if timeSpentInQueue > group.MaxTimeSpentInQueue {
			group.MaxTimeSpentInQueue = timeSpentInQueue
		}
		totalTimeSpentInQueue += timeSpentInQueue
	}

	group.AvgSkill = totalSkill / float64(len(group.Users))
	group.AvgLatency = totalLatency / float64(len(group.Users))
	group.AvgTimeSpentInQueue = totalTimeSpentInQueue / time.Duration(len(group.Users))
}
