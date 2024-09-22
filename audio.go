package goshazam

import (
	"bytes"
	"encoding/binary"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GenerateRawPCMInMemory(inputFile string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inputFile).
		Output("pipe:", ffmpeg.KwArgs{
			"f":      "s16le",
			"acodec": "pcm_s16le",
			"ar":     "16000",
			"ac":     "1",
		}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadSamplesFromBuffer(buf *bytes.Buffer) ([]int16, error) {
	samples := make([]int16, len(buf.Bytes())/2)
	err := binary.Read(buf, binary.LittleEndian, samples)
	if err != nil {
		return nil, err
	}
	return samples, nil
}
