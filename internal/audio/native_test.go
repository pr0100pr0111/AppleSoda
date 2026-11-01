//go:build darwin && cgo && !nocgo

package audio

import (
	"testing"
	"time"
)

func TestNativeCaptureLifecycle(t *testing.T) {
	c, err := NewNativeCapture()
	if err != nil {
		t.Skipf("native capture not available: %v", err)
	}

	if c.Source() != "native" {
		t.Fatalf("expected source native, got %s", c.Source())
	}

	if err := c.Start(); err != nil {
		t.Skipf("cannot start native capture (permissions or headless): %v", err)
	}

	if !c.IsRunning() {
		t.Fatal("expected running")
	}

	time.Sleep(100 * time.Millisecond)

	buf := make([]float32, 512)
	_ = c.ReadLatest(buf)
	_ = c.Level()

	if err := c.Stop(); err != nil {
		t.Fatalf("stop error: %v", err)
	}

	if c.IsRunning() {
		t.Fatal("expected not running after stop")
	}
}
