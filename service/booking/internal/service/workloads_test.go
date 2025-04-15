package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinTime(t *testing.T) {
	times := []time.Time{
		time.Now().Add(time.Hour),
		time.Now(),
		time.Now().Add(time.Hour * 2),
	}

	res := minTime(times...)
	assert.Equal(t, times[1], res)
}

func TestMaxTime(t *testing.T) {
	times := []time.Time{
		time.Now().Add(time.Hour),
		time.Now(),
		time.Now().Add(time.Hour * 2),
	}

	res := maxTime(times...)
	assert.Equal(t, times[2], res)
}
