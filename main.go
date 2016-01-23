package main

import (
	"log"
	"os"
	"os/user"
)

func main() {
	setup()
	video, _ := Search("king kunta")
	err := Download(video[0])
	log.Println(err)
}

func setup() {
	usr, _ := user.Current()

	os.MkdirAll(usr.HomeDir+"/.jable", 0777)

}
