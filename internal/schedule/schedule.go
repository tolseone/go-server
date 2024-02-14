package schedule

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/robfig/cron/v3"

	"go-server/internal/models"
	"go-server/pkg/logging"

)

func ScheduleTask() {
	logger := logging.GetLogger()

	c := cron.New()
	logger.Info("schedule task")

	_, err := c.AddFunc("@every 1s", func() {
		logger.Info("Running scheduled task...")

		model.CheckAndDeleteExpiredTokens()
		logger.Info("Scheduled task executed successfully")
	})

	if err != nil {
		logger.Error("Error scheduling task:", err)
	}

	c.Start()

	defer c.Stop()
}

func ScheduleTask2() {
	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Printf("Error creating scheduler: %v\n", err)
	}

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
		fmt.Printf("Error creating job: %v\n", err)
	}

	fmt.Println(j.ID())


	s.Start()


	select {
	case <-time.After(time.Minute):
	}


	err = s.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down the scheduler: %v\n", err)
	}
}
