# GoShazam

GoShazam is a Go library that provides an interface to the Shazam music recognition service. It allows you to recognize songs from audio files using Shazam's powerful audio fingerprinting technology.

## Features

- Recognize songs from audio files (MP3, WAV, OGG)
- Generate audio fingerprints
- Interface with Shazam's API
- Return full JSON response for maximum flexibility

## Installation

To install GoShazam, use `go get`:

```
go get github.com/kuudori/goshazam
```

## Usage

Here's a basic example of how to use GoShazam:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kuudori/goshazam"
)

func main() {
	client := goshazam.NewShazamClient()
	result, err := client.Recognize(context.Background(), "test.mp3")

	if err != nil {
		log.Fatalf("Error recognizing audio: %v", err)
	}
	rawResult := result.Raw() // Get raw JSON response
	fmt.Println(string(rawResult))

	serializedResult, _ := result.Serialize() // Serialize response
	fmt.Println(serializedResult)

	// keep it short

	// client.Recognize(context.Background(), "test.mp3").Serialize()
	// client.Recognize(context.Background(), "test.mp3").Raw()
}
```

## Examples

For more detailed examples, please check the `examples` folder in the repository.

## Supported Audio Formats

GoShazam supports recognizing songs from the following audio formats:

- MP3
- WAV
- OGG

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This library is for educational and personal use only. Make sure you comply with Shazam's terms of service when using this library.