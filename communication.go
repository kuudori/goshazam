package goshazam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type GeolocationResponse struct {
	Altitude  float64 `json:"altitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SignatureSong struct {
	Samples   uint32 `json:"samples"`
	Timestamp uint32 `json:"timestamp"`
	URI       string `json:"uri"`
}

type Signature struct {
	Geolocation GeolocationResponse `json:"geolocation"`
	Signature   SignatureSong       `json:"signature"`
	Timestamp   uint32              `json:"timestamp"`
	Timezone    string              `json:"timezone"`
}

func GetSignatureJSON(signature *DecodedSignature) (*Signature, error) {
	timestampMs := uint32(time.Now().UnixNano() / int64(time.Millisecond))
	samples := uint32(float32(signature.NumberSamples) / float32(signature.SampleRateHz) * 1000)

	uri, err := signature.EncodeToURI()
	if err != nil {
		return nil, fmt.Errorf("failed to encode signature to URI: %w", err)
	}

	return &Signature{
		Geolocation: GeolocationResponse{
			Altitude:  rand.Float64()*400 + 100,
			Latitude:  rand.Float64()*180 - 90,
			Longitude: rand.Float64()*360 - 180,
		},
		Signature: SignatureSong{
			Samples:   samples,
			Timestamp: timestampMs,
			URI:       uri,
		},
		Timestamp: timestampMs,
		Timezone:  "Europe/Paris",
	}, nil
}

func RecognizeFromVoice(ctx context.Context, requestData interface{}) (interface{}, error) {
	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()
	url := fmt.Sprintf("https://amp.shazam.com/discovery/v5/en-US/GB/iphone/-/tag/%s/%s", uuid1, uuid2)

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var userAgents = []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/85.0.4183.109 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/29.0 Mobile/15E148 Safari/605.1.15",
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	req.Header.Set("Content-Language", "en_US")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return result, nil
}
