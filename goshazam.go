package goshazam

import (
	"context"
	"encoding/json"
	"fmt"
)

// Recognize takes the context and path to the audio file and returns the result of recognition in JSON format
func Recognize(ctx context.Context, filePath string) (json.RawMessage, error) {
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

	result, err := RecognizeFromVoice(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("error recognizing from voice: %w", err)
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("error marshalling result to JSON: %w", err)
	}

	return jsonResult, nil
}
