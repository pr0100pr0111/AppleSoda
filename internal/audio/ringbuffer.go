package audio

import (
	"math"
	"sync"
)

type RingBuffer struct {
	mu       sync.RWMutex
	data     []float32
	writePos int
	count    int
	capacity int
}

func NewRingBuffer(capacity int) *RingBuffer {
	if capacity <= 0 {
		capacity = 16384
	}
	return &RingBuffer{
		data:     make([]float32, capacity),
		capacity: capacity,
	}
}

func (rb *RingBuffer) Write(samples []float32) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for _, s := range samples {
		rb.data[rb.writePos] = s
		rb.writePos = (rb.writePos + 1) % rb.capacity
		if rb.count < rb.capacity {
			rb.count++
		}
	}
}

func (rb *RingBuffer) ReadLatest(dst []float32) int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	n := len(dst)
	if n > rb.count {
		n = rb.count
	}
	if n == 0 {
		return 0
	}

	start := (rb.writePos - n + rb.capacity) % rb.capacity
	for i := 0; i < n; i++ {
		idx := (start + i) % rb.capacity
		dst[i] = rb.data[idx]
	}
	return n
}

func (rb *RingBuffer) Level() float64 {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	samplesToCheck := 1024
	if samplesToCheck > rb.count {
		samplesToCheck = rb.count
	}
	if samplesToCheck == 0 {
		return 0
	}

	var sumSq float64
	start := (rb.writePos - samplesToCheck + rb.capacity) % rb.capacity
	for i := 0; i < samplesToCheck; i++ {
		idx := (start + i) % rb.capacity
		v := float64(rb.data[idx])
		sumSq += v * v
	}

	rms := math.Sqrt(sumSq / float64(samplesToCheck))
	return rms
}

func (rb *RingBuffer) Reset() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.writePos = 0
	rb.count = 0
}
