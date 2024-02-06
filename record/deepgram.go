package record

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

// Implement the api.Callback interface
type DeepgramCallback struct {
	socket *websocket.Conn
}

var bufferedMessage = ""

func (c DeepgramCallback) Message(mr *api.MessageResponse) error {
	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)
	if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
		return nil
	}
	fmt.Printf("\nDeepgram: %s, IsFinal: %s, SpeechFinal: %s\n\n", sentence, strconv.FormatBool(mr.IsFinal), strconv.FormatBool(mr.SpeechFinal))

	if bufferedMessage += sentence + " "; mr.SpeechFinal {
		fmt.Printf("\n\nFull message: %s\n\n", bufferedMessage)
		bufferedMessage = ""
	}

	c.socket.WriteJSON(sentence)
	return nil
}

func (c DeepgramCallback) Metadata(md *api.MetadataResponse) error {
	fmt.Printf("\n[Metadata] Received\n")
	fmt.Printf("Metadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	fmt.Printf("Metadata.Channels: %d\n", md.Channels)
	fmt.Printf("Metadata.Created: %s\n\n", strings.TrimSpace(md.Created))
	return nil
}

func (c DeepgramCallback) UtteranceEnd(ur *api.UtteranceEndResponse) error {
	fmt.Printf("\n[UtteranceEnd] Received\n")
	return nil
}

func (c DeepgramCallback) Error(er *api.ErrorResponse) error {
	fmt.Printf("\n[Error] Received\n")
	fmt.Printf("Error.Type: %s\n", er.Type)
	fmt.Printf("Error.Message: %s\n", er.Message)
	fmt.Printf("Error.Description: %s\n\n", er.Description)
	return nil
}
