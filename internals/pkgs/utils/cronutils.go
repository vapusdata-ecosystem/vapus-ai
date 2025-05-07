package dmutils

import (
	"github.com/robfig/cron/v3"
)

func IsValidCron(expr string) bool {
	_, err := cron.ParseStandard(expr)
	return err == nil
}
