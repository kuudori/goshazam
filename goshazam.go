package goshazam

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// ShazamClient is a client for recognizing music using the Shazam service.
type ShazamClient struct {
	client     *http.Client
	headers    http.Header
	userAgents [12]string
	randMu     sync.Mutex
	rand       *rand.Rand
}

func NewShazamClient() *ShazamClient {
	return &ShazamClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: http.Header{
			"Content-Type":     {"application/json"},
			"Content-Language": {"en_US"},
		},
		userAgents: userAgents,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *ShazamClient) getRandomUserAgent() string {
	c.randMu.Lock()
	defer c.randMu.Unlock()
	return c.userAgents[c.rand.Intn(len(c.userAgents))]
}

func (c *ShazamClient) Do(req *http.Request) (*http.Response, error) {
	for key, values := range c.headers {
		req.Header[key] = values
	}
	req.Header.Set("User-Agent", c.getRandomUserAgent())
	return c.client.Do(req)
}

// Recognize processes an audio file and returns the recognition result.
func (c *ShazamClient) Recognize(ctx context.Context, filePath string) (*RecognizeResult, error) {
	rawPCM, err := GenerateRawPCMInMemory(filePath)
	if err != nil {
		return nil, fmt.Errorf("error generating raw PCM: %w", err)
	}

	samples, err := ReadSamplesFromBuffer(rawPCM)
	if err != nil {
		return nil, fmt.Errorf("error reading samples from buffer: %w", err)
	}

	sg := NewSignatureGenerator()
	signature := sg.MakeSignatureFromBuffer(samples)

	data, err := GetSignatureJSON(&signature)
	if err != nil {
		return nil, fmt.Errorf("error getting signature JSON: %w", err)
	}

	result, err := c.sendShazamRecognitionRequest(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("error recognizing from voice: %w", err)
	}

	return &RecognizeResult{rawData: result}, nil
}
