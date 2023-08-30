package invidious

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/image/webp"
)

type Video struct {
	Title            string         `json:"title"`
	VideoId          string         `json:"url"`
	Uploader         string         `json:"uploaderName"`
	Duration         int            `json:"duration"`
	Thumbnail        string         `json:"thumbnail"`
	ThumbnailUrl     string         `json:"thumbnailUrl"`
	AudioStreams     []audioStreams `json:"audioStreams"`
	FormatedDuration string
}

type audioStreams struct {
	Url     string `json:"url"`
	Quality string `json:"quality"`
}

type items struct {
	Items []Video `json:"items"`
}

func (v *Video) GetThumbnail() image.Image {
	var url string = v.Thumbnail
	if url == "" {
		url = v.ThumbnailUrl
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	image, err := webp.Decode(bufio.NewReader(resp.Body))
	if err != nil {
		log.Println(err)
	}

	return image
}

func GetSuggestions(instance, query string) ([]string, error) {
	var suggestions []string
	var url string = fmt.Sprintf("%s/suggestions?query=%s",
		instance, url.QueryEscape(query))

	var err = request("json", url, &suggestions)
	if err != nil {
		return nil, fmt.Errorf("GetSuggestions: %w", err)
	}

	return suggestions, nil
}

func Search(instance, query, sType string) ([]Video, error) {
	var url string = fmt.Sprintf("%s/search?q=%s&filter=%s",
		instance, url.QueryEscape(query), sType)

	var items = items{}
	var err = request("json", url, &items)
	if err != nil {
		return nil, fmt.Errorf("Search: %w", err)
	}

	for i, v := range items.Items {
		items.Items[i].VideoId = strings.Split(v.VideoId, "v=")[1]
		items.Items[i].FormatedDuration = getDuration(v.Duration)
	}

	return items.Items, nil
}

func GetVideo(instance, videoId string) Video {
	var url string = fmt.Sprintf("%s/streams/%s",
		instance, url.QueryEscape(videoId))

	var video = Video{}
	var err = request("json", url, &video)
	if err != nil {
		log.Println(err)
	}

	video.FormatedDuration = getDuration(video.Duration)

	return video
}

func GetInstances() []string {
	var instances []string
	var url string = "https://raw.githubusercontent.com/wiki/TeamPiped/Piped-Frontend/Instances.md"

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var api bool = false
	for _, v := range strings.Split(string(content), "\n") {
		if api && strings.Contains(v, "|") {
			instances = append(instances, strings.TrimSpace(strings.Split(v, "|")[1]))
		}
		if strings.Contains(v, "---") {
			api = true
		}
	}

	return instances
}

func request(t, url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if t == "json" {
		err = json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			return fmt.Errorf("request: %w", err)
		}
	} else if t == "string" {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("request: %w", err)
		}
		v = string(content)
	}

	return nil
}

func getDuration(duration int) string {
	if duration/60 > 60 {
		return fmt.Sprintf("%02d:%02d:%02d",
			duration/60/60, duration/60%60, duration%60)
	}
	return fmt.Sprintf("%02d:%02d", duration/60, duration%60)
}
