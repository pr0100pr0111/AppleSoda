//go:build !darwin || !cgo || nocgo

package audio

func NewNativeCapture() (Capture, error) {
	return nil, ErrNativeUnsupported
}
