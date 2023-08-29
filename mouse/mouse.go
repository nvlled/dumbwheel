package mouse

// #include <xdo.h>
// #include <linux/input.h>
// #cgo LDFLAGS: -lxdo
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func ruhOh(err error) {
	if err != nil {
		panic(err)
	}
}

type EventType int

const (
	EventNone EventType = iota
	EventOnMove
	EventOnDown
	EventOnUp
)

func (etype EventType) String() string {
	switch etype {
	case EventNone:
		return "-"
	case EventOnMove:
		return "mousemove"
	case EventOnDown:
		return "mousedown"
	case EventOnUp:
		return "mouseup"
	}
	return ""
}

type ButtonType int

const (
	ButtonNone ButtonType = iota
	ButtonLeft
	ButtonRight
	ButtonMiddle
	ButtonThumbDown
	ButtonThumbUp
)

func (etype ButtonType) String() string {
	switch etype {
	case ButtonNone:
		return "-"
	case ButtonLeft:
		return "left"
	case ButtonRight:
		return "right"
	case ButtonMiddle:
		return "middle"
	case ButtonThumbDown:
		return "thumbdown"
	case ButtonThumbUp:
		return "thumbup"
	}
	return ""
}

type Event struct {
	Type   EventType
	Button ButtonType
	RelX   int8
	RelY   int8
}

func (event Event) String() string {
	if event.Type == EventOnDown || event.Type == EventOnUp {
		return fmt.Sprintf("MouseEvent{type=%v, button=%v}", event.Type.String(), event.Button.String())
	} else if event.Type == EventOnMove {
		return fmt.Sprintf("MouseEvent{type=%v, rel_x=%v, rel_y=%v}", event.Type.String(), event.RelX, event.RelY)
	}
	return fmt.Sprintf("MouseEvent{type=%v}", event.Type.String())
}

type TimeVal struct {
	Seconds int64
	Millis  int64
}

type InputEvent struct {
	Time  TimeVal
	Type  uint16
	Code  uint16
	Value int32
}

var event C.struct_input_event

func ReadEvents(filename string) <-chan Event {
	ch := make(chan Event, 0)

	go func() {
		file, err := os.Open(filename)
		ruhOh(err)
		defer file.Close()
		fd := C.int(file.Fd())
		eventPtr := unsafe.Pointer(&event)
		for {
			var bytes C.ssize_t
			bytes = C.read(fd, eventPtr, C.sizeof_struct_input_event)
			if bytes < C.sizeof_struct_input_event {
				break
			}
			if event._type != 1 && event._type != 2 {
				continue
			}

			mouseEvent := Event{}

			event := InputEvent{
				Type:  uint16(event._type),
				Code:  uint16(event.code),
				Value: int32(event.value),
			}

			if event.Type != 1 && event.Type != 2 {
				continue
			}

			if event.Type == 1 {
				if event.Value == 1 {
					mouseEvent.Type = EventOnDown
				} else {
					mouseEvent.Type = EventOnUp

				}

				switch event.Code {
				case 272:
					mouseEvent.Button = ButtonLeft
				case 273:
					mouseEvent.Button = ButtonRight
				case 274:
					mouseEvent.Button = ButtonMiddle
				case 275:
					mouseEvent.Button = ButtonThumbDown
				case 276:
					mouseEvent.Button = ButtonThumbUp
				}
			} else if event.Type == 2 {
				mouseEvent.Button = ButtonNone
				mouseEvent.Type = EventOnMove
				if event.Code == 0 {
					mouseEvent.RelX = int8(event.Value)
				} else {
					mouseEvent.RelY = int8(event.Value)
				}
			}

			ch <- mouseEvent
		}

	}()

	return ch
}
