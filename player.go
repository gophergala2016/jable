package main

import (
	"io"

	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
)

type Player struct {
}

func (p *Player) Add(video Video) {

}

func (p *Player) Play(file string) error {
	handle, err := mpg123.Open(file)
	if err != nil {
		return err
	}
	defer handle.Close()

	ao.Initialize()
	defer ao.Shutdown()
	dev := ao.NewLiveDevice(AoSampleFormat(handle))
	defer dev.Close()

	if _, err := io.Copy(dev, handle); err != nil {
		return err
	}
	return nil
}

func (p *Player) Stop() {
	mpg123.Exit()
}

func NewPlayer() *Player {
	mpg123.Initialize()
	return &Player{}
}

func AoSampleFormat(handle *mpg123.Handle) *ao.SampleFormat {
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
