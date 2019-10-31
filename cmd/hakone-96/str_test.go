package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrade_WithParentheses(t *testing.T) {
	strs := []Str{
		makeStr(1.0, 2.0, "("),
		makeStr(3.0, 2.0, "1"),
		makeStr(5.0, 2.0, ")"),
	}
	g := grade(strs)
	es := makeStr(3.0, 2.0, "1")
	es.strType = GradeValue
	assert.Equal(t, GradeStr(es), g)
}

func TestCombineStr(t *testing.T) {
	strs := []Str{
		makeStr(1.0, 2.0, "("),
		makeStr(3.0, 2.0, "1"),
		makeStr(5.0, 2.0, ")"),
	}
	res := combineStr(strs)
	assert.Equal(t, makeStr(1.0, 2.0, "(1)"), res)
}

func TestTeamName(t *testing.T) {
	strs := []Str{
		makeStr(1.0, 2.0, "東"),
		makeStr(3.0, 2.0, "海"),
		makeStr(5.0, 2.0, "大"),
	}
	team := teamName(strs)
	es := makeStr(1.0, 2.0, "東海大")
	es.strType = TeamNameValue
	assert.Equal(t, TeamName(es), team)
}

func TestRunnerName(t *testing.T) {
	strs := []Str{
		makeStr(1.0, 2.0, "山"),
		makeStr(3.0, 2.0, "田"),
		makeStr(5.0, 2.0, "太"),
		makeStr(7.0, 2.0, "郎"),
	}

	runner := runnerName(strs)
	es := makeStr(1.0, 2.0, "山田太郎")
	es.strType = RunnerNameValue
	assert.Equal(t, RunnerName(es), runner)
}
