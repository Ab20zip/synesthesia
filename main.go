package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

const (
	width  = 800
	height = 600
)

func main() {
	f, err := os.Open("path_to_audio_file.mp3")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	scatter, err := plotter.NewScatter(nil)
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scatter)

	win, err := p.WriterTo(width, height, "png")
	if err != nil {
		log.Fatal(err)
	}

	for {
		if _, ok := streamer.(*beep.Buffer); !ok {
			break
		}

		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		streamer.Close()

		feature := calculateFeature(buffer)

		scatter.XYs = append(scatter.XYs, struct{ X, Y float64 }{float64(len(scatter.XYs)), feature})

		if err := p.Save(width, height, "synesthesia_plot.png"); err != nil {
			log.Fatal(err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func calculateFeature(buffer *beep.Buffer) float64 {
	return float64(len(buffer.Len()))
}
