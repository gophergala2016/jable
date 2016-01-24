package main

import (
	"fmt"
	"io"

	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
)

// Player handles audio player runtime
type Player struct {
	Lock                chan string
	Skip, Pause, Resume chan int
	Queue               chan Video
	Running, Playing    bool
}

// Add will add video to the queue
func (p *Player) Add(video *Video) {
	p.Queue <- *video
	println(fmt.Sprintf("%s added to the queue.", video.Title))
}

// Start will start playing the songs from the queue
func (p *Player) Start() {
	go func(p *Player) {
		for p.Running {
			select {
			case video := <-p.Queue:
				p.Lock <- video.ID
				println(fmt.Sprintf("\u266B Now playing %s", video.Title))
				returned()
				file := fmt.Sprintf("%s/.jable/%s.mp3", userDir, video.ID)
				p.Play(file)
				<-p.Lock
			}
		}
	}(p)
}

// Play starts playback from specified file
func (p *Player) Play(file string) {
	p.Playing = true
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
		p.Playing = false
		return
	}
	p.Playing = false
}

// Exit will exit the player
func (p *Player) Exit() {
	p.Running = false
	mpg123.Exit()
}

// NewPlayer returns a new Player instance
func NewPlayer() *Player {
	mpg123.Initialize()
	player := Player{}
	player.Running = true
	player.Queue = make(chan Video, 10)
	player.Lock = make(chan string, 1)
	player.Skip = make(chan int, 1)
	player.Pause = make(chan int, 1)
	player.Resume = make(chan int, 1)
	player.Playing = false
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
L:
	for {
		select {
		case <-p.Skip:
			println("Skipped.")
			returned()
			return 0, nil
		case <-p.Pause:
			println("Paused.")
			returned()
			<-p.Resume
			println("Resuming...")
			returned()
		default:
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if nw > 0 {
					written += int64(nw)
				}
				if ew != nil {
					err = ew
					break L
				}
				if nr != nw {
					err = io.ErrShortWrite
					break L
				}
			}
			if er == io.EOF {
				break L
			}
			if er != nil {
				err = er
				break L
			}
		}
	}
	return written, err
}
