package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/fatih/color"
)

var userDir string

var red, green, blue, bold *color.Color

func main() {
	setup()
	defer cleanup()
	player := NewPlayer()
	//player.Start()

	// player.Start()
	// video, _ := Search("smoke on the water")
	// err := Download(video[0])
	scanner := bufio.NewScanner(os.Stdin)
	blue.Println("Hello there! I am Jabble, and I will play you some music!")
	fmt.Println(`usage:`)
	fmt.Println(`\t\t play [search terms]`)
	fmt.Println(`\t\t exit`)
	fmt.Print("jable: ")
	for scanner.Scan() {
		cmd, args := parseCmd(scanner.Text())
		switch cmd {
		case "play":
			video, err := Search(args)
			handleErr(err)
			video.File, err = Download(video.ID)
			handleErr(err)
			player.Play(fmt.Sprintf("%s/.jable/%s.mp3", userDir, video.ID))
		case "exit":
			player.Stop()
			cleanup()
			os.Exit(0)
		}
		fmt.Print("jable: ")
	}
}

func setup() {
	usr, _ := user.Current()
	userDir = usr.HomeDir
	os.MkdirAll(userDir+"/.jable", 0777)

	red = color.New(color.FgRed)
	blue = color.New(color.FgBlue)
	green = color.New(color.FgGreen)
	bold = color.New(color.Bold)
}

func cleanup() {
	os.RemoveAll(userDir + "/.jable")
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
