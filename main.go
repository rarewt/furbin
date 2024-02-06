package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	microphone "rarewt/furbin/microphone"
	record "rarewt/furbin/record"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/live"
	"github.com/joho/godotenv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketMessage struct {
	Type string `json:"type"`
}

func main() {
	// Init microphone
	microphone.Initialize()

	// init deepgram lib
	client.InitWithDefault()

	// Configuration for the Deepgram client
	ctx := context.Background()

	// load .env props
	errenv := godotenv.Load()
	if errenv != nil {
		log.Fatal(".env file could not be loaded")
	}

	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	log.Println("apiKey: " + apiKey)

	// client + transcript options
	clientOptions := interfaces.ClientOptions{
		// EnableKeepAlive: true,
	}
	transcriptOptions := interfaces.LiveTranscriptionOptions{
		Language:    "en-US",
		Model:       "nova-2-conversationalai",
		Punctuate:   true,
		Encoding:    "linear16",
		Channels:    1,
		SampleRate:  16000,
		SmartFormat: true,
		// To get UtteranceEnd, the following must be set:
		InterimResults: true,
		UtteranceEndMs: "1000",
	}

	// Callback used to handle responses from Deepgram
	callback := record.DeepgramCallback{}

	// Create a new Deepgram LiveTranscription client with config options
	dgClient, err := client.New(ctx, apiKey, &clientOptions, transcriptOptions, callback)
	if err != nil {
		fmt.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// mic stuff
	mic, err := microphone.New(microphone.AudioConfig{
		InputChannels: 1,
		SamplingRate:  16000,
	})
	if err != nil {
		fmt.Printf("Mic initialisation failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start the mic
	err = mic.Start()
	if err != nil {
		fmt.Printf("mic.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	go func() {
		// feed the microphone stream to the deepgram client (this is a blocking call)
		mic.Stream(dgClient)
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close mic stream
	err = mic.Stop()
	if err != nil {
		fmt.Printf("mic.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}

	// teardown library
	microphone.Teardown()

	// close deepgram client
	dgClient.Stop()

	fmt.Printf("Exiting program\n")

}
