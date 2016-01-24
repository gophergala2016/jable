package main

import (
	"fmt"
	"os"
)

func execAdd(args string) {
	video, err := Search(args, 1)
	if err != nil {
		handleErr(err)
		return
	}
	video.File, err = Download(video.ID)
	if err != nil {
		handleErr(err)
		return
	}
	player.Add(video)
}

func execPause() {
	player.Pause <- 1
}

func execResume() {
	player.Resume <- 1
}

func execSkip() {
	player.Skip <- 1
}

func execHelp() {
	bold.Printf("%15s\t", "help")
	fmt.Printf("%s\n", "Prints this dialog.")
	bold.Printf("%15s\t", "play QUERY")
	fmt.Printf("%s\n", "Plays the first result from the query, adds it to the queue if already playing.")
	bold.Printf("%15s\t", "skip")
	fmt.Printf("%s\n", "Skips the current song, starts playing the next in queue.")
	bold.Printf("%15s\t", "pause")
	fmt.Printf("%s\n", "Pauses the playback for the current song.")
	bold.Printf("%15s\t", "resume")
	fmt.Printf("%s\n", "Resumes the playback for the current song.")
	bold.Printf("%15s\t", "exit")
	fmt.Printf("%s", "Quit Jable.")
	println("")
}

func execExit() {
	player.Exit()
	cleanup()
	os.Exit(0)
}
