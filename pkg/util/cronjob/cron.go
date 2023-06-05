package cronjob

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

type CronFunc func(ctx context.Context) error

type Cron struct {
	robfigCron *cron.Cron
}

func NewCron() *Cron {
	return &Cron{
		robfigCron: cron.New(),
	}
}

// CronHandler wraps cron handler
type CronHandler interface {
	RegisterCron(*Cron)
}

// AddFunc adds a func to the Cron to be run on the given schedule.
func (c *Cron) AddFunc(name, spec string, cmd CronFunc) {
	_, err := c.robfigCron.AddFunc(spec, c.WrapCronFunc(name, cmd))
	if err != nil {
		log.Println(fmt.Sprintf("[Cron] Error registering cron %s: %+v, skipping...", name, err))
		return
	}

	log.Println(fmt.Sprintf("[Cron] Cron %s is successfully registered, spec: %s", name, spec))
}

// WrapCronFunc wraps cron func
func (c *Cron) WrapCronFunc(cronName string, fn CronFunc) func() {
	return func() {
		err := fn(context.Background())
		if err != nil {
			log.Println(fmt.Sprintf("[Cron][%s] Cron execution return error: %+v", cronName, err))
			return
		}

		log.Println(fmt.Sprintf("[Cron][%s] Cron is successfully executed", cronName))
	}
}

// CountEntries returns the number of entries in the Cron.
func (c *Cron) CountEntries() int {
	return len(c.robfigCron.Entries())
}

// Start the cron scheduler in its own goroutine, or no-op if already started.
func (c *Cron) Start() {
	c.robfigCron.Start()
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing.
func (c *Cron) Stop() {
	c.robfigCron.Stop()
}
