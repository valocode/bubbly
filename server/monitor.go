package server

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/valocode/bubbly/integrations"
)

func (s *Server) initMonitoring() error {
	sched := gocron.NewScheduler(time.UTC)

	//
	// Setup the SPDX monitor
	//
	m, err := integrations.NewSPDXMonitor(
		integrations.WithStore(s.store),
	)
	if err != nil {
		return fmt.Errorf("initializing spdx monitoring: %w", err)
	}
	// Run it once before starting the schedule and make sure it's successful
	if err := m.Do(); err != nil {
		return fmt.Errorf("running initial SPDX import: %w", err)
	}
	spdxJob, err := sched.Cron(CronEveryDay).Do(m.Do)
	if err != nil {
		return fmt.Errorf("creating spdx job: %w", err)
	}
	s.bCtx.Logger.Info().Str("scheduledAt", spdxJob.ScheduledAtTime()).Msg("spdx job scheduled")

	//
	// Start the scheduler
	//
	sched.StartAsync()

	return nil
}
