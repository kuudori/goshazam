package goshazam

import (
	"encoding/json"
	"fmt"
	"time"
)

type RecognizeResponse struct {
	Matches   []Match  `json:"matches"`
	Location  Location `json:"location"`
	Timestamp int64    `json:"timestamp"`
	Timezone  string   `json:"timezone"`
	Track     Track    `json:"track"`
	TagID     string   `json:"tagid"`
}

type Match struct {
	ID            string  `json:"id"`
	Offset        float64 `json:"offset"`
	TimeSkew      float64 `json:"timeskew"`
	FrequencySkew float64 `json:"frequencyskew"`
}

type Location struct {
	Accuracy float64 `json:"accuracy"`
}

type Track struct {
	Layout    string    `json:"layout"`
	Type      string    `json:"type"`
	Key       string    `json:"key"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Images    Images    `json:"images"`
	Share     Share     `json:"share"`
	Hub       Hub       `json:"hub"`
	Sections  []Section `json:"sections"`
	URL       string    `json:"url"`
	Artists   []Artist  `json:"artists"`
	ISRC      string    `json:"isrc"`
	Genres    Genres    `json:"genres"`
	URLParams URLParams `json:"urlparams"`
	MyShazam  MyShazam  `json:"myshazam"`
}

type Images struct {
	Background string `json:"background"`
	CoverArt   string `json:"coverart"`
	CoverArtHQ string `json:"coverarthq"`
	JoeColor   string `json:"joecolor"`
}

type Share struct {
	Subject  string `json:"subject"`
	Text     string `json:"text"`
	Href     string `json:"href"`
	Image    string `json:"image"`
	Twitter  string `json:"twitter"`
	HTML     string `json:"html"`
	Avatar   string `json:"avatar"`
	Snapchat string `json:"snapchat"`
}

type Hub struct {
	Type        string     `json:"type"`
	Image       string     `json:"image"`
	Actions     []Action   `json:"actions"`
	Options     []Option   `json:"options"`
	Providers   []Provider `json:"providers"`
	Explicit    bool       `json:"explicit"`
	DisplayName string     `json:"displayname"`
}

type Action struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
	URI  string `json:"uri,omitempty"`
}

type Option struct {
	Caption             string     `json:"caption"`
	Actions             []Action   `json:"actions"`
	BeaconData          BeaconData `json:"beacondata"`
	Image               string     `json:"image"`
	Type                string     `json:"type"`
	ListCaption         string     `json:"listcaption"`
	OverflowImage       string     `json:"overflowimage"`
	ColourOverflowImage bool       `json:"colouroverflowimage"`
	ProviderName        string     `json:"providername"`
}

type BeaconData struct {
	Type         string `json:"type"`
	ProviderName string `json:"providername"`
}

type Provider struct {
	Caption string         `json:"caption"`
	Images  ProviderImages `json:"images"`
	Actions []Action       `json:"actions"`
	Type    string         `json:"type"`
}

type ProviderImages struct {
	Overflow string `json:"overflow"`
	Default  string `json:"default"`
}

type Section struct {
	Type      string     `json:"type"`
	MetaPages []MetaPage `json:"metapages,omitempty"`
	TabName   string     `json:"tabname"`
	Metadata  []Metadata `json:"metadata,omitempty"`
	URL       string     `json:"url,omitempty"`
}

type MetaPage struct {
	Image   string `json:"image"`
	Caption string `json:"caption"`
}

type Metadata struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Artist struct {
	ID     string `json:"id"`
	AdamID string `json:"adamid"`
}

type Genres struct {
	Primary string `json:"primary"`
}

type URLParams struct {
	TrackTitle  string `json:"tracktitle"`
	TrackArtist string `json:"trackartist"`
}

type MyShazam struct {
	Apple AppleMyShazam `json:"apple"`
}

type AppleMyShazam struct {
	Actions []Action `json:"actions"`
}

type RecognizeResult struct {
	rawData json.RawMessage
}

// Serialize converts JSON to RecognizeResponse
func (r *RecognizeResult) Serialize() (*RecognizeResponse, error) {
	var response RecognizeResponse
	if err := json.Unmarshal(r.rawData, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	response.Timestamp = time.Now().Unix()
	return &response, nil
}

// Raw returns raw JSON
func (r *RecognizeResult) Raw() json.RawMessage {
	return r.rawData
}
