package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	log "github.com/sirupsen/logrus"
)

var streamer beep.StreamSeekCloser
var format beep.Format
var buffer *beep.Buffer

// Initialize the crackling sound
func InitCrackling() {
	f, err := os.Open("crackle.wav")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err = wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/100))

	buffer = beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
}

// Finaly play the Crackle
func Crackle() {
	b := buffer.Streamer(0, buffer.Len())
	speaker.Play(b)
}
