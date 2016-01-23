package main

import (
	"fmt"
	"io"

	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
)

// Player handles audio player runtime
type Player struct {
	Lock    chan string
	Queue   chan string
	Error   chan error
	Running bool
}

// Add will add video to the queue
func (p *Player) Add(video *Video) {
	p.Queue <- video.ID
}

// Start will start playing the songs from the queue
func (p *Player) Start() {
	go func() {
		for p.Running {
			id := <-p.Queue
			p.Lock <- id
			file := fmt.Sprintf("%s/.jable/%s.mp3", userDir, id)
			p.Play(file)
			<-p.Lock
		}
	}()
}

// Play starts playback from specified file
func (p *Player) Play(file string) {
	handle, err := mpg123.Open(file)
	if err != nil {
		handleErr(err)
		return
	}
	defer handle.Close()

	ao.Initialize()
	defer ao.Shutdown()
	dev := ao.NewLiveDevice(aoSampleFormat(handle))
	defer dev.Close()

	if _, err := io.Copy(dev, handle); err != nil {
		handleErr(err)
		return
	}
}

// Stop will stop the player
func (p *Player) Stop() {
	p.Running = false
	mpg123.Exit()
}

// NewPlayer returns a new Player instance
func NewPlayer() *Player {
	mpg123.Initialize()
	player := Player{}
	player.Running = true
	player.Queue = make(chan string, 10)
	player.Lock = make(chan string)
	player.Error = make(chan error)
	return &player
}

func aoSampleFormat(handle *mpg123.Handle) *ao.SampleFormat {
	const bitsPerByte = 8

	rate, channels, encoding := handle.Format()

	return &ao.SampleFormat{
		BitsPerSample: handle.EncodingSize(encoding) * bitsPerByte,
		Rate:          int(rate),
		Channels:      channels,
		ByteFormat:    ao.FormatNative,
		Matrix:        nil,
	}
}
