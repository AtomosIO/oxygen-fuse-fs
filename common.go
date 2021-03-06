package oxygenfuse

import (
	"bazil.org/fuse"

	"crypto/rand"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
)

const (
	MIN_PASSWORD_LENGTH = 6
)

func Debug(msg fmt.Stringer) {
	fmt.Println(msg.String())
}

// An interruptChannel is a channel that signals that a request has been interrupted.
// Being able to receive from the channel means the request has been
// interrupted.
type interruptChannel chan struct{}

func (interruptChannel) String() string { return "fuse.Intr" }

type requestsInterruptMap struct {
	sync.RWMutex
	requestInterrupts map[fuse.RequestID]interruptChannel
}

func (requestsInterruptMap *requestsInterruptMap) Get(id fuse.RequestID) interruptChannel {
	//TODO: Implement

	return nil
}

func (requestsInterruptMap *requestsInterruptMap) Set(requestId fuse.RequestID, channel interruptChannel) {
	requestsInterruptMap.Lock()
	requestsInterruptMap.requestInterrupts[requestId] = channel
	requestsInterruptMap.Unlock()
}

func (requestsInterruptMap *requestsInterruptMap) Delete(requestId fuse.RequestID) {
	requestsInterruptMap.Lock()
	delete(requestsInterruptMap.requestInterrupts, requestId)
	requestsInterruptMap.Unlock()
}

func NewRequestInterruptsMap() *requestsInterruptMap {
	return &requestsInterruptMap{
		requestInterrupts: make(map[fuse.RequestID]interruptChannel),
	}

}

func opName(req fuse.Request) string {
	t := reflect.Indirect(reflect.ValueOf(req)).Type()
	s := t.Name()
	s = strings.TrimSuffix(s, "Request")
	return s
}

func NewEmptyReader() io.Reader {
	return &EmptyReader{}
}

type EmptyReader struct{}

func (reader *EmptyReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func RandomString(length int) string {
	return string(RandomAlphaNumericByteSlice(length))
}

func RandomProjectName() string {
	return RandomString(5)
}

func RandomAlphaNumericByteSlice(length int) []byte {
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = ALPHANUMERIC[b%byte(len(ALPHANUMERIC))]
	}
	return bytes
}

func RandomUsername() string {
	return string(RandomAlphaNumericByteSlice(5))
}

func RandomEmail() string {
	// Really hackie. Use a StringBuilder construct.
	output := append(RandomAlphaNumericByteSlice(5), byte('@'))
	output = append(output, RandomAlphaNumericByteSlice(5)...)
	output = append(output, []byte(".com")...)
	return string(output)
}

func RandomPassword() string {
	return string(RandomAlphaNumericByteSlice(MIN_PASSWORD_LENGTH))
}
