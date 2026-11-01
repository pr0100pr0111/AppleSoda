//go:build darwin && cgo && !nocgo

package audio

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework CoreAudio
#include "native_darwin.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

type NativeCapture struct {
	mu sync.Mutex
}

func NewNativeCapture() (Capture, error) {
	return &NativeCapture{}, nil
}

func (n *NativeCapture) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	code := int(C.applesoda_tap_start())
	if code != 0 {
		return fmt.Errorf("coreaudio tap failed with code %d", code)
	}
	return nil
}

func (n *NativeCapture) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	C.applesoda_tap_stop()
	return nil
}

func (n *NativeCapture) ReadLatest(dst []float32) int {
	if len(dst) == 0 {
		return 0
	}
	ptr := (*C.float)(unsafe.Pointer(&dst[0]))
	return int(C.applesoda_tap_read_samples(ptr, C.int(len(dst))))
}

func (n *NativeCapture) Level() float64 {
	return float64(C.applesoda_tap_get_level())
}

func (n *NativeCapture) Source() string {
	return "native"
}

func (n *NativeCapture) IsRunning() bool {
	return C.applesoda_tap_is_running() != 0
}
