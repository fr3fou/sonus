package main

import (
	"math"
	"time"

	"github.com/fr3fou/gusic/gusic"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxSamples          = 48000 * 5
	maxSamplesPerUpdate = 4096
)

type Key struct {
	Color rl.Color
	Note  gusic.Note
}

func NewKey(rect rl.Rectangle, note gusic.Note) *Key {
	return &Key{Note: note}
}

func (p *Key) Samples() []float32 {
	return samplesToFloat32(
		p.Note.Samples(
			// TODO, configurable params
			48000,
			math.Sin,
			gusic.NewLinearADSR(
				gusic.NewRatios(0.25, 0.25, 0.25, 0.25), 1.35, 0.35,
			),
		),
	)
}

func main() {
	rl.InitWindow(1650, 900, "goda - a simple music pad")

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice() // Close audio device (music streaming is automatically stopped)

	// Init raw audio stream (sample rate: 22050, sample size: 32bit-float, channels: 1-mono)
	stream := rl.InitAudioStream(48000, 32, 1)
	defer rl.CloseAudioStream(stream) // Close raw audio stream and delete buffers from RAM

	// Fill audio stream with some samples (sine wave)
	data := make([]float32, maxSamples)

	// NOTE: The generated MAX_SAMPLES do not fit to close a perfect loop
	// for that reason, there is a clip everytime audio stream is looped
	rl.PlayAudioStream(stream)

	totalSamples := int32(0)
	samplesLeft := int32(totalSamples)

	bpm := 200
	noteLength := 4

	breve := time.Minute / gusic.NoteDuration(bpm) * gusic.NoteDuration(noteLength) * 2
	semibreve := breve / 2
	// minim := semibreve / 2
	// crotchet := semibreve / 4
	quaver := semibreve / 8
	// semiquaver := semibreve / 16
	// demisemiquaver := semibreve / 32

	rl.SetTargetFPS(60)

	color := rl.Yellow
	for !rl.WindowShouldClose() {
		// Refill audio stream if required
		if rl.IsAudioStreamProcessed(stream) {
			numSamples := int32(0)
			if samplesLeft >= maxSamplesPerUpdate {
				numSamples = maxSamplesPerUpdate
			} else {
				numSamples = samplesLeft
			}

			rl.UpdateAudioStream(stream, data[totalSamples-samplesLeft:], numSamples)

			samplesLeft -= numSamples
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			coords := rl.GetMousePosition()
			if coords.X >= 100 && coords.Y >= 100 && coords.X <= 600+100 && coords.Y <= 600+100 {
				note := gusic.D(4, quaver, 0.125)
				samples := samplesToFloat32(
					note.Samples(
						// TODO, configurable params
						48000,
						math.Sin,
						gusic.NewLinearADSR(
							gusic.NewRatios(0.25, 0.25, 0.25, 0.25), 1.35, 0.35,
						),
					),
				)
				copy(data, samples)
				totalSamples = int32(len(samples))
				samplesLeft = totalSamples
			}
		}
		rl.DrawRectangle(100, 100, 600, 600, color)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func samplesToFloat32(in []float64) []float32 {
	samples := make([]float32, len(in))
	for i, v := range in {
		samples[i] = float32(v)
	}
	return samples
}
