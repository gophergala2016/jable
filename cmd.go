package main

import (
	"fmt"
	"os"
)

func execPlay(args string) {
	video, err := Search(args)
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

func execAdd(args string) {

}

func execHelp() {
	bold.Printf("%15s\t", "help")
	fmt.Printf("%s\n", "Prints this dialog.")
	bold.Printf("%15s\t", "play QUERY")
	fmt.Printf("%s\n", "Finds the first result on YouTube and plays it.")
	bold.Printf("%15s\t", "add QUERY")
	fmt.Printf("%s\n", "Finds the first result on YouTube and adds it to the queue.")
	bold.Printf("%15s\t", "exit")
	fmt.Printf("%s\n", "Quit Jable.")
}

func execExit() {
	player.Stop()
	cleanup()
	os.Exit(0)
}
