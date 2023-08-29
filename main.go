package main

import (
	"os"
	"time"

	"log/slog"

	"github.com/nvlled/carrot"
	"github.com/nvlled/dumbwheel/mouse"
	"github.com/nvlled/dumbwheel/xdo"
)

func main() {
	var programLevel = new(slog.LevelVar)
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	if os.Getenv("DEBUG") != "" {
		programLevel.Set(slog.LevelDebug)
	}

	var prevEvent *mouse.Event = nil
	var scrollScript *carrot.Script
	var loop *Interval

	move := 0
	maxMove := 1000
	xd := xdo.New()

	scrollScript = carrot.Start(func(ctrl *carrot.Control) {
		loopStart := time.Now().Add(-1 * time.Second)
		accelerated := false
		for {
			for prevEvent == nil {
				ctrl.Yield()
			}
			accelerated = time.Since(loopStart) < 300*time.Millisecond
			loopStart = time.Now()

			current := *prevEvent
			var mouseWheel xdo.MouseButton
			if current.Button == mouse.ButtonThumbUp {
				slog.Debug("start scroll up", "accelerated", accelerated)
				mouseWheel = xdo.MbWheelUp
			} else {
				slog.Debug("start scroll down", "accelerated", accelerated)
				mouseWheel = xdo.MbWheelDown
			}

			slog.Debug("scroll")
			xd.MouseClick(mouseWheel)

			sub := ctrl.StartAsync(func(ctrl *carrot.Control) {
				ctrl.Sleep(256 * time.Millisecond)
				for {
					if accelerated {
						xd.MouseClick(mouseWheel)
						xd.MouseClick(mouseWheel)
						if move > 200 {
							xd.MouseClick(mouseWheel)
						}
						if move > 500 {
							xd.MouseClick(mouseWheel)
							xd.MouseClick(mouseWheel)
						}
						if move > 1000 {
							xd.MouseClick(mouseWheel)
							xd.MouseClick(mouseWheel)
							xd.MouseClick(mouseWheel)
						}
						slog.Debug("scroll", "acceleration", move)
					} else {
						slog.Debug("scroll")
						xd.MouseClick(mouseWheel)
					}

					ctrl.Yield()
				}
			})

			loop.Start()
			for {
				if prevEvent.Type == mouse.EventOnUp && prevEvent.Button == current.Button {
					break
				}
				ctrl.Yield()
			}

			sub.Cancel()
			loop.Stop()
			prevEvent = nil
			slog.Debug("end scroll")
		}
	})

	loop = NewInterval(scrollScript.Update, 100*time.Millisecond)

	for event := range mouse.ReadEvents("/dev/input/event10") {
		if event.Type != mouse.EventOnMove || move == 0 {
			slog.Debug("read mouse event", "data", event)
		}

		if event.Type == mouse.EventOnMove {
			if move < maxMove {
				move++
			}
			continue
		}

		if event.Button == mouse.ButtonThumbUp || event.Button == mouse.ButtonThumbDown {
			prevEvent = &event
			move = 0
			scrollScript.Update()
		}
	}

}
