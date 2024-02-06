package interfaces

import "io"

// Microphone defines a interface for a Microphone implementation
type Microphone interface {
	Start() error
	Read() ([]int16, error)
	Stream(w io.Writer) error
	Mute()
	Unmute()
	Stop() error
}
