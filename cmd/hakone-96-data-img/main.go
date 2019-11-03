package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mike-neck/go-hakone-qualification/hakone"
	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"os"
)

func main() {
	file, err := os.Open("data/hakone-96-personal.jsonl")
	if err != nil {
		log.Fatalln("failed to open file: data/hakone-96-personal.jsonl by ", err)
	}
	defer func() {
		_ = file.Close()
	}()

	plotImg, err := plot.New()
	if err != nil {
		log.Fatalln("failed to prepare plot by", err)
	}
	plotImg.Title.Text = "Qualification Data"
	//plotImg.Title.Text = "ssssssssss"
	plotImg.X.Label.Text = "Persons"
	plotImg.Y.Label.Text = "Total Time"
	plotImg.Y.Tick.Marker = Tick{}

	grid := plotter.NewGrid()
	grid.Horizontal.Color = color.RGBA{R: 21, G: 21, B: 43, A: 0}
	plotImg.Add(grid)

	teamPlots := map[hakone.TeamName]*TeamPlot{
		"東京国際大学": NewTeamPlot("Tokyo Kokusai Univ", 0, 12, 192),
		"山梨学院大学": NewTeamPlot("Yamanashi Gakuin Univ", 21, 21, 127),
		//"筑波大学": NewTeamPlot("Tsukuba Univ", 13, 169, 169),
		"麗澤大学":  NewTeamPlot("Reitaku Univ", 192, 34, 0),
		"中央大学":  NewTeamPlot("Chuo Univ", 62, 62, 0),
		"上武大学":  NewTeamPlot("Joubu Univ", 168, 0, 194),
		"早稲田大学": NewTeamPlot("Waseda Univ", 62, 52, 10),
		"駿河台大学": NewTeamPlot("Surugadai Univ", 10, 14, 86),
	}

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		bytes := scanner.Bytes()
		var record hakone.Record
		err := json.Unmarshal(bytes, &record)
		if err != nil {
			fmt.Println("error at line:", i+1, "error: ", err, "json: ", string(bytes))
			continue
		}

		if p, ok := teamPlots[record.Team]; ok && p.Index < 10 {
			p.Append(record)
		}
	}

	params := make([]interface{}, len(teamPlots)*2)
	index := 0
	for n, p := range teamPlots {
		xy, err := p.ToPlot()
		if err != nil {
			log.Fatalln("failed to create plot data: ", p, "at:", n, "cause:", err)
		}

		params[index*2] = p.Name
		params[index*2+1] = xy
		index++
	}
	err = plotutil.AddLinePoints(plotImg, params...)
	if err != nil {
		log.Fatalln("failed to add points to plot, cause:", err)
	}

	if err = plotImg.Save(1440, 810, "build/hakone-96-img.png"); err != nil {
		log.Fatalln("failed to save file", err)
	}
}

type Tick struct{}

func (Tick) Ticks(min, max float64) []plot.Tick {
	return plot.DefaultTicks{}.Ticks(min, max)
}

type SinglePlot struct {
	Index int
	Sum   int
}

func (sp *SinglePlot) ToPlot() plotter.XY {
	return plotter.XY{
		X: float64(sp.Index),
		Y: float64(sp.Sum),
	}
}

type TeamPlot struct {
	Name  string
	Plots []SinglePlot
	Index int
	Sum   int
	Color color.Color
}

func NewTeamPlot(name string, red, green, blue uint8) *TeamPlot {
	return &TeamPlot{
		Name:  name,
		Index: 0,
		Sum:   0,
		Plots: make([]SinglePlot, 10),
		Color: color.RGBA{R: red, G: green, B: blue, A: 255},
	}
}

var even3Minutes30Seconds = (3*60+30)*21 + 21

func (tp *TeamPlot) Append(record hakone.Record) {
	tp.Sum += even3Minutes30Seconds - int(record.FinishTime)
	tp.Plots[tp.Index] = SinglePlot{
		Index: tp.Index + 1,
		Sum:   tp.Sum,
	}
	tp.Index += 1
}

func (tp *TeamPlot) ToPlot() (plotter.XYer, error) {
	xys := make(plotter.XYs, len(tp.Plots))
	for index, p := range tp.Plots {
		xys[index] = p.ToPlot()
	}
	line, err := plotter.NewLine(xys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create line")
	}
	line.LineStyle.Color = tp.Color
	line.LineStyle.Width = vg.Points(2)
	return line, nil
}
