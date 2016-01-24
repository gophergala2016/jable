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
	go func(p *Player) {
		fmt.Println("Player starting...")
		for p.Running {
			id := <-p.Queue
			fmt.Printf("Got video %s\n", id)
			p.Lock <- id
			fmt.Println("Successfully locked player.")
			file := fmt.Sprintf("%s/.jable/%s.mp3", userDir, id)
			p.Play(file)
			fmt.Println("Playback finished, releasing lock...")
			<-p.Lock
		}
	}(p)
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

	if _, err := p.copyBuffer(dev, handle); err != nil {
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
	player.Lock = make(chan string, 1)
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

func (p *Player) copyBuffer(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
