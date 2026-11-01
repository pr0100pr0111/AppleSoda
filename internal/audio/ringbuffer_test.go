package audio

import (
	"testing"
)

func TestRingBufferWriteAndReadLatest(t *testing.T) {
	rb := NewRingBuffer(8)

	src := []float32{1, 2, 3, 4, 5}
	rb.Write(src)

	dst := make([]float32, 3)
	n := rb.ReadLatest(dst)
	if n != 3 {
		t.Fatalf("expected 3 read, got %d", n)
	}
	if dst[0] != 3 || dst[1] != 4 || dst[2] != 5 {
		t.Fatalf("expected [3 4 5], got %v", dst)
	}

	more := []float32{6, 7, 8, 9, 10}
	rb.Write(more)

	dstAll := make([]float32, 8)
	n = rb.ReadLatest(dstAll)
	if n != 8 {
		t.Fatalf("expected 8 read, got %d", n)
	}
	expected := []float32{3, 4, 5, 6, 7, 8, 9, 10}
	for i, v := range expected {
		if dstAll[i] != v {
			t.Fatalf("at index %d expected %f, got %f", i, v, dstAll[i])
		}
	}
}

func TestRingBufferLevel(t *testing.T) {
	rb := NewRingBuffer(16)
	if rb.Level() != 0 {
		t.Fatal("expected level 0 on empty buffer")
	}

	rb.Write([]float32{1.0, 1.0, 1.0, 1.0})
	level := rb.Level()
	if level < 0.99 || level > 1.01 {
		t.Fatalf("expected level ~1.0, got %f", level)
	}
}
