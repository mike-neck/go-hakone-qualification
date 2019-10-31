package main

import (
	bytes2 "bytes"
	"github.com/ledongthuc/pdf"
)

type StrType int

const (
	String StrType = iota
	RunnerNameValue
	GradeValue
	TeamNameValue
	Time5km
	Time10km
	Time15km
	Time20km
	ResultTime
	Notes
	Rap5kmTo10km
	Rap10kmTo15km
	Rap15kmTo20km
)

type Str struct {
	empty   bool
	strType StrType
	xAxis   float64
	yAxis   float64
	value   string
}

var emptyStr = Str{empty: true}

func str(text pdf.Text) Str {
	return Str{empty: false, strType: String, xAxis: text.X, yAxis: text.Y, value: text.S}
}

func makeStr(x, y float64, value string) Str {
	return Str{empty: false, xAxis: x, yAxis: y, value: value}
}

type RunnerName Str

func runnerName(strs []Str) RunnerName {
	s := combineStr(strs)
	s.strType = RunnerNameValue
	return RunnerName(s)
}

type TeamName Str

func teamName(strs []Str) TeamName {
	s := combineStr(strs)
	s.strType = TeamNameValue
	return TeamName(s)
}

type GradeStr Str

func grade(strs []Str) GradeStr {
	count := 0
	for _, s := range strs {
		if s.value != "(" && s.value != ")" {
			count++
		}
	}
	res := make([]Str, count)
	index := 0
	for _, s := range strs {
		if s.value != "(" && s.value != ")" {
			res[index] = s
			index++
		}
	}
	s := combineStr(res)
	s.strType = GradeValue
	return GradeStr(s)
}

func combineStr(strs []Str) Str {
	buf := bytes2.NewBufferString("")
	for _, s := range strs {
		buf.WriteString(s.value)
	}
	return makeStr(strs[0].xAxis, strs[0].yAxis, buf.String())
}
