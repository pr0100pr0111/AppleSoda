package audio

import (
	"errors"
)

var (
	ErrCaptureAlreadyRunning = errors.New("audio capture is already running")
	ErrCaptureNotRunning     = errors.New("audio capture is not running")
	ErrNativeUnsupported     = errors.New("native CoreAudio process tap is not supported on this platform or build")
)

type Capture interface {
	Start() error
	Stop() error
	ReadLatest(dst []float32) int
	Level() float64
	Source() string
	IsRunning() bool
}

func New(preferredSource string) Capture {
	switch preferredSource {
	case "native":
		c, err := NewNativeCapture()
		if err == nil {
			return c
		}
		return NewSimulatedCapture()
	case "simulated":
		return NewSimulatedCapture()
	case "cava":
		return NewSimulatedCapture()
	default:
		c, err := NewNativeCapture()
		if err == nil {
			return c
		}
		return NewSimulatedCapture()
	}
}
