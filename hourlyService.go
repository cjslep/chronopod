package chronopod

import (
	"errors"
	"time"
)

type HourlyChronoService struct {
	GenericChronoService
	minutesToLaunch int
}

func NewHourlyChronoService(launchAtMinutes int, triggerFn func(time.Time)) (*HourlyChronoService, error) {
	if launchAtMinutes < 0 || launchAtMinutes >= 60 {
		return nil, errors.New("HourlyChronoService error: Can only launch service within 0-59 minutes inclusive")
	}
	return &HourlyChronoService{GenericChronoService{triggerFn, nil, nil, nil}, launchAtMinutes}, nil
}

func (c *HourlyChronoService) Start() {
	curT := time.Now()
	future := time.Date(curT.Year(), curT.Month(), curT.Day(), curT.Hour(), c.minutesToLaunch, 0, 0, time.Local)
	if curT.After(future) {
		future = future.Add(time.Hour)
	}
	c.GenericChronoService.Start(future, time.Hour)
}
