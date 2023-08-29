package xdo

// #include <stdlib.h>
// #include <xdo.h>
// #cgo LDFLAGS: -lxdo
import "C"
import (
	"strings"
	"time"
	"unicode"
	"unsafe"
)

const CURRENTWINDOW = 0

type MouseButton int

const (
	MbLeft MouseButton = iota + 1
	MbMid
	MbRight
	MbWheelUp
	MbWheelDown
)

const (
	KeyReturn = "Return"
	KeySpace  = "space"
	KeyDelete = "Delete"
)

type Window int

type Xdo struct {
	xdo       *C.xdo_t
	ctrlDown  bool
	altDown   bool
	shiftDown bool

	Window   int
	KeyDelay int
}

func New() *Xdo {
	x := new(Xdo)
	x.xdo = C.xdo_new(nil)
	x.KeyDelay = 12000
	x.Window = CURRENTWINDOW
	return x
}

func (t *Xdo) MouseMove(x, y int) {
	C.xdo_move_mouse_relative(t.xdo, C.int(x), C.int(y))
}

func (t *Xdo) MouseDown(mouseButton MouseButton) {
	C.xdo_mouse_down(t.xdo, C.Window(t.Window), C.int(mouseButton))
}
func (t *Xdo) MouseUp(mouseButton MouseButton) {
	C.xdo_mouse_up(t.xdo, C.Window(t.Window), C.int(mouseButton))
}
func (t *Xdo) MousePress(mouseButton MouseButton, pressed bool) {
	if pressed {
		t.MouseDown(mouseButton)
	} else {
		t.MouseUp(mouseButton)
	}
}

func (t *Xdo) MouseClick(mouseButton MouseButton) {
	t.MouseDown(mouseButton)
	time.Sleep(10 * time.Millisecond)
	t.MouseUp(mouseButton)
}

func (t *Xdo) KeyPress(keyseq string) {
	if t.shiftDown && isLetter(keyseq) {
		keyseq = strings.ToUpper(keyseq)
	}
	if t.ctrlDown {
		keyseq = "Control_L+" + keyseq
	}
	if t.altDown {
		keyseq = "Alt_L+" + keyseq
	}
	str := C.CString(keyseq)
	defer C.free(unsafe.Pointer(str))
	C.xdo_send_keysequence_window(t.xdo, C.Window(t.Window), str, C.useconds_t(t.KeyDelay))
}

func (t *Xdo) EnterText(text string) {
	str := C.CString(text)
	defer C.free(unsafe.Pointer(str))
	C.xdo_enter_text_window(t.xdo, C.Window(t.Window), str, C.useconds_t(t.KeyDelay))
}

func (t *Xdo) SetCtrl(val bool) {
	t.ctrlDown = val
}
func (t *Xdo) SetShift(val bool) {
	t.shiftDown = val
}
func (t *Xdo) ToggleCapsLock() {
	t.shiftDown = !t.shiftDown
}
func (t *Xdo) ToggleCtrl() {
	t.ctrlDown = !t.ctrlDown
}

func (t *Xdo) ToggleAlt() {
	t.altDown = !t.altDown
}

func (t *Xdo) IsCapsLock() bool {
	return t.shiftDown
}

func isLetter(s string) bool {
	if len(s) != 1 {
		return false
	}
	c := s[0]
	return unicode.IsLetter(rune(c))
}
