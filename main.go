package main

import (
	"bufio"
	"errors"
	"os"
	"os/signal"
	"os/user"
	"strings"

	"github.com/fatih/color"
)

var userDir string

var player *Player

var red, cyan, bold *color.Color

func main() {
	setup()
	defer cleanup()
	player = NewPlayer()
	player.Start()

	scanner := bufio.NewScanner(os.Stdin)
	returned()
	for scanner.Scan() {
		cmd, args := parseCmd(scanner.Text())
		switch cmd {
		case "skip":
			execSkip()
			returned()
		case "play":
			execAdd(args)
		case "pause":
			execPause()
			returned()
		case "resume":
			execResume()
			returned()
		case "exit":
			execExit()
			returned()
		case "help":
			execHelp()
			returned()
		default:
			handleErr(errors.New("Invalid command"))
		}
	}
}

func setup() {
	usr, _ := user.Current()
	userDir = usr.HomeDir
	os.MkdirAll(userDir+"/.jable", 0777)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			execExit()
		}
	}()

	red = color.New(color.FgRed)
	cyan = color.New(color.FgCyan)
	bold = color.New(color.Bold)

	cyan.Println(" \u266B Hello there! I am Jable, and I will play you some music!\u266B")
	println("Type 'help' for a list of commands.")

}

func cleanup() {
	os.RemoveAll(userDir + "/.jable")
}

func parseCmd(cmd string) (string, string) {
	args := strings.Split(cmd, " ")
	return args[0], strings.Join(args[1:], "")
}

func handleErr(err error) {
	red.Printf("\n%s\n", err)
}

func println(str string) {
	bold.Printf("\n%s\n", str)
}

func returned() {
	bold.Print("jable: ")
}
