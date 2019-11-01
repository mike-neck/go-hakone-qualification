package main

import (
	"encoding/json"
	"fmt"
	"github.com/ledongthuc/pdf"
	"github.com/mike-neck/go-hakone-qualification/hakone"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
)

func main() {
	file, reader, err := pdf.Open("data/hakone-96-personal.pdf")
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		log.Fatalln("error", "open file", "hakone-96-personal.pdf", err)
	}

	records := make([]hakone.Record, 0)

	maxPageNum := reader.NumPage()
	for pageNum := 1; pageNum <= maxPageNum; pageNum++ {
		page := reader.Page(pageNum)
		content := page.Content()
		texts := content.Text

		analyzer := NewAnalyzer()
		position := StartPosition()

		_, analyzer, position = analyzer.Take(position, texts) // discard header -> runner-name

		for i := 0; ; i++ {
			res, err := NextRecord(analyzer, position, texts)
			if err != nil {
				log.Fatalln("failed to load new record",
					"\nerror:", err,
					"\nindex:", i,
					"\nposition:", position,
					"\nrecords:", records)
			}
			if res.Record.Note == "" {
				records = append(records, res.Record)
			}
			position = res.Position
			if res.Done {
				break
			}
		}
	}

	jsonFile, err := os.Create("data/hakone-96-personal.jsonl")
	if err != nil {
		log.Fatalln("failed to open result file", err)
	}
	defer func() {
		_ = jsonFile.Close()
	}()
	encoder := json.NewEncoder(jsonFile)
	stdout := json.NewEncoder(os.Stdout)

	for idx, rec := range records {
		rec.Order = idx + 1
		_ = encoder.Encode(rec)
		_ = stdout.Encode(rec)
		fmt.Println("")
	}
}

type LoadResult struct {
	Record   hakone.Record
	Position Position
	Done     bool
}

func NextRecord(analyzer Analyzer, position Position, texts []pdf.Text) (*LoadResult, error) {
	if position.isOutOfRangeOf(texts) {
		result := LoadResult{Record: hakone.Record{}, Position: Position(len(texts)), Done: true}
		return &result, errors.New("already done")
	}
	an := analyzer
	pos := position
	runnerName, an, pos := an.Take(pos, texts)
	grade, an, pos := an.Take(pos, texts)
	team, an, pos := an.Take(pos, texts)

	var times Times
	t, an, pos := an.Take(pos, texts) // 5km
	if t != emptyStr {
		time, err := hakone.NewTime(t.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format time of TimeOf5km(%s) at %v", t.value, pos)
		}
		times.TimeOf5km = time
		t = emptyStr
	}
	if _, ok := an.(*NoteAnalyzer); ok == false {
		tim, a, p := an.Take(pos, texts) // 10km
		t = tim
		an = a
		pos = p
	}
	if t != emptyStr {
		time, err := hakone.NewTime(t.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format time of TimeOf10km(%s) at %v", t.value, pos)
		}
		times.TimeOf10km = time
		t = emptyStr
	}
	if _, ok := an.(*NoteAnalyzer); ok == false {
		tim, a, p := an.Take(pos, texts) // 15km
		t = tim
		an = a
		pos = p
	}
	if t != emptyStr {
		time, err := hakone.NewTime(t.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format time of TimeOf15km(%s) at %v", t.value, pos)
		}
		times.TimeOf15km = time
		t = emptyStr
	}
	if _, ok := an.(*NoteAnalyzer); ok == false {
		tim, a, p := an.Take(pos, texts) // 20km
		t = tim
		an = a
		pos = p
	}
	if t != emptyStr {
		time, err := hakone.NewTime(t.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format time of TimeOf20km(%s) at %v", t.value, pos)
		}
		times.TimeOf20km = time
		t = emptyStr
	}
	if _, ok := an.(*NoteAnalyzer); ok == false {
		tim, a, p := an.Take(pos, texts) // half
		t = tim
		an = a
		pos = p
	}
	if t != emptyStr {
		time, err := hakone.NewTime(t.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format time of TimeOfFinish(%s) at %v", t.value, pos)
		}
		times.TimeOfFinish = time
		t = emptyStr
	}

	var note Str
	if _, ok := an.(*NoteAnalyzer); ok { // note
		n, a, p := an.Take(pos, texts)
		note = n
		an = a
		pos = p
	}
	_, an, pos = an.Take(pos, texts) // discard

	var rap Raps
	rapTime := emptyStr
	if _, ok := an.(*RapTo10kmAnalyzer); ok {
		rt, a, p := an.Take(pos, texts)
		rapTime = rt
		an = a
		pos = p
	}
	if rapTime != emptyStr {
		rt, err := hakone.NewTime(rapTime.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format rap time of RapTo10km(%s) at %v", rapTime.value, pos)
		}
		rap.RapTo10km = rt
		rapTime = emptyStr
	}
	if _, ok := an.(*RapTo15kmAnalyzer); ok {
		rt, a, p := an.Take(pos, texts)
		rapTime = rt
		an = a
		pos = p
	}
	if rapTime != emptyStr {
		rt, err := hakone.NewTime(rapTime.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format rap time of RapTo15km(%s) at %v", rapTime.value, pos)
		}
		rap.RapTo15km = rt
		rapTime = emptyStr
	}
	if _, ok := an.(*RapTo20kmAnalyzer); ok {
		rt, a, p := an.Take(pos, texts)
		rapTime = rt
		an = a
		pos = p
	}
	if rapTime != emptyStr {
		rt, err := hakone.NewTime(rapTime.value)
		if err != nil {
			return &LoadResult{}, errors.Wrapf(err, "invalid format rap time of RapTo20km(%s) at %v", rapTime.value, pos)
		}
		rap.RapTo20km = rt
		rapTime = emptyStr
	}

	_, finished := an.(*DoneAnalyzer)
	_, succeeded := an.(*RunnerNameAnalyzer)
	if !finished && !succeeded {
		var d DefaultAnalyzer
		current, _, _ := d.Take(pos, texts)
		return &LoadResult{}, errors.New(
			fmt.Sprintf("invalid finish status at position: %v(%v), analyzer: %s(%v)", pos, current, AnalyzerName(an), an))
	}

	record := hakone.Record{
		Runner:            hakone.Runner(runnerName.value),
		Grade:             hakone.Grade(grade.value),
		Team:              hakone.TeamName(team.value),
		TimeOf5km:         times.TimeOf5km,
		TimeOf10km:        times.TimeOf10km,
		TimeOf15km:        times.TimeOf15km,
		TimeOf20km:        times.TimeOf20km,
		FinishTime:        times.TimeOfFinish,
		RapFrom5kmTo10km:  rap.RapTo10km,
		RapFrom10kmTo15km: rap.RapTo15km,
		RapFrom15kmTo20km: rap.RapTo20km,
		Note:              hakone.Note(note.value),
	}

	result := LoadResult{Record: record, Position: pos, Done: finished}

	return &result, nil
}

func AnalyzerName(a interface{}) string {
	if t := reflect.TypeOf(a); t.Kind() == reflect.Ptr {
		return fmt.Sprintf("*%s", t.Elem().Name())
	} else {
		return t.Name()
	}
}

type Times struct {
	TimeOf5km    hakone.Time
	TimeOf10km   hakone.Time
	TimeOf15km   hakone.Time
	TimeOf20km   hakone.Time
	TimeOfFinish hakone.Time
}

type Raps struct {
	RapTo10km hakone.Time
	RapTo15km hakone.Time
	RapTo20km hakone.Time
}
