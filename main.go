package main

import (
	"math"

	"github.com/fr3fou/beep/beep"
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	topMargin = 671
)

type Key struct {
	rl.Rectangle
	Texture        rl.Texture2D
	PressedTexture rl.Texture2D
	beep.SingleNote
	KeyboardKey int
	IsSemitone  bool
	IsActive    bool
}

func NewKey(note beep.SingleNote, isSemitone bool, texture rl.Texture2D, pressedTexture rl.Texture2D, keyboardKey int) Key {
	return Key{SingleNote: note, IsSemitone: isSemitone, Texture: texture, PressedTexture: pressedTexture, KeyboardKey: keyboardKey}
}

func (k *Key) Draw() {
	if !k.IsActive {
		rl.DrawTexturePro(k.Texture, rl.NewRectangle(0, 0, float32(k.Texture.Width), float32(k.Texture.Height)), k.Rectangle, rl.NewVector2(0, 0), 0, rl.White)
	} else {
		rl.DrawTexturePro(k.PressedTexture, rl.NewRectangle(0, 0, float32(k.Texture.Width), float32(k.Texture.Height)), k.Rectangle, rl.NewVector2(0, 0), 0, rl.White)
	}
}

const (
	sampleRate   = 48000
	streamBuffer = 4096
)

func main() {
	width := int32(896)
	height := int32(896)
	rl.InitWindow(width, height, "goda - a simple music synth")

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	stream := rl.InitAudioStream(sampleRate, 32, 1)
	defer rl.CloseAudioStream(stream)

	data := make([]float32, streamBuffer)

	rl.PlayAudioStream(stream)

	volume := float32(0.125)

	_keys := []Key{}
	whiteKeys := []Key{}
	blackKeys := []Key{}
	startOctave := 1
	lastOctave := 4
	octaveCount := lastOctave - startOctave + 1 // +1 because it's inclusive

	whiteWidth := int(width / (int32(octaveCount) * 7)) // 7 is white keys per octave
	blackWidth := int(0.75 * float64(whiteWidth))

	generatorMap := map[string]beep.Generator{
		"Sin":      math.Sin,
		"Sawtooth": beep.Sawtooth(2 * math.Pi),
		"Square":   beep.Square(1),
		"Triangle": beep.Triangle(2 * math.Pi),
	}
	generators := []string{"Sin", "Sawtooth", "Square", "Triangle"}
	generatorIndex := 0

	adsr := beep.NewEasedADSR(func(t float64) float64 { return t * t }, beep.NewRatios(0.25, 0.25, 0.25, 0.25), 1.25, 0.25)

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
		_keys = append(_keys,
			NewKey(beep.C(i, beep.NoDuration, 0), false, whiteTexture, whitePressedTexture, -1),
			NewKey(beep.CS(i, beep.NoDuration, 0), true, blackTexture, blackPressedTexture, -1),
			NewKey(beep.D(i, beep.NoDuration, 0), false, whiteTexture, whitePressedTexture, -1),
			NewKey(beep.DS(i, beep.NoDuration, 0), true, blackTexture, blackPressedTexture, -1),
			NewKey(beep.E(i, beep.NoDuration, 0), false, whiteTexture, whitePressedTexture, -1),
			NewKey(beep.F(i, beep.NoDuration, 0), false, whiteTexture, whitePressedTexture, -1),
			NewKey(beep.FS(i, beep.NoDuration, 0), true, blackTexture, blackPressedTexture, -1),
			NewKey(beep.G(i, beep.NoDuration, 0), false, whiteTexture, whitePressedTexture, -1),
			NewKey(beep.GS(i, beep.NoDuration, 0), true, blackTexture, blackPressedTexture, -1),
			NewKey(beep.A(i, beep.NoDuration, 0), false, whiteTexture, whitePressedTexture, -1),
			NewKey(beep.AS(i, beep.NoDuration, 0), true, blackTexture, blackPressedTexture, -1),
			NewKey(beep.B(i, beep.NoDuration, 0), false, whiteTexture, whitePressedTexture, -1),
		)
	}
	keyboardKeys := []int{
		// ---
		rl.KeyA,
		rl.KeyW,
		rl.KeyS,
		rl.KeyE,
		rl.KeyD,
		rl.KeyF,
		rl.KeyT,
		rl.KeyG,
		rl.KeyY,
		rl.KeyH,
		rl.KeyU,
		rl.KeyJ,
		// ---
		rl.KeyK,
		rl.KeyO,
		rl.KeyL,
		rl.KeyP,
		rl.KeySemicolon,
		rl.KeyApostrophe,
		rl.KeyRightBracket,
		rl.KeyBackSlash,
		-1,
		-1,
		-1,
		-1,
	}

	for i := 0; i < 12*2; i++ { // loop for 2 octaves
		_keys[i+12].KeyboardKey = keyboardKeys[i]
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
			zeroSamples(data)
			for _, k := range blackKeys {
				if k.IsActive {
					samples := make([]float64, streamBuffer)
					for i := 0; i < streamBuffer; i++ {
						phi := (k.Frequency / sampleRate)
						t := 2.0 * math.Pi * float64(i) * (phi - math.Floor(phi))
						samples[i] += k.SampleAt(t, sampleRate, generatorMap[generators[generatorIndex]])
					}
					beep.ApplyADSR(samples, adsr)
					for i := range samples {
						data[i] += float32(samples[i])
					}
				}
			}
			for _, k := range whiteKeys {
				if k.IsActive {
					samples := make([]float64, streamBuffer)
					for i := 0; i < streamBuffer; i++ {
						phi := (k.Frequency / sampleRate)
						t := 2.0 * math.Pi * float64(i) * (phi - math.Floor(phi))
						samples[i] += k.SampleAt(t, sampleRate, generatorMap[generators[generatorIndex]])
					}
					beep.ApplyADSR(samples, adsr)
					for i := range samples {
						data[i] += float32(samples[i])
					}
				}
			}
			rl.UpdateAudioStream(stream, data, streamBuffer)
		}

		pos := rl.GetMousePosition()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Handling presses
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			hasFound := false

			for i, key := range blackKeys {
				if rl.CheckCollisionPointRec(pos, key.Rectangle) {
					hasFound = true
					blackKeys[i].IsActive = true
					continue
				}
				blackKeys[i].IsActive = false
			}

			for i, key := range whiteKeys {
				if !hasFound && rl.CheckCollisionPointRec(pos, key.Rectangle) {
					whiteKeys[i].IsActive = true
					continue
				}
				whiteKeys[i].IsActive = false
			}
		} else {
			for i := range blackKeys {
				blackKeys[i].IsActive = false
				blackKeys[i].SingleNote.Volume = float64(volume)
			}

			for i := range whiteKeys {
				whiteKeys[i].IsActive = false
				whiteKeys[i].SingleNote.Volume = float64(volume)
			}
		}

		for i, key := range whiteKeys {
			if rl.IsKeyDown(int32(key.KeyboardKey)) {
				whiteKeys[i].IsActive = true
				continue
			}
			// whiteKeys[i].IsActive = false
		}

		for i, key := range blackKeys {
			if rl.IsKeyDown(int32(key.KeyboardKey)) {
				blackKeys[i].IsActive = true
				continue
			}
			// blackKeys[i].IsActive = false
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
		// adsr = beep.NewIdentityADSR()
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

func adsrInput(ratios beep.ADSRRatios) beep.ADSRRatios {
	return ratios
}

func volumeInput(volume float32) float32 {
	return raygui.SliderBar(rl.NewRectangle(50, topMargin-75, 4*100+4*3, 25), volume, 0, 0.3)
}

func zeroSamples(data []float32) {
	for i := range data {
		data[i] = 0.0
	}
}

func float64toFloat32(in []float64) []float32 {
	samples := make([]float32, len(in))
	for i, v := range in {
		samples[i] = float32(v)
	}
	return samples
}

func float32toFloat64(in []float32) []float64 {
	samples := make([]float64, len(in))
	for i, v := range in {
		samples[i] = float64(v)
	}
	return samples
}
