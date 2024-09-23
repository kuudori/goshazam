package goshazam

import (
	"gonum.org/v1/gonum/dsp/fourier"
	"math"
	"sync"
)

var hannWindow []float64

func init() {
	hannWindow = make([]float64, fftSize)
	N := len(hannWindow)
	for i := 0; i < N; i++ {
		hannWindow[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(N-1)))
	}
}

type FrequencyBand int

type FrequencyPeak struct {
	FFTPassNumber             uint32
	PeakMagnitude             float64
	CorrectedPeakFrequencyBin uint16
	SampleRateHz              uint32
}

type DecodedSignature struct {
	SampleRateHz              uint32
	NumberSamples             uint32
	FrequencyBandToSoundPeaks map[FrequencyBand][]FrequencyPeak
}

type SignatureGenerator struct {
	ringBufferOfSamples          []int16
	reorderedRingBufferOfSamples []float64
	fftOutputs                   [][]float64
	spreadFFTOutputs             [][]float64
	ringBufferOfSamplesIndex     int
	fftOutputsIndex              int
	spreadFFTOutputsIndex        int
	numSpreadFFTsDone            uint32
	signature                    DecodedSignature
	mu                           sync.Mutex
}

func NewSignatureGenerator() *SignatureGenerator {
	s := &SignatureGenerator{
		ringBufferOfSamples:          make([]int16, fftSize),
		reorderedRingBufferOfSamples: make([]float64, fftSize),
		fftOutputs:                   make([][]float64, numFFTs),
		spreadFFTOutputs:             make([][]float64, numFFTs),
		signature: DecodedSignature{
			FrequencyBandToSoundPeaks: make(map[FrequencyBand][]FrequencyPeak),
		},
	}
	for i := range s.fftOutputs {
		s.fftOutputs[i] = make([]float64, fftOutputSize)
	}
	for i := range s.spreadFFTOutputs {
		s.spreadFFTOutputs[i] = make([]float64, fftOutputSize)
	}
	return s
}

func (s *SignatureGenerator) MakeSignatureFromBuffer(s16Mono16kHzBuffer []int16) DecodedSignature {
	s.signature = DecodedSignature{
		SampleRateHz:              sampleRate,
		FrequencyBandToSoundPeaks: make(map[FrequencyBand][]FrequencyPeak),
	}

	maxSamples := int(maxTimeSeconds * float64(sampleRate))
	if len(s16Mono16kHzBuffer) > maxSamples {
		s16Mono16kHzBuffer = s16Mono16kHzBuffer[:maxSamples]
	}
	s.signature.NumberSamples = uint32(len(s16Mono16kHzBuffer))

	for i := 0; i+128 <= len(s16Mono16kHzBuffer); i += 128 {
		chunk := s16Mono16kHzBuffer[i : i+128]
		s.doFFT(chunk)
		s.doPeakSpreading()
		s.numSpreadFFTsDone++

		if s.numSpreadFFTsDone >= 46 {
			s.doPeakRecognition()
		}

	}
	return s.signature
}

func (s *SignatureGenerator) doFFT(s16Mono16kHzBuffer []int16) {
	for i := 0; i < len(s16Mono16kHzBuffer); i++ {
		s.ringBufferOfSamples[(s.ringBufferOfSamplesIndex+i)%fftSize] = s16Mono16kHzBuffer[i]
	}
	s.ringBufferOfSamplesIndex = (s.ringBufferOfSamplesIndex + len(s16Mono16kHzBuffer)) % fftSize

	startIndex := (s.ringBufferOfSamplesIndex + fftSize - fftSize) % fftSize

	for i := 0; i < fftSize; i++ {
		s.reorderedRingBufferOfSamples[i] = float64(s.ringBufferOfSamples[(startIndex+i)%fftSize]) * hannWindow[i]
	}
	fft := fourier.NewFFT(fftSize)
	complexFFTResults := fft.Coefficients(nil, s.reorderedRingBufferOfSamples)
	realFFTResults := s.fftOutputs[s.fftOutputsIndex]

	for i := 0; i < fftOutputSize; i++ {
		realPart := real(complexFFTResults[i])
		imagPart := imag(complexFFTResults[i])
		magnitudeSquared := (realPart*realPart + imagPart*imagPart) / float64(1<<17)
		realFFTResults[i] = math.Max(magnitudeSquared, 1e-10)
	}
	s.fftOutputsIndex = (s.fftOutputsIndex + 1) % numFFTs
}

func (s *SignatureGenerator) doPeakSpreading() {
	fftOutputsPosition := (s.fftOutputsIndex - 1 + numFFTs) % numFFTs
	originLastFFT := s.fftOutputs[fftOutputsPosition]

	temporaryArray1 := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		temporaryArray1[i] = originLastFFT
	}
	temporaryArray1[1] = append(temporaryArray1[1][1:], temporaryArray1[1][:1]...)
	temporaryArray1[2] = append(temporaryArray1[2][2:], temporaryArray1[2][:2]...)

	originLastFFTSp := make([]float64, fftOutputSize)
	for i := 0; i < fftOutputSize-3; i++ {
		originLastFFTSp[i] = math.Max(
			temporaryArray1[0][i],
			math.Max(temporaryArray1[1][i], temporaryArray1[2][i]),
		)
	}
	copy(originLastFFTSp[fftOutputSize-3:], originLastFFT[fftOutputSize-3:])

	offsets := []int{-1, -3, -6}
	positions := make([]int, len(offsets))
	for i, offset := range offsets {
		positions[i] = (s.spreadFFTOutputsIndex + offset + numFFTs) % numFFTs
	}

	temporaryArray2 := make([][]float64, 4)
	temporaryArray2[0] = originLastFFTSp
	for i, pos := range positions {
		temporaryArray2[i+1] = s.spreadFFTOutputs[pos]
	}

	for i := 1; i <= 3; i++ {
		for j := 0; j < fftOutputSize; j++ {
			temporaryArray2[i][j] = math.Max(temporaryArray2[i][j], temporaryArray2[i-1][j])
		}
		s.spreadFFTOutputs[positions[i-1]] = temporaryArray2[i]
	}

	s.spreadFFTOutputs[s.spreadFFTOutputsIndex] = originLastFFTSp
	s.spreadFFTOutputsIndex = (s.spreadFFTOutputsIndex + 1) % numFFTs
}

func (s *SignatureGenerator) doPeakRecognition() {
	fftMinus46 := s.fftOutputs[(s.fftOutputsIndex-46+numFFTs)%numFFTs]
	fftMinus49 := s.spreadFFTOutputs[(s.spreadFFTOutputsIndex-49+numFFTs)%numFFTs]

	for binPosition := 10; binPosition <= 1014; binPosition++ {
		isMagnitudeAboveThreshold := fftMinus46[binPosition] >= minPeakMagnitude
		isLocalMax := fftMinus46[binPosition] >= fftMinus49[binPosition-1]
		if isMagnitudeAboveThreshold && isLocalMax {
			var maxNeighborInFFTMinus49 float64
			for _, offset := range []int{-10, -7, -4, -3, 1, 2, 5, 8} {
				neighborIndex := binPosition + offset
				if neighborIndex >= 0 && neighborIndex < fftOutputSize {
					maxNeighborInFFTMinus49 = math.Max(maxNeighborInFFTMinus49, fftMinus49[neighborIndex])
				}
			}
			if fftMinus46[binPosition] > maxNeighborInFFTMinus49 {
				maxNeighborInOtherAdjacentFFTs := maxNeighborInFFTMinus49
				for _, offset := range []int{-53, -45, 165, 172, 179, 186, 193, 200, 214, 221, 228, 235, 242, 249} {
					idx := (s.spreadFFTOutputsIndex + offset + numFFTs) % numFFTs
					otherFFT := s.spreadFFTOutputs[idx]
					binIdx := binPosition - 1
					if binIdx >= 0 && binIdx < fftOutputSize {
						maxNeighborInOtherAdjacentFFTs = math.Max(maxNeighborInOtherAdjacentFFTs, otherFFT[binIdx])
					}
				}
				if fftMinus46[binPosition] > maxNeighborInOtherAdjacentFFTs {
					fftPassNumber := s.numSpreadFFTsDone - 46
					peakMagnitude := math.Log(math.Max(minPeakMagnitude, fftMinus46[binPosition]))*1477.4 + 6144.0

					peakMagnitudeBefore := math.Log(math.Max(minPeakMagnitude, fftMinus46[binPosition-1]))*1477.4 + 6144.0
					peakMagnitudeAfter := math.Log(math.Max(minPeakMagnitude, fftMinus46[binPosition+1]))*1477.4 + 6144.0

					peakVariation1 := peakMagnitude*2.0 - peakMagnitudeBefore - peakMagnitudeAfter
					if peakVariation1 <= 0 {
						continue
					}
					peakVariation2 := (peakMagnitudeAfter - peakMagnitudeBefore) * 32.0 / peakVariation1
					correctedPeakFrequencyBin := uint16(binPosition*64) + uint16(peakVariation2+0.5)

					frequencyHz := float64(correctedPeakFrequencyBin) * (16000.0 / 2.0 / 1024.0 / 64.0)
					var frequencyBand FrequencyBand
					switch {
					case frequencyHz > 250 && frequencyHz < 520:
						frequencyBand = _250_520
					case frequencyHz >= 520 && frequencyHz < 1450:
						frequencyBand = _520_1450
					case frequencyHz >= 1450 && frequencyHz < 3500:
						frequencyBand = _1450_3500
					case frequencyHz >= 3500 && frequencyHz <= 5500:
						frequencyBand = _3500_5500
					default:
						continue
					}
					if frequencyBand == _3500_5500 {
						continue
					}
					s.signature.FrequencyBandToSoundPeaks[frequencyBand] = append(
						s.signature.FrequencyBandToSoundPeaks[frequencyBand],
						FrequencyPeak{
							FFTPassNumber:             fftPassNumber,
							PeakMagnitude:             peakMagnitude,
							CorrectedPeakFrequencyBin: correctedPeakFrequencyBin,
							SampleRateHz:              sampleRate,
						})
				}
			}
		}
	}
}
