package invidious

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Video struct {
	Rtype          string           `json:"type"`
	Title          string           `json:"title"`
	VideoId        string           `json:"videoId"`
	Author         string           `json:"author"`
	Length         int              `json:"lengthSeconds"`
	Format         []format         `json:"adaptiveFormats"`
	VideoThumbnail []videoThumbnail `json:"videoThumbnails"`
}

type format struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

type videoThumbnail struct {
	Url string `json:"url"`
}

type instances struct {
	Api bool   `json:"api"`
	Url string `json:"uri"`
}

type suggestion struct {
	Query       string   `json:"query"`
	Suggestions []string `json:"suggestions"`
}

func GetSuggestions(instance, query string) ([]string, error) {
	var sg suggestion
	var url string = fmt.Sprintf("%s/api/v1/search/suggestions?q=%s",
		instance, url.QueryEscape(query))

	var err = request(url, &sg)
	if err != nil {
		return nil, fmt.Errorf("GetSuggestions: %w", err)
	}

	return sg.Suggestions, nil
}

func Search(instance, query, sType string, page int) ([]Video, error) {
	var url string = fmt.Sprintf("%s/api/v1/search?q=%s&type=%s&page=%d",
		instance, url.QueryEscape(query), sType, page)

	var videos = []Video{}
	var err = request(url, &videos)
	if err != nil {
		return nil, fmt.Errorf("Search: %w", err)
	}

	return videos, nil
}

func GetVideo(instance, videoId string) Video {
	var url string = fmt.Sprintf("%s/api/v1/videos/%s",
		instance, url.QueryEscape(videoId))

	var video = Video{}
	var err = request(url, &video)
	if err != nil {
		log.Println(err)
	}

	return video
}

func GetInstances() []string {
	var instance []string
	var ins = [][]instances{}
	var url string = "https://api.invidious.io/instances.json"

	request(url, &ins)

	for _, v := range ins {
		if v[1].Api == true {
			instance = append(instance, v[1].Url)
		}
	}
	return instance
}

func request(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
