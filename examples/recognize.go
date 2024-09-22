package examples

import (
	"context"
	"fmt"
	"goshazam"
	"log"
)

func exampleRecognize() {
	jsonResult, err := goshazam.Recognize(context.Background(), "test.mp3")
	if err != nil {
		log.Fatalf("Error recognizing audio: %v", err)
	}
	fmt.Println(string(jsonResult))
}
