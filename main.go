package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxSamples          = 22050
	maxSamplesPerUpdate = 4096
)

func main() {
	rl.InitWindow(800, 800, "gad - a simple music pad")

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

	totalSamples := int32(maxSamples)
	samplesLeft := int32(totalSamples)

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

			// Reset samples feeding (loop audio)
			if samplesLeft <= 0 {
				samplesLeft = totalSamples
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			coords := rl.GetMousePosition()
			if coords.X >= 100 && coords.Y >= 100 && coords.X <= 600+100 && coords.Y <= 600+100 {
				fmt.Println("clicked")
			}
		}
		rl.DrawRectangle(100, 100, 600, 600, color)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
