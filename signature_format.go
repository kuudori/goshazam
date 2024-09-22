package goshazam

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"sort"
)

const DataURIPrefix = "data:audio/vnd.shazam.sig;base64,"

type RawSignatureHeader struct {
	Magic1                             uint32
	CRC32                              uint32
	SizeMinusHeader                    uint32
	Magic2                             uint32
	Void1                              [12]byte
	ShiftedSampleRateID                uint32
	Void2                              [8]byte
	NumberSamplesPlusDividedSampleRate uint32
	FixedValue                         uint32
}

func (ds *DecodedSignature) EncodeToBinary() ([]byte, error) {
	header := RawSignatureHeader{
		Magic1:     0xCAFE2580,
		Magic2:     0x94119C00,
		FixedValue: (15 << 19) + 0x40000,
	}

	switch ds.SampleRateHz {
	case 8000:
		header.ShiftedSampleRateID = 1 << 27
	case 11025:
		header.ShiftedSampleRateID = 2 << 27
	case 16000:
		header.ShiftedSampleRateID = 3 << 27
	case 32000:
		header.ShiftedSampleRateID = 4 << 27
	case 44100:
		header.ShiftedSampleRateID = 5 << 27
	case 48000:
		header.ShiftedSampleRateID = 6 << 27
	default:
		return nil, fmt.Errorf("invalid sample rate passed when encoding Shazam packet")
	}

	header.NumberSamplesPlusDividedSampleRate = ds.NumberSamples + uint32(float32(ds.SampleRateHz)*0.24)

	sortedBands := make([]FrequencyBand, 0, len(ds.FrequencyBandToSoundPeaks))
	for band := range ds.FrequencyBandToSoundPeaks {
		sortedBands = append(sortedBands, band)
	}
	sort.Slice(sortedBands, func(i, j int) bool {
		return sortedBands[i] < sortedBands[j]
	})

	var contentsBuf bytes.Buffer
	contentsBuf.Grow(len(ds.FrequencyBandToSoundPeaks) * 1024)

	for _, band := range sortedBands {
		peaks := ds.FrequencyBandToSoundPeaks[band]

		var peaksBuf bytes.Buffer
		peaksBuf.Grow(len(peaks) * 8)
		var fftPassNumber uint32

		for _, peak := range peaks {
			if peak.FFTPassNumber-fftPassNumber >= 255 {
				peaksBuf.WriteByte(0xff)
				binary.Write(&peaksBuf, binary.LittleEndian, uint32(peak.FFTPassNumber))
				fftPassNumber = peak.FFTPassNumber
			}

			peaksBuf.WriteByte(byte(peak.FFTPassNumber - fftPassNumber))
			binary.Write(&peaksBuf, binary.LittleEndian, uint16(peak.PeakMagnitude))
			binary.Write(&peaksBuf, binary.LittleEndian, peak.CorrectedPeakFrequencyBin)

			fftPassNumber = peak.FFTPassNumber
		}

		peaksBytes := peaksBuf.Bytes()

		binary.Write(&contentsBuf, binary.LittleEndian, 0x60030040+uint32(band))
		binary.Write(&contentsBuf, binary.LittleEndian, uint32(len(peaksBytes)))
		contentsBuf.Write(peaksBytes)
		paddingSize := (4 - len(peaksBytes)%4) % 4
		contentsBuf.Write(make([]byte, paddingSize))
	}

	header.SizeMinusHeader = uint32(contentsBuf.Len() + 8)

	var buf bytes.Buffer
	buf.Grow(binary.Size(header) + contentsBuf.Len() + 8)
	binary.Write(&buf, binary.LittleEndian, header)
	binary.Write(&buf, binary.LittleEndian, uint32(0x40000000))
	binary.Write(&buf, binary.LittleEndian, header.SizeMinusHeader)
	buf.Write(contentsBuf.Bytes())

	bufferBytes := buf.Bytes()
	header.CRC32 = crc32.ChecksumIEEE(bufferBytes[8:])

	buf.Reset()
	buf.Grow(len(bufferBytes))
	binary.Write(&buf, binary.LittleEndian, header)
	buf.Write(bufferBytes[binary.Size(header):])

	return buf.Bytes(), nil
}

func (ds *DecodedSignature) EncodeToURI() (string, error) {
	b, err := ds.EncodeToBinary()
	if err != nil {
		return "", err
	}
	return DataURIPrefix + base64.StdEncoding.EncodeToString(b), nil
}
