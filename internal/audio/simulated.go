package audio

import (
	"math"
	"sync"
	"time"
)

type SimulatedCapture struct {
	rb      *RingBuffer
	stopCh  chan struct{}
	running bool
	mu      sync.Mutex
	phase   float64
}

func NewSimulatedCapture() *SimulatedCapture {
	return &SimulatedCapture{
		rb: NewRingBuffer(16384),
	}
}

func (s *SimulatedCapture) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return ErrCaptureAlreadyRunning
	}

	s.running = true
	s.stopCh = make(chan struct{})

	go s.generateLoop()
	return nil
}

func (s *SimulatedCapture) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false
	close(s.stopCh)
	return nil
}

func (s *SimulatedCapture) generateLoop() {
	ticker := time.NewTicker(33 * time.Millisecond)
	defer ticker.Stop()

	chunkSize := 1024
	buf := make([]float32, chunkSize)

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			for i := 0; i < chunkSize; i++ {
				fi := float64(i)
				s1 := math.Sin(s.phase+fi*0.05) * 0.4
				s2 := math.Sin(s.phase*0.5+fi*0.12) * 0.3
				s3 := math.Sin(s.phase*2.0+fi*0.01) * 0.2
				buf[i] = float32(s1 + s2 + s3)
			}
			s.phase += 0.3
			s.mu.Unlock()

			s.rb.Write(buf)
		}
	}
}

func (s *SimulatedCapture) ReadLatest(dst []float32) int {
	return s.rb.ReadLatest(dst)
}

func (s *SimulatedCapture) Level() float64 {
	return s.rb.Level()
}

func (s *SimulatedCapture) Source() string {
	return "simulated"
}

func (s *SimulatedCapture) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
