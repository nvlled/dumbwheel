package main

import (
	"time"

	"github.com/nvlled/quest"
)

type Interval struct {
	fn      func()
	running bool
	delay   time.Duration
	task    quest.VoidTask
}

func NewInterval(fn func(), delay time.Duration) *Interval {
	interval := &Interval{fn, false, delay, quest.NewVoidTask()}
	go interval.loop()
	return interval
}

func (interval *Interval) Start() {
	if !interval.running {
		interval.running = true
		interval.task.Resolve(quest.Void{})
	}
}

func (interval *Interval) Stop() {
	interval.task.Reset()
	interval.running = false
}

func (interval *Interval) IsRunning() bool {
	return interval.running
}

func (interval *Interval) loop() {
	for {
		interval.task.Await()
		for interval.running {
			interval.fn()
			time.Sleep(interval.delay)
		}
	}
}
