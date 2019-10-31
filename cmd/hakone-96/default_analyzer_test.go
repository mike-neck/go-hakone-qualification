package main

import (
	"github.com/ledongthuc/pdf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultSourceOp_Take_Duplicate_Case(t *testing.T) {
	t1 := pdf.Text{X: 10.0, Y: 30.0, W: 0, S: "test"}
	t2 := pdf.Text{X: 10.0, Y: 30.0, W: 0, S: "failure"}
	texts := []pdf.Text{t1, t2}

	var op DefaultAnalyzer

	str, _, pos := op.Take(0, texts)

	assert.Equal(t, Position(2), pos)
	assert.Equal(t, 10.0, str.xAxis)
	assert.Equal(t, 30.0, str.yAxis)
	assert.Equal(t, "test", str.value)
}

func TestDefaultSourceOp_Take_Single_Case(t *testing.T) {
	t1 := pdf.Text{X: 10.0, Y: 30.0, W: 0, S: "test"}
	t2 := pdf.Text{X: 12.0, Y: 30.0, W: 0, S: "next"}
	texts := []pdf.Text{t1, t2}

	var op DefaultAnalyzer

	str, _, pos := op.Take(0, texts)

	assert.Equal(t, Position(1), pos)
	assert.Equal(t, 10.0, str.xAxis)
	assert.Equal(t, 30.0, str.yAxis)
	assert.Equal(t, "test", str.value)
}
