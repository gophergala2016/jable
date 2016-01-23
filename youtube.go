package main

import (
	"errors"
	"net/http"

	"google.golang.org/api/youtube/v3"

	"github.com/google/google-api-go-client/googleapi/transport"
)

// Video is a YouTube video
type Video struct {
	ID, Title, File string
}

// Search will find the first video for the search terms
func Search(term string) (*Video, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyBzqzgWz6_tucORR3NAGw9XC6qPq0ORanc"},
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(term).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	// Group video, channel, and playlist results in separate lists.
	videos := []Video{}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos = append(videos, Video{item.Id.VideoId, item.Snippet.Title, ""})
		}
	}

	if len(videos) < 1 {
		return nil, errors.New("No videos found")
	}
	return &videos[0], nil
}
