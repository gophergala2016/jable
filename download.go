package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const downloadAPI = "http://www.youtubeinmp3.com/fetch/?format=JSON&video=http://www.youtube.com/watch?v=%s"

func Download(video *Video) error {
	resp, err := http.Get(fmt.Sprintf(downloadAPI, video.id))
	if err != nil {
		return err
	}
	var metadata map[string]string
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return err
	}
	resp, err = http.Get(metadata["link"])
	if err != nil {
		return err
	}
	file, err := os.Open(fmt.Sprintf("~/.jable/%s.mp3", video.id))
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
