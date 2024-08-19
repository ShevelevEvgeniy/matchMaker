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

func (e *FormedGroupEvent) Handle(_ context.Context, group dto.Group) {
	e.log.Info("formed group event received")

	e.fillingMetrics(&group)

	e.printGroupInfo(group)
}

func (e *FormedGroupEvent) printGroupInfo(group dto.Group) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("\nGroup #%d:\n", group.GroupID))
	builder.WriteString(fmt.Sprintf("Skill: min = %.2f, max = %.2f, avg = %.2f\n", group.MinSkill, group.MaxSkill, group.AvgSkill))
	builder.WriteString(fmt.Sprintf("Latency: min = %.2f, max = %.2f, avg = %.2f\n", group.MinLatency, group.MaxLatency, group.AvgLatency))
	builder.WriteString(fmt.Sprintf("Time Spent in Queue: min = %v, max = %v, avg = %v\n",
		group.MinTimeSpentInQueue, group.MaxTimeSpentInQueue, group.AvgTimeSpentInQueue))

	builder.WriteString("Users: ")
	for i, user := range group.Users {
		if i == 0 {
			builder.WriteString(fmt.Sprintf("%s", user.Name))
			continue
		}

		builder.WriteString(fmt.Sprintf(", %s", user.Name))
	}

	fmt.Println(builder.String())
}

func (e *FormedGroupEvent) fillingMetrics(group *dto.Group) {
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

	for _, user := range group.Users {
		if user.Skill != 0 {
			group.MinSkill = math.Min(group.MinSkill, user.Skill)
			group.MaxSkill = math.Max(group.MaxSkill, user.Skill)
			totalSkill += user.Skill
		}

		if user.Latency != 0 {
			group.MinLatency = math.Min(group.MinLatency, user.Latency)
			group.MaxLatency = math.Max(group.MaxLatency, user.Latency)
			totalLatency += user.Latency
		}

		timeSpentInQueue := time.Since(user.SearchStartTime)
		if timeSpentInQueue < group.MinTimeSpentInQueue {
			group.MinTimeSpentInQueue = timeSpentInQueue
		}
		if timeSpentInQueue > group.MaxTimeSpentInQueue {
			group.MaxTimeSpentInQueue = timeSpentInQueue
		}
		totalTimeSpentInQueue += timeSpentInQueue
	}

	if len(group.Users) > 0 {
		group.AvgSkill = totalSkill / float64(len(group.Users))
		group.AvgLatency = totalLatency / float64(len(group.Users))
		group.AvgTimeSpentInQueue = totalTimeSpentInQueue / time.Duration(len(group.Users))
	}
}
