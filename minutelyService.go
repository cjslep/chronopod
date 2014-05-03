package chronopod

import (
	"errors"
	"time"
)

type MinutelyChronoService struct {
	GenericChronoService
	secondsToLaunch int
}

func NewMinutelyChronoService(launchAtSeconds int, triggerFn func(time.Time)) (*MinutelyChronoService, error) {
	if launchAtSeconds < 0 || launchAtSeconds >= 60 {
		return nil, errors.New("MinutelyChronoService error: Can only launch service within 0-59 seconds inclusive")
	}
	return &MinutelyChronoService{GenericChronoService{triggerFn, nil, nil, nil}, launchAtSeconds}, nil
}

func (c *MinutelyChronoService) Start() {
	curT := time.Now()
	future := time.Date(curT.Year(), curT.Month(), curT.Day(), curT.Hour(), curT.Minute(), c.secondsToLaunch, 0, time.Local)
	if curT.After(future) {
		future = future.Add(time.Minute)
	}
	c.GenericChronoService.Start(future, time.Minute)
}
