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

var player *Player

var red, cyan, bold *color.Color

func main() {
	setup()
	defer cleanup()
	player = NewPlayer()
	player.Start()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd, args := parseCmd(scanner.Text())
		switch cmd {
		case "play":
			execPlay(args)
		case "help":
			execHelp()
		case "exit":
			execExit()
		default:
			red.Printf("Invalid command ")
			bold.Printf("%s %s\n", cmd, args)
		}
		bold.Print("jable: ")
	}
}

func setup() {
	usr, _ := user.Current()
	userDir = usr.HomeDir
	os.MkdirAll(userDir+"/.jable", 0777)

	red = color.New(color.FgRed)
	cyan = color.New(color.FgCyan)
	bold = color.New(color.Bold)

	cyan.Println("Hello there! I am Jable, and I will play you some music!")
	fmt.Print("Type")
	bold.Print(" help ")
	fmt.Println("for list of commands.")
	bold.Print("jable: ")
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
		red.Println(err)
	}
}
