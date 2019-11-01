package main

import (
	"bytes"
	"encoding/json"
	"github.com/ledongthuc/pdf"
	"github.com/mike-neck/go-hakone-qualification/hakone"
	"log"
	"os"
	"strings"
)

func main() {
	closeable, reader, err := pdf.Open("data/hakone-96-teams.pdf")
	if err != nil {
		log.Fatalln("failed to open data file", err)
	}
	defer func() {
		_ = closeable.Close()
	}()

	page := reader.Page(1)
	content := page.Content()
	texts := content.Text

	if len(texts) <= 0 {
		log.Fatalln("no data")
	}

	firstYAxis := texts[0].Y
	operator := NewOperator(firstYAxis)

	for _, text := range texts {
		operator.Operate(text)
	}
	names := operator.Names()
	teams := make([]hakone.Team, len(names))
	for index, name := range names {
		teams[index] = hakone.Team{Id: index + 1, Name: name}
	}

	jsonFile, err := os.Create("data/hakone-96-teams.jsonl")
	if err != nil {
		log.Fatalln("failed to create result json file", err)
	}
	defer func() {
		_ = jsonFile.Close()
	}()

	encoder := json.NewEncoder(jsonFile)
	stdout := json.NewEncoder(os.Stdout)
	for _, t := range teams {
		_ = encoder.Encode(t)
		_ = stdout.Encode(t)
	}
}

var invalidBytes = []byte{239, 191, 189}

func IsInvalidText(text pdf.Text) bool {
	return !IsValidText(text)
}

func IsValidText(text pdf.Text) bool {
	if text.S == "\n" {
		return false
	}
	bs := []byte(text.S)
	if len(bs) != 3 {
		return true
	}
	for i := 0; i < 3; i++ {
		if bs[i] != invalidBytes[i] {
			return true
		}
	}
	return false
}

type Operator interface {
	Operate(text pdf.Text)
	IsNewLine(text pdf.Text) bool
	Print()
	NewLineStart(text pdf.Text)
	Append(text pdf.Text)
	Names() []string
}

type BufferOperator struct {
	names  []string
	yAxis  float64
	buffer *bytes.Buffer
}

func (bo *BufferOperator) Operate(text pdf.Text) {
	if IsInvalidText(text) {
		return
	}
	if bo.IsNewLine(text) {
		bo.Print()
		bo.NewLineStart(text)
	} else {
		bo.Append(text)
	}
}

func (bo *BufferOperator) IsNewLine(text pdf.Text) bool {
	return bo.yAxis != text.Y
}

func (bo *BufferOperator) Print() {
	text := bo.buffer.String()
	if strings.Contains(text, "大学") &&
		!strings.Contains(text, "第96回") &&
		!strings.Contains(text, "人数") {
		bo.names = append(bo.names, text)
	}
}

func (bo *BufferOperator) NewLineStart(text pdf.Text) {
	var buffer bytes.Buffer
	buffer.WriteString(text.S)
	bo.yAxis = text.Y
	bo.buffer = &buffer
}

func (bo *BufferOperator) Append(text pdf.Text) {
	bo.buffer.WriteString(text.S)
}

func (bo *BufferOperator) Names() []string {
	return bo.names
}

func NewOperator(yAxis float64) Operator {
	var buffer bytes.Buffer
	operator := BufferOperator{yAxis: yAxis, buffer: &buffer, names: make([]string, 0)}
	return &operator
}
