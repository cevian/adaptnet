package adaptnet

import (
	"fmt"
	"time"
)

type BufferFiller interface {
	FillBuffer() int //returns seconds of video filled in Ms
}

type SimulatedBuffer struct {
	targetLevelMs  int
	totalRunMs     int
	currentLevelMs int
	bufFill        BufferFiller
}

func NewSimulatedBuffer(targetLevelMs int, totalRunMs int, bufFill BufferFiller) *SimulatedBuffer {
	return &SimulatedBuffer{targetLevelMs, totalRunMs, 0, bufFill}
}

func (t *SimulatedBuffer) Run() error {
	runStart := time.Now()
	totalRunDuration := time.Duration(t.totalRunMs) * time.Millisecond

	for time.Since(runStart) < totalRunDuration {
		forStart := time.Now()
		if t.currentLevelMs > t.targetLevelMs {
			time.Sleep(time.Duration(t.currentLevelMs-t.targetLevelMs) * time.Millisecond)
		}

		t.currentLevelMs += t.bufFill.FillBuffer()

		drained := time.Since(forStart)
		t.currentLevelMs -= int(drained / time.Millisecond)
		fmt.Println("Buffer Level", t.currentLevelMs)
		if t.currentLevelMs < 0 {
			fmt.Println("Buffer Level Negative PAUSE")
			t.currentLevelMs = 0

		}
	}

	return nil
}

func (t *SimulatedBuffer) Stop() error {
	return nil
}
