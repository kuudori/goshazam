package goshazam

const (
	sampleRate     = 16000
	maxTimeSeconds = 6
	fftSize        = 2048
	fftOutputSize  = fftSize/2 + 1
	numFFTs        = 256
)

const (
	_250_520 FrequencyBand = iota
	_520_1450
	_1450_3500
	_3500_5500
)
