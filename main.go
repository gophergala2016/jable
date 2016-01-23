package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
)

var userDir string

func main() {
	setup()
	player := NewPlayer()
	//player.Start()

	// player.Start()
	// video, _ := Search("smoke on the water")
	// err := Download(video[0])
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("jable: ")
	for scanner.Scan() {
		cmd, args := parseCmd(scanner.Text())
		switch cmd {
		case "play":
			video, err := Search(args)
			video.File, err = Download(video.ID)
			handleErr(err)
			player.Play(fmt.Sprintf("%s/.jable/%s.mp3", userDir, video.ID))
		case "exit":
			player.Stop()
			os.Exit(0)
		}
		fmt.Print("jable: ")
	}
}

func setup() {
	usr, _ := user.Current()
	userDir = usr.HomeDir
	os.MkdirAll(userDir+"/.jable", 0777)
}

func parseCmd(cmd string) (string, string) {
	args := strings.Split(cmd, " ")
	return args[0], strings.Join(args[1:], "")
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
