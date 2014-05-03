package chronopod

import (
	"time"
)

type GenericUniformlyVaryingChronoService struct {
	triggerFn    func(time.Time)
	chronoticker *time.Ticker
	catchUpTimer *time.Timer
	quit         chan bool
}

func NewGenericUniformlyVaryingChronoService(triggerFn func(time.Time)) *GenericUniformlyVaryingChronoService {
	return &GenericUniformlyVaryingChronoService{triggerFn, nil, nil, nil}
}

func (c *GenericUniformlyVaryingChronoService) Start(futureTime time.Time, durationLower, durationUpper time.Duration) {
	if c.quit != nil {
		return
	}

	c.quit = make(chan bool)
	curT := time.Now()
	if c.catchUpTimer == nil {
		c.catchUpTimer = time.NewTimer(futureTime.Sub(curT))
	} else {
		c.catchUpTimer.Reset(futureTime.Sub(curT))
	}

	go chronoServiceLoop(durationLower, durationUpper, c.quit, c.catchUpTimer, c.chronoticker, c.triggerFn)
}

func (c *GenericUniformlyVaryingChronoService) Stop() {
	if c.quit != nil {
		c.quit <- true
	}
}
