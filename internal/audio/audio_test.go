package audio

import (
	"testing"
	"time"
)

func TestSimulatedCapture(t *testing.T) {
	c := NewSimulatedCapture()
	if c.Source() != "simulated" {
		t.Fatalf("expected source simulated, got %s", c.Source())
	}
	if c.IsRunning() {
		t.Fatal("expected not running initially")
	}

	if err := c.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	if !c.IsRunning() {
		t.Fatal("expected running after start")
	}

	time.Sleep(100 * time.Millisecond)

	buf := make([]float32, 256)
	n := c.ReadLatest(buf)
	if n == 0 {
		t.Fatal("expected non-zero samples read from simulated capture")
	}

	lvl := c.Level()
	if lvl <= 0 {
		t.Fatalf("expected positive level, got %f", lvl)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("stop error: %v", err)
	}
	if c.IsRunning() {
		t.Fatal("expected not running after stop")
	}
}

func TestNewFactory(t *testing.T) {
	c := New("simulated")
	if c.Source() != "simulated" {
		t.Fatalf("expected simulated source, got %s", c.Source())
	}

	autoC := New("auto")
	if autoC == nil {
		t.Fatal("expected non-nil capture for auto")
	}
}
