package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"log/slog"

	"github.com/nvlled/carrot"
	"github.com/nvlled/dumbwheel/mouse"
	"github.com/nvlled/dumbwheel/xdo"
)

func findMouseEventDevice() string {
	dir := "/dev/input/by-id"
	entries, err := os.ReadDir(dir)
	ruhOh(err)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), "-event-mouse") {
			return dir + "/" + entry.Name()
		}
	}

	return ""
}

func usage() {
	program := os.Args[0]
	fmt.Printf("usage: %v [input event filename]\n", program)
	fmt.Printf("if no filename is given, it will automatically look for one.\n")
	fmt.Printf("\n")
	fmt.Printf("Examples: %v /dev/input/10\n", program)
	fmt.Printf("          %v /dev/input/by-id/usb-MOSART_Semi._2.4G_INPUT_DEVICE-if01-event-mouse\n", program)
	fmt.Printf("          %v /dev/input/by-path/pci-0000:00:14.0-usb-0:2:1.1-event-mouse\n", program)
	fmt.Printf("\n")
	fmt.Printf("Note: when reading from /dev/input/by-id or /dev/input/by-path, look for filenames\n")
	fmt.Printf("      that ends with '-event-mouse'\n")
	os.Exit(1)
}

func main() {
	deviceName := ""
	if len(os.Args) < 2 {
		deviceName = findMouseEventDevice()
	} else {
		deviceName = os.Args[1]
	}

	if deviceName == "" {
		usage()
	}

	slog.Info("using input event device")
	slog.Info(deviceName)

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

	for event := range mouse.ReadEvents(deviceName) {
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

func ruhOh(err error) {
	if err != nil {
		panic(err)
	}
}
