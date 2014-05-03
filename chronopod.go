package chronopod

import (
	"math/rand"
	"time"
)

const ()

type Chronopod struct {
}

func NewChronopod() *Chronopod {
	return nil
}

// chronoServiceLoop provides the engine for driving repeated time-based events. It assumes a timer given to it already is ticking down
// to when it should begin its interval-based ticking. Thus a long pause or a consistent starting point can be chosen before the time
// period of ticks begins.
//
// 'dura' is the duration between intervals once regular ticking begins.
// 'quit' is the channel used to signal this function to terminate, which will close and nil the quit channel.
// 'catchUpTimer' is an already-started timer that ticks down when to begin regular interval calls to the provided function.
// 'chronoticker' may be a nil ticker, unused until after the 'catchUpTimer' finishes. It then repeatedly triggers after 'dura'
//      length of time until a message is received on 'quit'.
// 'triggerFn' is the function or unit of work spawned in its own goroutine when the 'catchUpTimer' or 'chronoticker' activates.
func chronoServiceLoop(dura time.Duration, quit chan bool, catchUpTimer *time.Timer, chronoticker *time.Ticker, triggerFn func(time.Time)) {
	if chronoticker == nil {
		chronoticker = time.NewTicker(dura)
		chronoticker.Stop()
	}
	for {
		select {
		case _ = <-quit:
			close(quit)
			if chronoticker != nil {
				chronoticker.Stop()
			}
			if catchUpTimer != nil {
				catchUpTimer.Stop()
			}
			quit = nil
			return
		case t := <-(*catchUpTimer).C:
			go triggerFn(t)
			catchUpTimer.Stop()
			chronoticker = time.NewTicker(dura)
		case t := <-chronoticker.C:
			go triggerFn(t)
		}
	}
}

// varyingChronoServiceLoop behaves the same as 'chronoServiceLoop', with the caveat that the tick interval is inconsistent and
// random. The upper and lower limit to a uniform distribution is provided instead of a deterministic value.
//
// 'duraLower' is the lower bound of a uniform distribution of interval times.
// 'duraUpper' is the upper bound of a uniform distribution of interval times.
func varyingChronoServiceLoop(duraLower, duraUpper time.Duration, quit chan bool, catchUpTimer *time.Timer, chronoticker *time.Ticker, triggerFn func(time.Time)) {
	if chronoticker == nil {
		chronoticker = time.NewTicker(duraUpper)
		chronoticker.Stop()
	}
	for {
		select {
		case _ = <-quit:
			close(quit)
			if chronoticker != nil {
				chronoticker.Stop()
			}
			if catchUpTimer != nil {
				catchUpTimer.Stop()
			}
			quit = nil
			return
		case t := <-(*catchUpTimer).C:
			go triggerFn(t)
			catchUpTimer.Stop()
			dura := time.Duration(rand.Int63n(int64(duraUpper-duraLower)) + int64(duraLower))
			chronoticker = time.NewTicker(dura)
		case t := <-chronoticker.C:
			go triggerFn(t)
			dura := time.Duration(rand.Int63n(int64(duraUpper-duraLower)) + int64(duraLower))
			chronoticker = time.NewTicker(dura)
		}
	}
}
