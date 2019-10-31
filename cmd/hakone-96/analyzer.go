package main

import (
	"github.com/ledongthuc/pdf"
	"regexp"
	"strings"
)

func StartPosition() Position {
	return 0
}

type Position int

func (p Position) isOutOfRangeOf(a []pdf.Text) bool {
	i := int(p)
	return len(a) <= i
}

func (p Position) isAtTheEndOf(a []pdf.Text) bool {
	i := int(p)
	return i == len(a)-1
}

type Analyzer interface {
	Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position)
}

func NewAnalyzer() Analyzer {
	var def DefaultAnalyzer
	analyzer := DiscardingHeaderAnalyzer(def)
	return &analyzer
}

type DefaultAnalyzer struct {
}

func (a *DefaultAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	max := len(texts)
	if pos.isOutOfRangeOf(texts) {
		return emptyStr, a, Position(max)
	}
	current := texts[pos]
	if pos.isAtTheEndOf(texts) {
		return str(current), a, Position(max)
	}
	next := texts[pos+1]
	if current.X == next.X && current.Y == next.Y {
		return str(current), a, pos + 2
	} else {
		return str(current), a, pos + 1
	}
}

type Collector interface {
	Add(str Str) Collector
}

type EmptyCollector struct {
}

func (e *EmptyCollector) Add(str Str) Collector {
	return e
}

func (d *DefaultAnalyzer) Seek(start Position, texts []pdf.Text, continueIfTrue func(Str) bool) Position {
	var c EmptyCollector
	pos, _ := d.seekAndCollect(start, texts, &c, continueIfTrue)
	return pos
}

type ArrayCollector struct {
	items []Str
}

func (a *ArrayCollector) Add(str Str) Collector {
	items := append(a.items, str)
	return &ArrayCollector{items: items}
}

func (d *DefaultAnalyzer) SeekAndCollect(start Position, texts []pdf.Text, continueIfTrue func(Str) bool) (Position, []Str) {
	result := make([]Str, 0)
	collector := &ArrayCollector{items: result}
	pos, c := d.seekAndCollect(start, texts, collector, continueIfTrue)
	if ac, ok := c.(*ArrayCollector); ok {
		return pos, ac.items
	}
	return pos, result
}

func (d *DefaultAnalyzer) seekAndCollect(start Position, texts []pdf.Text, collector Collector, continueIfTrue func(Str) bool) (Position, Collector) {
	result := collector
	pos := start
	max := len(texts)
	if pos.isOutOfRangeOf(texts) {
		return Position(max), result
	}
	finishIfTrue := func(s Str) bool {
		return !continueIfTrue(s)
	}
	for {
		current, _, next := d.Take(pos, texts)
		if finishIfTrue(current) {
			return pos, result
		}
		pos = next
		result = result.Add(current)
		if pos.isOutOfRangeOf(texts) {
			return Position(max), result
		}
	}
}

type DiscardingHeaderAnalyzer DefaultAnalyzer

func (d *DiscardingHeaderAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	delegate := DefaultAnalyzer(*d)
	max := len(texts)
	p := pos
	if pos.isOutOfRangeOf(texts) {
		return emptyStr, d, Position(max)
	}
	p = delegate.Seek(p, texts, Str.isNotHyphen)
	p = d.SeekToHeaderFinish(p, texts)
	p = delegate.Seek(p, texts, Str.isNotNumber)

	analyzer := RunnerNameAnalyzer(delegate)
	return emptyStr, &analyzer, p
}

func (d *DiscardingHeaderAnalyzer) SeekToHeaderFinish(start Position, texts []pdf.Text) Position {
	analyzer := DefaultAnalyzer(*d)
	pos := start
	if pos.isOutOfRangeOf(texts) {
		return Position(len(texts))
	}
	for startHyphen := false; startHyphen == false; {
		current, _, p := analyzer.Take(pos, texts)
		if current.isHyphen() {
			next, _, _ := analyzer.Take(p, texts)
			if next.isHyphen() {
				startHyphen = true
			}
		}
		pos = p
		if pos.isOutOfRangeOf(texts) {
			return Position(len(texts))
		}
	}
	for endHyphen := false; endHyphen == false; {
		current, _, p := analyzer.Take(pos, texts)
		if current.isNotHyphen() {
			endHyphen = true
		}
		pos = p
		if pos.isOutOfRangeOf(texts) {
			return Position(len(texts))
		}
	}
	return pos
}

func (s Str) isHyphen() bool {
	return s.value == "-"
}

func (s Str) isNotHyphen() bool {
	return !s.isHyphen()
}

type RunnerNameAnalyzer DefaultAnalyzer

func (ra *RunnerNameAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	delegate := DefaultAnalyzer(*ra)
	p := delegate.Seek(pos, texts, Str.isNumber)
	p, strs := delegate.SeekAndCollect(p, texts, Str.isNotParenthesisNeitherAlNum)
	cs := combineStr(strs)
	cs.strType = RunnerNameValue
	analyzer := GradeAnalyzer(delegate)
	return cs, &analyzer, p
}

var numberPattern = regexp.MustCompile("^[0-9]$")

func (s Str) isNumber() bool {
	return numberPattern.MatchString(s.value)
}

func (s Str) isNotNumber() bool {
	return !s.isNumber()
}

var parenthesisOrAlNum = regexp.MustCompile("^[0-9()a-zA-Z]$")

func (s Str) isParenthesisOrAlNum() bool {
	return parenthesisOrAlNum.MatchString(s.value)
}

func (s Str) isNotParenthesisNeitherAlNum() bool {
	return !s.isParenthesisOrAlNum()
}

type GradeAnalyzer DefaultAnalyzer

func (ga *GradeAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	delegate := DefaultAnalyzer(*ga)
	p, strs := delegate.SeekAndCollect(pos, texts, Str.isParenthesisOrAlNum)
	cs := combineStr(strs)
	if !cs.endsWith(")") {
		cs = combineStr([]Str{cs, makeStr(cs.xAxis, cs.yAxis, ")")})
	}
	cs.strType = GradeValue
	analyzer := TeamAnalyzer(delegate)
	return cs, &analyzer, p
}

func (s Str) endsWith(prefix string) bool {
	return strings.HasSuffix(s.value, prefix)
}

type TeamAnalyzer DefaultAnalyzer

func (ta *TeamAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	delegate := DefaultAnalyzer(*ta)
	p, strs := delegate.SeekAndCollect(pos, texts, Str.isNotNumber)
	name := combineStr(strs)
	name.strType = TeamNameValue

	analyzer := TimeAnalyzer{delegate: delegate, expectedYAxis: name.yAxis}
	to5kmAnalyzer := To5kmAnalyzer(analyzer)
	return name, &to5kmAnalyzer, p
}

type TimeAnalyzer struct {
	delegate      DefaultAnalyzer
	expectedYAxis float64
}

func (t *TimeAnalyzer) takeWithoutParenthesis(pos Position, texts []pdf.Text) (Str, Position) {
	max := len(texts)
	if pos.isOutOfRangeOf(texts) {
		return emptyStr, Position(max)
	}
	analyzer := t.delegate
	p := pos
	first, _, nextPos := analyzer.Take(p, texts)
	if first.isNotNumber() {
		return emptyStr, Position(max)
	}
	p = nextPos
	second, _, nextPos := analyzer.Take(p, texts)
	if second.isNumber() {
		return t.takeMinutesTime(pos, texts)
	} else if second.isColon() {
		return t.takeHourTime(pos, texts)
	}
	return emptyStr, Position(max)
}

func (t *TimeAnalyzer) takeWithParenthesis(pos Position, texts []pdf.Text) (Str, Position) {
	max := len(texts)
	if pos.isOutOfRangeOf(texts) {
		return emptyStr, Position(max)
	}
	analyzer := t.delegate
	p := pos
	first, _, nextPos := analyzer.Take(p, texts)
	if first.empty {
		return emptyStr, Position(len(texts))
	} else if first.isNotParenthesis() && first.yAxis != t.expectedYAxis {
		return emptyStr, p
	} else if first.isNotParenthesis() {
		return emptyStr, Position(len(texts))
	}
	p = nextPos
	result, nextPos := t.takeWithoutParenthesis(p, texts)
	p = nextPos
	end, _, nextPos := analyzer.Take(p, texts)
	if end.isNotParenthesis() {
		return emptyStr, Position(max)
	}
	p = nextPos
	return result, p
}

func (s Str) isColon() bool {
	return s.value == ":"
}

func (s Str) isTimeChar(isColon bool) bool {
	if isColon {
		return s.value == ":"
	}
	return s.isNumber()
}

var parenthesisChars = regexp.MustCompile("^[()]$")

func (s Str) isParenthesis() bool {
	return parenthesisChars.MatchString(s.value)
}

func (s Str) isNotParenthesis() bool {
	return !s.isParenthesis()
}

func (s Str) isNotTimeChar(isColon bool) bool {
	return !s.isTimeChar(isColon)
}

func (t *TimeAnalyzer) takeMinutesTime(pos Position, texts []pdf.Text) (Str, Position) {
	return t.takeTime(pos, 5, texts, func(i int) bool {
		return i == 2
	})
}

func (t *TimeAnalyzer) takeHourTime(pos Position, texts []pdf.Text) (Str, Position) {
	return t.takeTime(pos, 7, texts, func(i int) bool {
		return i == 1 || i == 4
	})
}

func (t *TimeAnalyzer) takeTime(pos Position, size int, texts []pdf.Text, colonPosition func(int) bool) (Str, Position) {
	strs := make([]Str, size)
	p := pos
	for i := 0; i < size; i++ {
		current, _, next := t.delegate.Take(p, texts)
		if current.isNotTimeChar(colonPosition(i)) || current.yAxis != t.expectedYAxis {
			return emptyStr, Position(len(texts))
		}
		strs[i] = current
		p = next
	}
	return combineStr(strs), p
}

func delegationAnalyzeTime(
	strType StrType,
	ta *TimeAnalyzer,
	pos Position,
	texts []pdf.Text,
	toNextAnalyzer func(Position) Analyzer) (Str, Analyzer, Position) {
	timeStr, next := ta.takeWithoutParenthesis(pos, texts)
	if next.isOutOfRangeOf(texts) {
		noteAnalyzer := NoteAnalyzer{ta.delegate, ta.expectedYAxis}
		return emptyStr, &noteAnalyzer, pos
	}
	timeStr.strType = strType
	nt := toNextAnalyzer(next)
	return timeStr, nt, next
}

type To5kmAnalyzer TimeAnalyzer

func (t *To5kmAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*t)
	return delegationAnalyzeTime(Time5km, &analyzer, pos, texts, func(np Position) Analyzer {
		n := To10kmAnalyzer(analyzer)
		return &n
	})
}

type To10kmAnalyzer TimeAnalyzer

func (t *To10kmAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*t)
	return delegationAnalyzeTime(Time10km, &analyzer, pos, texts, func(np Position) Analyzer {
		n := To15kmAnalyzer(analyzer)
		return &n
	})
}

type To15kmAnalyzer TimeAnalyzer

func (t *To15kmAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*t)
	return delegationAnalyzeTime(Time15km, &analyzer, pos, texts, func(np Position) Analyzer {
		n := To20kmAnalyzer(analyzer)
		return &n
	})
}

type To20kmAnalyzer TimeAnalyzer

func (t *To20kmAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*t)
	return delegationAnalyzeTime(Time20km, &analyzer, pos, texts, func(np Position) Analyzer {
		n := FinishTimeAnalyzer(analyzer)
		return &n
	})
}

type FinishTimeAnalyzer TimeAnalyzer

func (t *FinishTimeAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*t)
	return delegationAnalyzeTime(ResultTime, &analyzer, pos, texts, func(np Position) Analyzer {
		def := analyzer.delegate
		ns, _, _ := def.Take(np, texts)
		next := DiscardingRomeAndNativeAnalyzer{def, ns.yAxis}
		return &next
	})
}

type NoteAnalyzer struct {
	delegate      DefaultAnalyzer
	expectedYAxis float64
}

func (n *NoteAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	max := len(texts)
	if pos.isOutOfRangeOf(texts) {
		return emptyStr, n, Position(max)
	}
	p, strs := n.delegate.SeekAndCollect(pos, texts, func(s Str) bool {
		return s.isNoteChar(n.expectedYAxis)
	})
	cs := combineStr(strs)
	cs.strType = Notes

	nextStr, _, _ := n.delegate.Take(p, texts)

	da := DiscardingRomeAndNativeAnalyzer{n.delegate, nextStr.yAxis}

	return cs, &da, p
}

var noteChars = regexp.MustCompile("^[A-Z0-9]$")

func (s Str) isNoteChar(sameLineYAxis float64) bool {
	return noteChars.MatchString(s.value) && s.yAxis == sameLineYAxis
}

func (s Str) isNotNoteChar(sameLineYAxis float64) bool {
	return s.isNoteChar(sameLineYAxis)
}

type DiscardingRomeAndNativeAnalyzer struct {
	delegate      DefaultAnalyzer
	expectedYAxis float64
}

func (d *DiscardingRomeAndNativeAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	delegate := d.delegate
	p := pos
	p = delegate.Seek(p, texts, func(s Str) bool {
		return s.isNotParenthesis() && s.yAxis == d.expectedYAxis
	})
	analyzer := TimeAnalyzer{delegate: delegate, expectedYAxis: d.expectedYAxis}
	next := RapTo10kmAnalyzer(analyzer)
	return emptyStr, &next, p
}

type RapTo10kmAnalyzer TimeAnalyzer

func (ra *RapTo10kmAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*ra)
	time, next := analyzer.takeWithParenthesis(pos, texts)

	if time.empty && next.isOutOfRangeOf(texts) {
		return emptyStr, analyzeDone, Position(len(texts))
	} else if time.empty {
		def := analyzer.delegate
		na := RunnerNameAnalyzer(def)
		return emptyStr, &na, pos
	}
	time.strType = Rap5kmTo10km

	if next.isOutOfRangeOf(texts) {
		return time, analyzeDone, Position(len(texts))
	}

	na := RapTo15kmAnalyzer(analyzer)
	return time, &na, next
}

type RapTo15kmAnalyzer TimeAnalyzer

func (ra *RapTo15kmAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*ra)
	time, next := analyzer.takeWithParenthesis(pos, texts)

	if time.empty && next.isOutOfRangeOf(texts) {
		return emptyStr, analyzeDone, Position(len(texts))
	} else if time.empty {
		def := analyzer.delegate
		na := RunnerNameAnalyzer(def)
		return emptyStr, &na, pos
	}
	time.strType = Rap10kmTo15km

	if next.isOutOfRangeOf(texts) {
		return time, analyzeDone, Position(len(texts))
	}

	na := RapTo20kmAnalyzer(analyzer)
	return time, &na, next
}

type RapTo20kmAnalyzer TimeAnalyzer

func (ra *RapTo20kmAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	analyzer := TimeAnalyzer(*ra)
	time, next := analyzer.takeWithParenthesis(pos, texts)
	if time.empty && next.isOutOfRangeOf(texts) {
		return emptyStr, analyzeDone, Position(len(texts))
	} else if time.empty {
		def := analyzer.delegate
		na := RunnerNameAnalyzer(def)
		return emptyStr, &na, pos
	}

	time.strType = Rap15kmTo20km

	if next.isOutOfRangeOf(texts) {
		return time, analyzeDone, Position(len(texts))
	}

	def := analyzer.delegate
	na := RunnerNameAnalyzer(def)

	return time, &na, next
}

type DoneAnalyzer struct {
}

var analyzeDone = &DoneAnalyzer{}

func (d *DoneAnalyzer) Take(pos Position, texts []pdf.Text) (Str, Analyzer, Position) {
	return emptyStr, d, Position(len(texts))
}
