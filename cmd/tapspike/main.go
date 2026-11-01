package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pr0100pr0111/AppleSoda/internal/audio"
)

func main() {
	sourceName := "native"
	if len(os.Args) > 1 {
		sourceName = os.Args[1]
	}

	cap := audio.New(sourceName)
	fmt.Printf("AppleSoda audio capture spike\n")
	fmt.Printf("Requested source: %s | Active source: %s\n", sourceName, cap.Source())

	if err := cap.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start capture: %v\n", err)
		os.Exit(1)
	}
	defer cap.Stop()

	fmt.Printf("Capture running. Monitoring signal level (press Ctrl+C to stop)...\n\n")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	sampleBuf := make([]float32, 1024)
	startTime := time.Now()

	for {
		select {
		case <-sigCh:
			fmt.Printf("\nStopping capture...\n")
			return
		case <-ticker.C:
			elapsed := time.Since(startTime).Seconds()
			readCount := cap.ReadLatest(sampleBuf)
			lvl := cap.Level()

			meterWidth := 30
			scaled := int(math.Min(1.0, lvl*10.0) * float64(meterWidth))
			if scaled < 0 {
				scaled = 0
			}
			bar := strings.Repeat("█", scaled) + strings.Repeat("░", meterWidth-scaled)

			fmt.Printf("\rt=%.1fs | samples=%4d | level=%7.5f [%s]", elapsed, readCount, lvl, bar)
		}
	}
}
