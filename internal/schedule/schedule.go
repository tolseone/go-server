package schedule

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"

	"go-server/internal/models"
	"go-server/pkg/logging"
)

func ScheduleTask() {
	logger := logging.GetLogger()
	logger.Info("Connected logger to Scheduler")

	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Infof("Error creating scheduler: %v\n", err)
	}
	logger.Infof("Scheduler created: %v\n", s)

	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				model.CheckAndDeleteExpiredTokens()
			},
		),
	)
	if err != nil {
		logger.Infof("Error creating job: %v\n", err)
	}

	fmt.Println(j.ID())

	s.Start()

	for {
		time.Sleep(1 * time.Hour)
	}
}
