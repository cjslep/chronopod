package chronopod

import (
	"time"
)

type GenericChronoService struct {
	triggerFn    func(time.Time)
	chronoticker *time.Ticker
	catchUpTimer *time.Timer
	quit         chan bool
}

func NewGenericChronoService(triggerFn func(time.Time)) *GenericChronoService {
	return &GenericChronoService{triggerFn, nil, nil, nil}
}

func (c *GenericChronoService) Start(futureTime time.Time, duration time.Duration) {
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

	go chronoServiceLoop(duration, c.quit, c.catchUpTimer, c.chronoticker, c.triggerFn)
}

func (c *GenericChronoService) Stop() {
	if c.quit != nil {
		c.quit <- true
	}
}
