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
    // Recognize a song from an audio file
    jsonResult, err := goshazam.Recognize(context.Background(), "test.mp3")
    if err != nil {
        log.Fatalf("Error recognizing audio: %v", err)
    }
    fmt.Println(string(jsonResult))
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