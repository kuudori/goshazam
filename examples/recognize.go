package examples

import (
	"context"
	"fmt"
	"github.com/kuudori/goshazam"
	"log"
)

func testExample() {
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
