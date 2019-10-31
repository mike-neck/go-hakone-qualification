package main

import (
	"fmt"
	"github.com/ledongthuc/pdf"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func makeText(x, y float64, value string) pdf.Text {
	return pdf.Text{X: x, Y: y, W: 0.0, S: value}
}

func makeBytes(x, y float64) pdf.Text {
	return pdf.Text{X: x, Y: y, W: 0.0, S: string([]byte{239, 191, 189})}
}

func TestDiscardingHeaderAnalyzer_Take(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "2"), // 0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "0"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, "1"),
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "9"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "備"),
		makeBytes(9.0, 2.0),
		makeText(11.0, 2.0, "　"), // 10
		makeBytes(11.0, 2.0),
		makeText(13.0, 2.0, "考"),
		makeBytes(13.0, 2.0),
		makeText(15.0, 2.0, "-"),
		makeBytes(15.0, 2.0),
		makeText(17.0, 2.0, "-"),
		makeBytes(17.0, 2.0),
		makeText(19.0, 2.0, "-"),
		makeBytes(19.0, 2.0),
		makeText(21.0, 2.0, "-"), // 20
		makeBytes(21.0, 2.0),
		makeText(23.0, 2.0, "-"),
		makeBytes(23.0, 2.0),
		makeText(25.0, 2.0, "-"),
		makeBytes(25.0, 2.0),
		makeText(27.0, 2.0, "1"), // 26
		makeBytes(27.0, 2.0),
		makeText(29.0, 2.0, "3"),
		makeBytes(29.0, 2.0),
		makeText(31.0, 2.0, "長"), // 30
		makeBytes(31.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := DiscardingHeaderAnalyzer(def)

	s, _, pos := analyzer.Take(0, texts)
	assert.True(t, s.empty)
	assert.Equal(t, Position(26), pos)
}

func TestStrIsNumber(t *testing.T) {
	assert.True(t, Str{value: "0"}.isNumber())
	assert.True(t, Str{value: "1"}.isNumber())
	assert.True(t, Str{value: "2"}.isNumber())
	assert.True(t, Str{value: "3"}.isNumber())
	assert.True(t, Str{value: "4"}.isNumber())
	assert.True(t, Str{value: "5"}.isNumber())
	assert.True(t, Str{value: "6"}.isNumber())
	assert.True(t, Str{value: "7"}.isNumber())
	assert.True(t, Str{value: "8"}.isNumber())
	assert.True(t, Str{value: "9"}.isNumber())
	assert.False(t, Str{value: "("}.isNumber())
}

func TestStrIsParenthesisOrAlNum(t *testing.T) {
	assert.True(t, Str{value: "0"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "1"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "2"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "3"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "4"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "5"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "6"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "7"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "8"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "9"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "("}.isParenthesisOrAlNum())
	assert.True(t, Str{value: "M"}.isParenthesisOrAlNum())
	assert.True(t, Str{value: ")"}.isParenthesisOrAlNum())
}

func TestRunnerNameAnalyzer_Take(t *testing.T) {
	texts := []pdf.Text{
		makeText(201.0, 20.0, "3"), // 0
		makeBytes(201.0, 20.0),
		makeText(203.0, 20.0, "0"),
		makeBytes(203.0, 20.0),
		makeText(205.0, 20.0, "1"),
		makeBytes(205.0, 20.0),
		makeText(207.0, 20.0, "4"),
		makeBytes(207.0, 20.0),
		makeText(209.0, 20.0, "2"),
		makeBytes(209.0, 20.0),
		makeText(211.0, 20.0, "石"), // 10
		makeBytes(211.0, 20.0),
		makeText(213.0, 20.0, "田"),
		makeBytes(213.0, 20.0),
		makeText(215.0, 20.0, "　"),
		makeBytes(215.0, 20.0),
		makeText(217.0, 20.0, "三"),
		makeBytes(217.0, 20.0),
		makeText(219.0, 20.0, "成"),
		makeBytes(219.0, 20.0),
		makeText(221.0, 20.0, "("), // 20
		makeBytes(221.0, 20.0),
		makeText(221.0, 20.0, "3"),
		makeBytes(221.0, 20.0),
		makeText(221.0, 20.0, ")"),
		makeBytes(221.0, 20.0),
	}

	var def DefaultAnalyzer
	analyzer := RunnerNameAnalyzer(def)

	result, _, pos := analyzer.Take(0, texts)

	es := makeStr(211.0, 20.0, "石田　三成")
	es.strType = RunnerNameValue
	assert.Equal(t, Position(20), pos)
	assert.Equal(t, es, result)
}

func TestGradeAnalyzer_Take(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "("), //0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "3"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, ")"),
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "東"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "京"),
		makeBytes(9.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := GradeAnalyzer(def)

	result, _, pos := analyzer.Take(0, texts)

	es := makeStr(1.0, 2.0, "(3)")
	es.strType = GradeValue

	assert.Equal(t, Position(6), pos)
	assert.Equal(t, es, result)
}

func TestGradeAnalyzer_Take2(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "("), //0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "3"),
		makeBytes(3.0, 2.0),
		makeText(7.0, 2.0, "東"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "京"),
		makeBytes(9.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := GradeAnalyzer(def)

	result, _, pos := analyzer.Take(0, texts)

	es := makeStr(1.0, 2.0, "(3)")
	es.strType = GradeValue

	assert.Equal(t, Position(4), pos)
	assert.Equal(t, es, result)
}

func TestStrEndsWith(t *testing.T) {
	s1 := makeStr(1.0, 2.0, "(3)")
	assert.True(t, s1.endsWith(")"))
}

func TestTeamAnalyzer_Take(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "東"),
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "京"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, "大"),
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "1"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "7"),
		makeBytes(9.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := TeamAnalyzer(def)

	result, _, pos := analyzer.Take(0, texts)

	es := makeStr(1.0, 2.0, "東京大")
	es.strType = TeamNameValue

	assert.Equal(t, Position(6), pos)
	assert.Equal(t, es, result)
}

func TestTimeAnalyzerTakeWithoutParenthesis_SuccessMinutes(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "大"), // 0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "学"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, "1"), // 4
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "1"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, ":"),
		makeBytes(9.0, 2.0),
		makeText(11.0, 2.0, "3"),
		makeBytes(11.0, 2.0),
		makeText(13.0, 2.0, "0"),
		makeBytes(13.0, 2.0),
		makeText(15.0, 2.0, "1"), // 14
		makeBytes(15.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := TimeAnalyzer{delegate: def, expectedYAxis: 2.0}

	result, pos := analyzer.takeWithoutParenthesis(4, texts)

	es := makeStr(5.0, 2.0, "11:30")

	assert.Equal(t, es, result)
	assert.Equal(t, Position(14), pos)
}

func TestTimeAnalyzerTakeWithoutParenthesis_SuccessHours(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "大"), // 0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "学"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, "1"), // 4
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, ":"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "1"),
		makeBytes(9.0, 2.0),
		makeText(11.0, 2.0, "3"),
		makeBytes(11.0, 2.0),
		makeText(13.0, 2.0, ":"),
		makeBytes(13.0, 2.0),
		makeText(15.0, 2.0, "3"),
		makeBytes(15.0, 2.0),
		makeText(17.0, 2.0, "0"),
		makeBytes(17.0, 2.0),
		makeText(19.0, 2.0, "1"), // 18
		makeBytes(19.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := TimeAnalyzer{delegate: def, expectedYAxis: 2.0}

	result, pos := analyzer.takeWithoutParenthesis(4, texts)

	es := makeStr(5.0, 2.0, "1:13:30")

	assert.Equal(t, es, result)
	assert.Equal(t, Position(18), pos)
}

func TestTimeAnalyzerTakeWithParenthesis_SuccessMinutes(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "大"), // 0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "学"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, "("), // 4
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "1"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "7"),
		makeBytes(9.0, 2.0),
		makeText(11.0, 2.0, ":"),
		makeBytes(11.0, 2.0),
		makeText(13.0, 2.0, "2"),
		makeBytes(13.0, 2.0),
		makeText(15.0, 2.0, "3"),
		makeBytes(15.0, 2.0),
		makeText(17.0, 2.0, ")"),
		makeBytes(17.0, 2.0),
		makeText(19.0, 2.0, "1"), // 18
		makeBytes(19.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := TimeAnalyzer{delegate: def, expectedYAxis: 2.0}

	result, pos := analyzer.takeWithParenthesis(4, texts)

	es := makeStr(7.0, 2.0, "17:23")

	assert.Equal(t, es, result)
	assert.Equal(t, Position(18), pos)
}

func TestTimeAnalyzerTakeMinutesTime(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "大"), // 0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "学"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, "1"), // 4
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "1"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, ":"),
		makeBytes(9.0, 2.0),
		makeText(11.0, 2.0, "3"),
		makeBytes(11.0, 2.0),
		makeText(13.0, 2.0, "0"),
		makeBytes(13.0, 2.0),
		makeText(15.0, 2.0, "1"), // 14
		makeBytes(15.0, 2.0),
	}

	var def DefaultAnalyzer
	analyzer := TimeAnalyzer{delegate: def, expectedYAxis: 2.0}

	result, pos := analyzer.takeMinutesTime(4, texts)

	es := makeStr(5.0, 2.0, "11:30")

	assert.Equal(t, Position(14), pos)
	assert.Equal(t, es, result)
}

func TestTimeAnalyzerTakeMinutesTime_Fail(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "大"), // 0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "学"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, "D"), // 4
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "N"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "S"),
		makeBytes(9.0, 2.0),
		makeText(11.0, 4.0, "A"),
		makeBytes(11.0, 4.0),
		makeText(13.0, 4.0, "O"),
		makeBytes(13.0, 4.0),
		makeText(15.0, 4.0, "I"), // 14
		makeBytes(15.0, 4.0),
	}

	var def DefaultAnalyzer
	analyzer := TimeAnalyzer{delegate: def, expectedYAxis: 2.0}

	result, pos := analyzer.takeMinutesTime(4, texts)

	assert.Equal(t, Position(16), pos)
	assert.Equal(t, emptyStr, result)
}

func TestNoteAnalyzer_Take(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.0, 2.0, "1"), // 0
		makeBytes(1.0, 2.0),
		makeText(3.0, 2.0, "2"),
		makeBytes(3.0, 2.0),
		makeText(5.0, 2.0, ":"),
		makeBytes(5.0, 2.0),
		makeText(7.0, 2.0, "3"),
		makeBytes(7.0, 2.0),
		makeText(9.0, 2.0, "4"),
		makeBytes(9.0, 2.0),
		makeText(11.0, 2.0, "D"), //10
		makeBytes(11.0, 2.0),
		makeText(13.0, 2.0, "Q"),
		makeBytes(13.0, 2.0),
		makeText(15.0, 2.0, "D"),
		makeBytes(15.0, 2.0),
		makeText(17.0, 2.0, "Q"),
		makeBytes(17.0, 2.0),
		makeText(19.0, 2.0, "2"),
		makeBytes(19.0, 2.0),
		makeText(1.0, 4.0, "Y"), // 20
		makeBytes(1.0, 4.0),
		makeText(3.0, 4.0, "A"),
		makeBytes(3.0, 4.0),
	}

	var def DefaultAnalyzer
	analyzer := NoteAnalyzer{delegate: def, expectedYAxis: 2.0}

	result, _, pos := analyzer.Take(Position(10), texts)

	es := makeStr(11.0, 2.0, "DQDQ2")
	es.strType = Notes

	assert.Equal(t, Position(20), pos)
	assert.Equal(t, es, result)
}

func TestData(t *testing.T) {
	data := [][]string{
		{"17:15",
			"34:45",
			"52:12",
			"DQ",
			"DQ2",
		},
		{"ISHIDA",
			"岐阜",
			"(17:30)",
		},
	}
	for i, ds := range data {
		yAxis := 21.0 + float64(i)*3.2
		for j, d := range ds {
			xAxisBase := 1.7 + float64(j)*9.7
			items := strings.Split(d, "")
			for k, item := range items {
				xAxis := xAxisBase + float64(k)*1.3
				fmt.Printf("makeText(%.1f, %.1f, \"%s\"),\n", xAxis, yAxis, item)
				fmt.Printf("makeBytes(%.1f, %.1f),\n", xAxis, yAxis)
			}
		}
	}
}

func TestTimeAnalyzers(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.7, 21.0, "1"),
		makeBytes(1.7, 21.0),
		makeText(3.0, 21.0, "7"),
		makeBytes(3.0, 21.0),
		makeText(4.3, 21.0, ":"),
		makeBytes(4.3, 21.0),
		makeText(5.6, 21.0, "1"),
		makeBytes(5.6, 21.0),
		makeText(6.9, 21.0, "5"),
		makeBytes(6.9, 21.0),
		makeText(11.4, 21.0, "3"),
		makeBytes(11.4, 21.0),
		makeText(12.7, 21.0, "4"),
		makeBytes(12.7, 21.0),
		makeText(14.0, 21.0, ":"),
		makeBytes(14.0, 21.0),
		makeText(15.3, 21.0, "4"),
		makeBytes(15.3, 21.0),
		makeText(16.6, 21.0, "5"),
		makeBytes(16.6, 21.0),
		makeText(21.1, 21.0, "5"),
		makeBytes(21.1, 21.0),
		makeText(22.4, 21.0, "2"),
		makeBytes(22.4, 21.0),
		makeText(23.7, 21.0, ":"),
		makeBytes(23.7, 21.0),
		makeText(25.0, 21.0, "1"),
		makeBytes(25.0, 21.0),
		makeText(26.3, 21.0, "2"),
		makeBytes(26.3, 21.0),
		makeText(30.8, 21.0, "1"),
		makeBytes(30.8, 21.0),
		makeText(32.1, 21.0, ":"),
		makeBytes(32.1, 21.0),
		makeText(33.4, 21.0, "1"),
		makeBytes(33.4, 21.0),
		makeText(34.7, 21.0, "1"),
		makeBytes(34.7, 21.0),
		makeText(36.0, 21.0, ":"),
		makeBytes(36.0, 21.0),
		makeText(37.3, 21.0, "3"),
		makeBytes(37.3, 21.0),
		makeText(38.6, 21.0, "4"),
		makeBytes(38.6, 21.0),
		makeText(40.5, 21.0, "1"),
		makeBytes(40.5, 21.0),
		makeText(41.8, 21.0, ":"),
		makeBytes(41.8, 21.0),
		makeText(43.1, 21.0, "1"),
		makeBytes(43.1, 21.0),
		makeText(44.4, 21.0, "5"),
		makeBytes(44.4, 21.0),
		makeText(45.7, 21.0, ":"),
		makeBytes(45.7, 21.0),
		makeText(47.0, 21.0, "2"),
		makeBytes(47.0, 21.0),
		makeText(48.3, 21.0, "3"),
		makeBytes(48.3, 21.0),
		makeText(1.7, 24.2, "I"),
		makeBytes(1.7, 24.2),
		makeText(3.0, 24.2, "S"),
		makeBytes(3.0, 24.2),
		makeText(4.3, 24.2, "H"),
		makeBytes(4.3, 24.2),
		makeText(5.6, 24.2, "I"),
		makeBytes(5.6, 24.2),
		makeText(6.9, 24.2, "D"),
		makeBytes(6.9, 24.2),
		makeText(8.2, 24.2, "A"),
		makeBytes(8.2, 24.2),
		makeText(11.4, 24.2, "岐"),
		makeBytes(11.4, 24.2),
		makeText(12.7, 24.2, "阜"),
		makeBytes(12.7, 24.2),
		makeText(21.1, 24.2, "("),
		makeBytes(21.1, 24.2),
		makeText(22.4, 24.2, "1"),
		makeBytes(22.4, 24.2),
		makeText(23.7, 24.2, "7"),
		makeBytes(23.7, 24.2),
		makeText(25.0, 24.2, ":"),
		makeBytes(25.0, 24.2),
		makeText(26.3, 24.2, "3"),
		makeBytes(26.3, 24.2),
		makeText(27.6, 24.2, "0"),
		makeBytes(27.6, 24.2),
		makeText(28.9, 24.2, ")"),
		makeBytes(28.9, 24.2),
		makeText(30.8, 24.2, "("),
		makeBytes(30.8, 24.2),
		makeText(32.1, 24.2, "1"),
		makeBytes(32.1, 24.2),
		makeText(33.4, 24.2, "7"),
		makeBytes(33.4, 24.2),
		makeText(34.7, 24.2, ":"),
		makeBytes(34.7, 24.2),
		makeText(36.0, 24.2, "2"),
		makeBytes(36.0, 24.2),
		makeText(37.3, 24.2, "7"),
		makeBytes(37.3, 24.2),
		makeText(38.6, 24.2, ")"),
		makeBytes(38.6, 24.2),
		makeText(40.5, 24.2, "("),
		makeBytes(40.5, 24.2),
		makeText(41.8, 24.2, "1"),
		makeBytes(41.8, 24.2),
		makeText(43.1, 24.2, "9"),
		makeBytes(43.1, 24.2),
		makeText(44.4, 24.2, ":"),
		makeBytes(44.4, 24.2),
		makeText(45.7, 24.2, "2"),
		makeBytes(45.7, 24.2),
		makeText(47.0, 24.2, "2"),
		makeBytes(47.0, 24.2),
		makeText(48.3, 24.2, ")"),
		makeBytes(48.3, 24.2),
	}
	makeTime := func(x, y float64, value string, strType StrType) Str {
		str := makeStr(x, y, value)
		str.strType = strType
		return str
	}

	var def DefaultAnalyzer
	ta := TimeAnalyzer{delegate: def, expectedYAxis: 21.0}
	a := To5kmAnalyzer(ta)
	var analyzer Analyzer
	analyzer = &a

	pos := Position(0)
	timeOf5km, analyzer, pos := analyzer.Take(pos, texts)
	timeOf10km, analyzer, pos := analyzer.Take(pos, texts)
	timeOf15km, analyzer, pos := analyzer.Take(pos, texts)
	timeOf20km, analyzer, pos := analyzer.Take(pos, texts)
	timeOfHalf, analyzer, pos := analyzer.Take(pos, texts)
	mayBeEmpty, analyzer, pos := analyzer.Take(pos, texts)
	rapFrom5kmTo10km, analyzer, pos := analyzer.Take(pos, texts)
	rapFrom10kmTo15km, analyzer, pos := analyzer.Take(pos, texts)
	rapFrom15kmTo20km, analyzer, pos := analyzer.Take(pos, texts)

	_, ok := analyzer.(*DoneAnalyzer)
	assert.True(t, ok, "テキスト終了")
	assert.Equal(t, Position(len(texts)), pos, "解析終了ポジション")
	assert.Equal(t, makeTime(1.7, 21.0, "17:15", Time5km), timeOf5km, "5kmタイム")
	assert.Equal(t, makeTime(11.4, 21.0, "34:45", Time10km), timeOf10km, "10kmタイム")
	assert.Equal(t, makeTime(21.1, 21.0, "52:12", Time15km), timeOf15km, "15kmタイム")
	assert.Equal(t, makeTime(30.8, 21.0, "1:11:34", Time20km), timeOf20km, "20kmタイム")
	assert.Equal(t, makeTime(40.5, 21.0, "1:15:23", ResultTime), timeOfHalf, "ハーフマラソンタイム")
	assert.Equal(t, emptyStr, mayBeEmpty)
	assert.Equal(t, makeTime(22.4, 24.2, "17:30", Rap5kmTo10km), rapFrom5kmTo10km, "ラップ5-10")
	assert.Equal(t, makeTime(32.1, 24.2, "17:27", Rap10kmTo15km), rapFrom10kmTo15km, "ラップ10-15")
	assert.Equal(t, makeTime(41.8, 24.2, "19:22", Rap15kmTo20km), rapFrom15kmTo20km, "ラップ15-20")
}

func TestTimeAnalyzers_FinishThenNextLine(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.7, 21.0, "1"),
		makeBytes(1.7, 21.0),
		makeText(3.0, 21.0, "7"),
		makeBytes(3.0, 21.0),
		makeText(4.3, 21.0, ":"),
		makeBytes(4.3, 21.0),
		makeText(5.6, 21.0, "1"),
		makeBytes(5.6, 21.0),
		makeText(6.9, 21.0, "5"),
		makeBytes(6.9, 21.0),
		makeText(11.4, 21.0, "3"),
		makeBytes(11.4, 21.0),
		makeText(12.7, 21.0, "4"),
		makeBytes(12.7, 21.0),
		makeText(14.0, 21.0, ":"),
		makeBytes(14.0, 21.0),
		makeText(15.3, 21.0, "4"),
		makeBytes(15.3, 21.0),
		makeText(16.6, 21.0, "5"),
		makeBytes(16.6, 21.0),
		makeText(21.1, 21.0, "5"),
		makeBytes(21.1, 21.0),
		makeText(22.4, 21.0, "2"),
		makeBytes(22.4, 21.0),
		makeText(23.7, 21.0, ":"),
		makeBytes(23.7, 21.0),
		makeText(25.0, 21.0, "1"),
		makeBytes(25.0, 21.0),
		makeText(26.3, 21.0, "2"),
		makeBytes(26.3, 21.0),
		makeText(30.8, 21.0, "D"),
		makeBytes(30.8, 21.0),
		makeText(32.1, 21.0, "Q"),
		makeBytes(32.1, 21.0),
		makeText(40.5, 21.0, "D"),
		makeBytes(40.5, 21.0),
		makeText(41.8, 21.0, "Q"),
		makeBytes(41.8, 21.0),
		makeText(43.1, 21.0, "2"),
		makeBytes(43.1, 21.0),
		makeText(1.7, 24.2, "I"),
		makeBytes(1.7, 24.2),
		makeText(3.0, 24.2, "S"),
		makeBytes(3.0, 24.2),
		makeText(4.3, 24.2, "H"),
		makeBytes(4.3, 24.2),
		makeText(5.6, 24.2, "I"),
		makeBytes(5.6, 24.2),
		makeText(6.9, 24.2, "D"),
		makeBytes(6.9, 24.2),
		makeText(8.2, 24.2, "A"),
		makeBytes(8.2, 24.2),
		makeText(11.4, 24.2, "岐"),
		makeBytes(11.4, 24.2),
		makeText(12.7, 24.2, "阜"),
		makeBytes(12.7, 24.2),
		makeText(21.1, 24.2, "("),
		makeBytes(21.1, 24.2),
		makeText(22.4, 24.2, "1"),
		makeBytes(22.4, 24.2),
		makeText(23.7, 24.2, "7"),
		makeBytes(23.7, 24.2),
		makeText(25.0, 24.2, ":"),
		makeBytes(25.0, 24.2),
		makeText(26.3, 24.2, "3"),
		makeBytes(26.3, 24.2),
		makeText(27.6, 24.2, "0"),
		makeBytes(27.6, 24.2),
		makeText(28.9, 24.2, ")"),
		makeBytes(28.9, 24.2),
		makeText(1.7, 27.4, "1"),
		makeBytes(1.7, 27.4),
		makeText(3.0, 27.4, "7"),
		makeBytes(3.0, 27.4),
	}

	var def DefaultAnalyzer
	ta := TimeAnalyzer{delegate: def, expectedYAxis: 21.0}
	a := To5kmAnalyzer(ta)
	var analyzer Analyzer
	analyzer = &a

	pos := Position(0)
	time5km, analyzer, pos := analyzer.Take(pos, texts)
	time10km, analyzer, pos := analyzer.Take(pos, texts)
	time15km, analyzer, pos := analyzer.Take(pos, texts)
	mayBeEmpty, analyzer, pos := analyzer.Take(pos, texts)
	mayNote, analyzer, pos := analyzer.Take(pos, texts)
	mayBeEmpty2, analyzer, pos := analyzer.Take(pos, texts)
	rap5kmTo10km, analyzer, pos := analyzer.Take(pos, texts)
	mayBeEmpty3, analyzer, pos := analyzer.Take(pos, texts) // 解析失敗 -> 次の行

	_, ok := analyzer.(*RunnerNameAnalyzer)
	assert.True(t, ok, "次の行", analyzer)
	assert.Equal(t, Position(len(texts)-4), pos, "解析終了ポジション")

	assert.Equal(t, Time5km, time5km.strType, "5kmタイム", time5km)
	assert.Equal(t, Time10km, time10km.strType, "10kmタイム", time10km)
	assert.Equal(t, Time15km, time15km.strType, "15kmタイム", time15km)
	assert.Equal(t, emptyStr, mayBeEmpty2, "20km取得失敗")
	assert.Equal(t, Notes, mayNote.strType, "ノートタイプ")
	assert.Equal(t, "DQDQ2", mayNote.value, "ノートDQDQ2")
	assert.Equal(t, emptyStr, mayBeEmpty, "読み飛ばし")
	assert.Equal(t, Rap5kmTo10km, rap5kmTo10km.strType, "ラップ5km-10km", rap5kmTo10km)
	assert.Equal(t, emptyStr, mayBeEmpty3, "解析不可")
}

func TestTimeAnalyzers_Finish(t *testing.T) {
	texts := []pdf.Text{
		makeText(1.7, 21.0, "1"),
		makeBytes(1.7, 21.0),
		makeText(3.0, 21.0, "7"),
		makeBytes(3.0, 21.0),
		makeText(4.3, 21.0, ":"),
		makeBytes(4.3, 21.0),
		makeText(5.6, 21.0, "1"),
		makeBytes(5.6, 21.0),
		makeText(6.9, 21.0, "5"),
		makeBytes(6.9, 21.0),
		makeText(11.4, 21.0, "3"),
		makeBytes(11.4, 21.0),
		makeText(12.7, 21.0, "4"),
		makeBytes(12.7, 21.0),
		makeText(14.0, 21.0, ":"),
		makeBytes(14.0, 21.0),
		makeText(15.3, 21.0, "4"),
		makeBytes(15.3, 21.0),
		makeText(16.6, 21.0, "5"),
		makeBytes(16.6, 21.0),
		makeText(21.1, 21.0, "5"),
		makeBytes(21.1, 21.0),
		makeText(22.4, 21.0, "2"),
		makeBytes(22.4, 21.0),
		makeText(23.7, 21.0, ":"),
		makeBytes(23.7, 21.0),
		makeText(25.0, 21.0, "1"),
		makeBytes(25.0, 21.0),
		makeText(26.3, 21.0, "2"),
		makeBytes(26.3, 21.0),
		makeText(30.8, 21.0, "D"),
		makeBytes(30.8, 21.0),
		makeText(32.1, 21.0, "Q"),
		makeBytes(32.1, 21.0),
		makeText(40.5, 21.0, "D"),
		makeBytes(40.5, 21.0),
		makeText(41.8, 21.0, "Q"),
		makeBytes(41.8, 21.0),
		makeText(43.1, 21.0, "2"),
		makeBytes(43.1, 21.0),
		makeText(1.7, 24.2, "I"),
		makeBytes(1.7, 24.2),
		makeText(3.0, 24.2, "S"),
		makeBytes(3.0, 24.2),
		makeText(4.3, 24.2, "H"),
		makeBytes(4.3, 24.2),
		makeText(5.6, 24.2, "I"),
		makeBytes(5.6, 24.2),
		makeText(6.9, 24.2, "D"),
		makeBytes(6.9, 24.2),
		makeText(8.2, 24.2, "A"),
		makeBytes(8.2, 24.2),
		makeText(11.4, 24.2, "岐"),
		makeBytes(11.4, 24.2),
		makeText(12.7, 24.2, "阜"),
		makeBytes(12.7, 24.2),
		makeText(21.1, 24.2, "("),
		makeBytes(21.1, 24.2),
		makeText(22.4, 24.2, "1"),
		makeBytes(22.4, 24.2),
		makeText(23.7, 24.2, "7"),
		makeBytes(23.7, 24.2),
		makeText(25.0, 24.2, ":"),
		makeBytes(25.0, 24.2),
		makeText(26.3, 24.2, "3"),
		makeBytes(26.3, 24.2),
		makeText(27.6, 24.2, "0"),
		makeBytes(27.6, 24.2),
		makeText(28.9, 24.2, ")"),
		makeBytes(28.9, 24.2),
	}

	var def DefaultAnalyzer
	ta := TimeAnalyzer{delegate: def, expectedYAxis: 21.0}
	a := To5kmAnalyzer(ta)
	var analyzer Analyzer
	analyzer = &a

	pos := Position(0)
	time5km, analyzer, pos := analyzer.Take(pos, texts)
	time10km, analyzer, pos := analyzer.Take(pos, texts)
	time15km, analyzer, pos := analyzer.Take(pos, texts)
	mayBeEmpty, analyzer, pos := analyzer.Take(pos, texts)
	mayNote, analyzer, pos := analyzer.Take(pos, texts)
	mayBeEmpty2, analyzer, pos := analyzer.Take(pos, texts)
	rap5kmTo10km, analyzer, pos := analyzer.Take(pos, texts)
	mayBeEmpty3, analyzer, pos := analyzer.Take(pos, texts) // 終了

	_, ok := analyzer.(*DoneAnalyzer)
	assert.True(t, ok, "終了", analyzer)
	assert.Equal(t, Position(len(texts)), pos, "解析終了ポジション")

	assert.Equal(t, Time5km, time5km.strType, "5kmタイム", time5km)
	assert.Equal(t, Time10km, time10km.strType, "10kmタイム", time10km)
	assert.Equal(t, Time15km, time15km.strType, "15kmタイム", time15km)
	assert.Equal(t, emptyStr, mayBeEmpty2, "20km取得失敗")
	assert.Equal(t, Notes, mayNote.strType, "ノートタイプ")
	assert.Equal(t, "DQDQ2", mayNote.value, "ノートDQDQ2")
	assert.Equal(t, emptyStr, mayBeEmpty, "読み飛ばし")
	assert.Equal(t, Rap5kmTo10km, rap5kmTo10km.strType, "ラップ5km-10km", rap5kmTo10km)
	assert.Equal(t, emptyStr, mayBeEmpty3, "終了")
}
