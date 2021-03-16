package main

import (
	"math"
	"strconv"
	"time"

	"github.com/fr3fou/gusic/gusic"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Key struct {
	rl.Rectangle
	gusic.Note
	IsSemitone bool
}

func NewKey(note gusic.Note, isSemitone bool) Key {
	return Key{Note: note, IsSemitone: isSemitone}
}

func (k *Key) Samples() []float32 {
	return samplesToFloat32(
		k.Note.Samples(
			// TODO, configurable params
			48000,
			math.Sin,
			gusic.NewLinearADSR(
				gusic.NewRatios(0.25, 0.25, 0.25, 0.25), 1.35, 0.35,
			),
		),
	)
}

func (k *Key) Draw() {
	color := rl.White
	if k.IsSemitone {
		color = rl.Black
	}
	rl.DrawRectangleRec(k.Rectangle, color)
}

func main() {
	width := int32(1680)
	height := int32(900)
	rl.InitWindow(width, height, "goda - a simple music synth")

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	stream := rl.InitAudioStream(48000, 32, 1)
	defer rl.CloseAudioStream(stream)

	// maxSamples := 48000 * 5
	// maxSamplesPerUpdate := 4096

	// data := make([]float32, maxSamples)

	rl.PlayAudioStream(stream)

	// totalSamples := int32(0)
	// samplesLeft := int32(totalSamples)

	bpm := 200
	noteLength := 4

	breve := time.Minute / gusic.NoteDuration(bpm) * gusic.NoteDuration(noteLength) * 2
	semibreve := breve / 2
	// minim := semibreve / 2
	// crotchet := semibreve / 4
	quaver := semibreve / 8
	// semiquaver := semibreve / 16
	// demisemiquaver := semibreve / 32

	volume := 0.125

	keys := []Key{}
	whiteKeys := []Key{}
	blackKeys := []Key{}
	startOctave := 3
	lastOctave := 5
	octaveCount := lastOctave - startOctave + 1 // +1 because it's inclusive

	whiteWidth := int(width / (int32(octaveCount) * 7)) // 7 is white keys per octave
	blackWidth := int(0.75 * float64(whiteWidth))

	topMargin := 350

	for i := startOctave; i <= lastOctave; i++ {
		// TODO: set duration to 0 and update it based on hold duration
		keys = append(keys,
			NewKey(gusic.C(i, quaver, volume), false),
			NewKey(gusic.CS(i, quaver, volume), true),
			NewKey(gusic.D(i, quaver, volume), false),
			NewKey(gusic.DS(i, quaver, volume), true),
			NewKey(gusic.E(i, quaver, volume), false),
			NewKey(gusic.F(i, quaver, volume), false),
			NewKey(gusic.FS(i, quaver, volume), true),
			NewKey(gusic.G(i, quaver, volume), false),
			NewKey(gusic.GS(i, quaver, volume), true),
			NewKey(gusic.A(i, quaver, volume), false),
			NewKey(gusic.AS(i, quaver, volume), true),
			NewKey(gusic.B(i, quaver, volume), false),
		)
	}

	for _, key := range keys {
		if !key.IsSemitone {
			whiteKeys = append(whiteKeys, key)
		} else {
			blackKeys = append(blackKeys, key)
		}
	}

	for i := range whiteKeys {
		whiteKeys[i].Rectangle = rl.NewRectangle(
			float32(i*whiteWidth),
			float32(topMargin),
			float32(whiteWidth),
			float32(height-int32(topMargin)),
		)
	}
	counter := 0
	gapCount := 0
	for i := range blackKeys {
		if counter == 2 || counter == 5 {
			gapCount++
		}
		if counter == 5 {
			counter = 0
		}
		x := whiteWidth - blackWidth/2 + i*whiteWidth + gapCount*whiteWidth
		y := topMargin

		blackKeys[i].Rectangle = rl.NewRectangle(
			float32(x),
			float32(y),
			float32(blackWidth),
			float32(height-int32(topMargin))*0.6,
		)
		counter++
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Refill audio stream if required
		if rl.IsAudioStreamProcessed(stream) {
			// numSamples := int32(0)
			// if samplesLeft >= maxSamplesPerUpdate {
			// 	numSamples = maxSamplesPerUpdate
			// } else {
			// 	numSamples = samplesLeft
			// }

			// rl.UpdateAudioStream(stream, data[totalSamples-samplesLeft:], numSamples)

			// samplesLeft -= numSamples
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			// coords := rl.GetMousePosition()
		}

		for i, key := range whiteKeys {
			key.Draw()
			rl.DrawRectangle(int32(i*whiteWidth), int32(topMargin), 1, height-int32(topMargin), rl.Gray)
		}

		for _, key := range blackKeys {
			key.Draw()
		}

		rl.DrawLineEx(rl.NewVector2(0, float32(topMargin)), rl.NewVector2(float32(width), float32(topMargin)), 3, rl.Red)
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

const twelfthrootof2 float64 = 1.059463094359

var (
	a4 = 440
	// https://github.com/fr3fou/gusic/blob/72a7e32d5644ed6d123e365d416fdca51a268161/gusic/step.go#L38
	// c0    = float64(a4) * math.Pow(twelfthrootof2, float64(-4*12-9))
	c0    = float64(a4) * math.Pow(2, -4.75)
	notes = []string{
		"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B",
	}
)

func note(freq float64) string {
	h := int(math.Round(12 * math.Log2(freq/c0)))
	octave := h / 12

	return notes[h%12] + strconv.Itoa(octave)
}
