package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
)

const downloadAPI = "http://www.youtubeinmp3.com/fetch/?format=JSON&video=http://www.youtube.com/watch?v=%s"

func Download(video Video) error {
	resp, err := http.Get(fmt.Sprintf(downloadAPI, video.Id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var metadata map[string]string
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return err
	}
	fmt.Println(metadata)
	fileData, err := http.Get(metadata["link"])
	if err != nil {
		return err
	}
	defer fileData.Body.Close()
	usr, err := user.Current()

	file, err := os.Create(fmt.Sprintf("%s/.jable/%s.mp3", usr.HomeDir, video.Id))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, fileData.Body)
	if err != nil {
		return err
	}
	return nil
}
