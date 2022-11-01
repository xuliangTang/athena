package task

import (
	"github.com/robfig/cron/v3"
	"sync"
)

var taskCron *cron.Cron
var taskCronOnce sync.Once

func GetCron() *cron.Cron {
	taskCronOnce.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
	})
	return taskCron
}
