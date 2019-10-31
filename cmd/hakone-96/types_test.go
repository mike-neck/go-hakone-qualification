package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTime(t *testing.T) {
	time, err := NewTime("1:01:23")
	assert.Nil(t, err)
	assert.Equal(t, Time(61*60+23), time)
}

func TestNewTime_SubHour(t *testing.T) {
	time, err := NewTime("59:52")
	assert.Nil(t, err)
	assert.Equal(t, Time(59*60+52), time)
}

func TestNewTime_Failure(t *testing.T) {
	_, err := NewTime("DNS")
	assert.NotNil(t, err)
}

func TestNewTime_InvalidTimeNumber(t *testing.T) {
	_, err := NewTime("10:DD")
	assert.NotNil(t, err)
}
