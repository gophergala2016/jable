package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

const checkURI = "https://d.yt-downloader.org/check.php?callback=jable&v=%s&f=mp3&_=%d"
const convertURI = "https://d.yt-downloader.org/progress.php?callback=jable&id=%s&_=%d"
const downloadURI = "http://%s.yt-downloader.org/download.php?id=%s"

var errDownload = errors.New("There was an error processing the song, please try again.")

var serverIDs = map[string]string{
	"1":  "gpkio",
	"2":  "hpbnj",
	"3":  "macsn",
	"4":  "pikku",
	"5":  "fgkzc",
	"6":  "hmqbu",
	"7":  "kyhxj",
	"8":  "nwwxj",
	"9":  "sbist",
	"10": "ditrj",
	"11": "qypbr",
	"12": "wiyqr",
	"13": "xxvcy",
	"14": "afyzk",
	"15": "kjzmv",
	"16": "txrys",
	"17": "kzrzi",
	"18": "rmira",
	"19": "umbbo",
	"20": "aigkk",
	"21": "qgxhg",
	"22": "twrri",
	"23": "fkaph",
}

// Download video with the given id and return the filepath it is stored in
func Download(id string) (string, error) {
	metadata := check(id)
	println("Fetching, converting, buffering...")
	if metadata["sid"] != "0" && metadata["ce"] != "0" {
		return fetch(id, metadata["hash"], serverIDs[metadata["sid"]])
	}
	return convert(id, metadata["hash"])
}

func fetch(videoID, hash, sid string) (string, error) {
	fileData, err := http.Get(fmt.Sprintf(downloadURI, sid, hash))
	if err != nil {
		return "", errDownload
	}
	defer fileData.Body.Close()

	filename := fmt.Sprintf("%s/.jable/%s.mp3", userDir, videoID)
	file, err := os.Create(fmt.Sprintf("%s/.jable/%s.mp3", userDir, videoID))
	if err != nil {
		return "", errDownload
	}
	defer file.Close()
	_, err = io.Copy(file, fileData.Body)
	if err != nil {
		return "", errDownload
	}
	return filename, nil
}

func convert(id, hash string) (string, error) {
	metadata := jsonPRequest(fmt.Sprintf(convertURI, hash, time.Now().Unix()))
	if metadata["progress"] == "3" {
		return fetch(id, hash, serverIDs[metadata["sid"]])
	}
	if metadata["error"] != "" {
		return "", errDownload
	}
	time.Sleep(3 * time.Second)
	return convert(id, hash)

}

func check(id string) map[string]string {
	return jsonPRequest(fmt.Sprintf(checkURI, id, time.Now().Unix()))
}

func jsonPRequest(url string) map[string]string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	r, _ := regexp.Compile(`jable\((\{.*\})\)`)

	jsonp, _ := ioutil.ReadAll(resp.Body)

	metadataString := r.FindStringSubmatch(string(jsonp))

	var metadata map[string]string
	json.Unmarshal([]byte(metadataString[1]), &metadata)

	return metadata
}
