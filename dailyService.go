package chronopod

import (
	"errors"
	"time"
)

type DailyChronoService struct {
	GenericChronoService
	hourToLaunch int
}

func NewDailyChronoService(launchAtHour int, triggerFn func(time.Time)) (*DailyChronoService, error) {
	if launchAtHour < 0 || launchAtHour >= 60 {
		return nil, errors.New("DailyChronoService error: Can only launch service within 0-23 hours inclusive")
	}
	return &DailyChronoService{GenericChronoService{triggerFn, nil, nil, nil}, launchAtHour}, nil
}

func (c *DailyChronoService) Start() {
	curT := time.Now()
	future := time.Date(curT.Year(), curT.Month(), curT.Day(), c.hourToLaunch, 0, 0, 0, time.Local)
	if curT.After(future) {
		future = future.Add(24 * time.Hour)
	}
	c.GenericChronoService.Start(future, 24*time.Hour)
}
