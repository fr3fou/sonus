package main

import (
	"math"
	"time"

	"github.com/fr3fou/gusic/gusic"
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	topMargin = 350
)

type Key struct {
	rl.Rectangle
	Texture        rl.Texture2D
	PressedTexture rl.Texture2D
	gusic.SingleNote
	IsSemitone bool
	IsActive   bool
}

func NewKey(note gusic.SingleNote, isSemitone bool, texture rl.Texture2D, pressedTexture rl.Texture2D) Key {
	return Key{SingleNote: note, IsSemitone: isSemitone, Texture: texture, PressedTexture: pressedTexture}
}

func (k *Key) Samples(generator gusic.Generator, adsr gusic.ADSR) []float32 {
	return samplesToFloat32(
		k.SingleNote.Samples(
			// TODO, configurable params
			48000,
			generator,
			adsr,
		),
	)
}

func (k *Key) Draw() {
	if !k.IsActive {
		rl.DrawTexturePro(k.Texture, rl.NewRectangle(0, 0, float32(k.Texture.Width), float32(k.Texture.Height)), k.Rectangle, rl.NewVector2(0, 0), 0, rl.White)
	} else {
		rl.DrawTexturePro(k.PressedTexture, rl.NewRectangle(0, 0, float32(k.Texture.Width), float32(k.Texture.Height)), k.Rectangle, rl.NewVector2(0, 0), 0, rl.White)
	}
}

func main() {
	width := int32(1680)
	height := int32(900)
	rl.InitWindow(width, height, "goda - a simple music synth")

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	stream := rl.InitAudioStream(48000, 32, 1)
	defer rl.CloseAudioStream(stream)

	maxSamples := 48000 * 5
	maxSamplesPerUpdate := int32(4096)

	data := make([]float32, maxSamples)

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

	volume := float32(0.125)

	_keys := []Key{}
	whiteKeys := []Key{}
	blackKeys := []Key{}
	startOctave := 3
	lastOctave := 5
	octaveCount := lastOctave - startOctave + 1 // +1 because it's inclusive

	whiteWidth := int(width / (int32(octaveCount) * 7)) // 7 is white keys per octave
	blackWidth := int(0.75 * float64(whiteWidth))

	generatorMap := map[string]gusic.Generator{
		"Sin":      math.Sin,
		"Sawtooth": gusic.Sawtooth(2 * math.Pi),
		"Square":   gusic.Square(2 * math.Pi),
		"Triangle": gusic.Triangle(2 * math.Pi),
	}
	generators := []string{"Sin", "Sawtooth", "Square", "Triangle"}
	generatorIndex := 0

	adsr := gusic.NewIdentityADSR()

	whiteTexture := rl.LoadTexture("white.png")
	blackTexture := rl.LoadTexture("black.png")
	whitePressedTexture := rl.LoadTexture("white_pressed.png")
	blackPressedTexture := rl.LoadTexture("black_pressed.png")
	sinTexture := rl.LoadTexture("sin.png")
	sawtoothTexture := rl.LoadTexture("sawtooth.png")
	squareTexture := rl.LoadTexture("square.png")
	triangleTexture := rl.LoadTexture("triangle.png")

	raygui.LoadGuiStyle("zahnrad.style")

	for i := startOctave; i <= lastOctave; i++ {
		// TODO: set duration to 0 and update it based on hold duration
		_keys = append(_keys,
			NewKey(gusic.C(i, quaver, 0), false, whiteTexture, whitePressedTexture),
			NewKey(gusic.CS(i, quaver, 0), true, blackTexture, blackPressedTexture),
			NewKey(gusic.D(i, quaver, 0), false, whiteTexture, whitePressedTexture),
			NewKey(gusic.DS(i, quaver, 0), true, blackTexture, blackPressedTexture),
			NewKey(gusic.E(i, quaver, 0), false, whiteTexture, whitePressedTexture),
			NewKey(gusic.F(i, quaver, 0), false, whiteTexture, whitePressedTexture),
			NewKey(gusic.FS(i, quaver, 0), true, blackTexture, blackPressedTexture),
			NewKey(gusic.G(i, quaver, 0), false, whiteTexture, whitePressedTexture),
			NewKey(gusic.GS(i, quaver, 0), true, blackTexture, blackPressedTexture),
			NewKey(gusic.A(i, quaver, 0), false, whiteTexture, whitePressedTexture),
			NewKey(gusic.AS(i, quaver, 0), true, blackTexture, blackPressedTexture),
			NewKey(gusic.B(i, quaver, 0), false, whiteTexture, whitePressedTexture),
		)
	}

	for _, key := range _keys {
		if !key.IsSemitone {
			whiteKeys = append(whiteKeys, key)
		} else {
			blackKeys = append(blackKeys, key)
		}
	}

	for i := range whiteKeys {
		rect := rl.NewRectangle(
			float32(i*whiteWidth),
			float32(topMargin),
			float32(whiteWidth),
			float32(height-int32(topMargin)),
		)
		whiteKeys[i].Rectangle = rect
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

		rect := rl.NewRectangle(
			float32(whiteWidth-blackWidth/2+i*whiteWidth+gapCount*whiteWidth),
			float32(topMargin),
			float32(blackWidth),
			float32(height-int32(topMargin))*0.6,
		)

		blackKeys[i].Rectangle = rect
		counter++
	}

	rl.SetTargetFPS(60)
	iconScale := float32(0.5)

	for !rl.WindowShouldClose() {
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

		pos := rl.GetMousePosition()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Handling presses
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			hasFound := false

			for i, key := range blackKeys {
				if rl.CheckCollisionPointRec(pos, key.Rectangle) {
					samples := key.Samples(generatorMap[generators[generatorIndex]], adsr)
					copy(data, samples)
					totalSamples = int32(len(samples))
					samplesLeft = totalSamples
					hasFound = true
					blackKeys[i].IsActive = true
					continue
				}
				blackKeys[i].IsActive = false
			}

			for i, key := range whiteKeys {
				if !hasFound && rl.CheckCollisionPointRec(pos, key.Rectangle) {
					samples := key.Samples(generatorMap[generators[generatorIndex]], adsr)
					copy(data, samples)
					totalSamples = int32(len(samples))
					samplesLeft = totalSamples
					whiteKeys[i].IsActive = true
					continue
				}
				whiteKeys[i].IsActive = false
			}
		} else {
			for i := range whiteKeys {
				whiteKeys[i].IsActive = false
				whiteKeys[i].SingleNote.Volume = float64(volume)
			}

			for i := range blackKeys {
				blackKeys[i].IsActive = false
				blackKeys[i].SingleNote.Volume = float64(volume)
			}
		}

		// Rendering white keys
		for i, key := range whiteKeys {
			key.Draw()
			rl.DrawRectangle(int32(i*whiteWidth), int32(topMargin), 1, height-int32(topMargin), rl.Gray)
		}

		// Rendering black keys
		for _, key := range blackKeys {
			key.Draw()
		}

		// Rendering settings
		generatorIndex = generatorInput(sinTexture, sawtoothTexture, squareTexture, triangleTexture, generatorIndex, generators, iconScale)
		adsr = gusic.NewIdentityADSR()
		volume = volumeInput(volume)

		// Rendering soundwave
		for i := 0; i < 4*100+4*3; i++ {
			rl.DrawPixelV(rl.NewVector2(float32(50+i), 50+50+3+float32(float32(sinTexture.Height)*iconScale)+50+100*data[i]), rl.Red)
		}

		// Rendering decorations
		rl.DrawLineEx(rl.NewVector2(0, float32(topMargin)), rl.NewVector2(float32(width), float32(topMargin)), 3, rl.Red)
		rl.DrawText("Goda", int32(width-rl.MeasureText("Goda", 50)-50), int32(50), 50, rl.White)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func generatorInput(sinTexture, sawtoothTexture, squareTexture, triangleTexture rl.Texture2D, generatorIndex int, generators []string, iconScale float32) int {
	rl.DrawTextureEx(sinTexture, rl.NewVector2(
		100*0+50-(iconScale*float32(sawtoothTexture.Width))/2+50,
		50+50+5,
	), 0, float32(iconScale), rl.Red)
	rl.DrawTextureEx(sawtoothTexture, rl.NewVector2(
		100*1+3+50-(iconScale*float32(sawtoothTexture.Width))/2+50,
		50+50+5,
	), 0, float32(iconScale), rl.Red)
	rl.DrawTextureEx(squareTexture, rl.NewVector2(
		100*2+3+50-(iconScale*float32(squareTexture.Width))/2+50,
		50+50+5,
	), 0, float32(iconScale), rl.Red)
	rl.DrawTextureEx(triangleTexture, rl.NewVector2(
		100*3+3+50-(iconScale*float32(triangleTexture.Width))/2+50,
		50+50+5,
	), 0, float32(iconScale), rl.Red)
	return raygui.ToggleGroup(rl.NewRectangle(50, 50, 100, 50), generators, generatorIndex)
}

func adsrInput(ratios gusic.ADSRRatios) gusic.ADSRRatios {
	return ratios
}

func volumeInput(volume float32) float32 {
	return raygui.SliderBar(rl.NewRectangle(50, topMargin-75, 4*100+4*3, 25), volume, 0, 0.3)
}

func samplesToFloat32(in []float64) []float32 {
	samples := make([]float32, len(in))
	for i, v := range in {
		samples[i] = float32(v)
	}
	return samples
}
