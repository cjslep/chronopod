package chronopod

import (
	"fmt"
	"testing"
	"time"
)

func testPrinter(t time.Time) {
	fmt.Printf("Test Printer ticked at: %s\n", t)
}

func notifyWhenTicked(sender chan<- time.Time) func(time.Time) {
	return func(t time.Time) {
		fmt.Printf("Got time %s\n", t)
		sender <- t
	}
}

func TestMinutelyChronoServiceConstructor(t *testing.T) {
	_, err := NewMinutelyChronoService(0, testPrinter)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
}

func TestMinutelyChronoServiceLowInvalidConstructor(t *testing.T) {
	_, err := NewMinutelyChronoService(-1, testPrinter)
	if err == nil {
		t.Errorf("Expected error!\n")
	}
}

func TestMinutelyChronoServiceHighInvalidConstructor(t *testing.T) {
	_, err := NewMinutelyChronoService(60, testPrinter)
	if err == nil {
		t.Errorf("Expected error!\n")
	}
}

func TestMinutelyChronoServiceStopNoStart(t *testing.T) {
	receiver := make(chan time.Time)
	c, err := NewMinutelyChronoService(0, notifyWhenTicked(receiver))
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
	c.Stop()
}

func TestMinutelyChronoServiceStopEarly(t *testing.T) {
	receiver := make(chan time.Time)
	c, err := NewMinutelyChronoService(0, notifyWhenTicked(receiver))
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
	c.Start()
	c.Stop()
}

func TestMinutelyChronoService(t *testing.T) {
	receiver := make(chan time.Time)
	fn := notifyWhenTicked(receiver)
	c, err := NewMinutelyChronoService((time.Now().Second()+1)%60, fn)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
	c.Start()
	timeout := time.NewTicker(time.Second * 3)
	select {
	case time := <-receiver:
		t.Logf("Received: %s\n", time)
	case _ = <-timeout.C:
		t.Errorf("Timed out!\n")
	}
	c.Stop()
}

func TestTwoMinutelyChronoService(t *testing.T) {
	receiver := make(chan time.Time)
	fn := notifyWhenTicked(receiver)
	c, err := NewMinutelyChronoService((time.Now().Second()+1)%60, fn)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
	c.Start()
	timeout := time.NewTicker(time.Second * 63)
	count := 0
OuterLoop:
	for {
		select {
		case time := <-receiver:
			count++
			t.Logf("Received: %s\nCount: %v\n", time, count)
			if count == 2 {
				break OuterLoop
			}
		case _ = <-timeout.C:
			t.Errorf("Timed out!\n")
			break OuterLoop
		}
	}
	c.Stop()
}

func TestThreeMinutelyChronoService(t *testing.T) {
	receiver := make(chan time.Time)
	fn := notifyWhenTicked(receiver)
	c, err := NewMinutelyChronoService((time.Now().Second()+1)%60, fn)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
	c.Start()
	timeout := time.NewTicker(time.Second * 123)
	count := 0
OuterLoop:
	for {
		select {
		case time := <-receiver:
			count++
			t.Logf("Received: %s\nCount: %v\n", time, count)
			if count == 3 {
				break OuterLoop
			}
		case _ = <-timeout.C:
			t.Errorf("Timed out!\n")
			break OuterLoop
		}
	}
	c.Stop()
}
