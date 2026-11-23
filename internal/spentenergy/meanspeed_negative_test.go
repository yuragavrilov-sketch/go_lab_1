package spentenergy

import (
	"testing"
	"time"
)

func TestMeanSpeedNegativeSteps(t *testing.T) {
	if got := MeanSpeed(-1000, 1.75, time.Hour); got != 0 {
		t.Fatalf("MeanSpeed() = %v, ожидался 0 при отрицательных шагах", got)
	}
}
