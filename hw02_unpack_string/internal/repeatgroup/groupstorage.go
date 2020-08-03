package repeatgroup

import (
	"bytes"
	"fmt"
)

// Unpacker can create string from Unpacker, assembly string from internal view.
type Unpacker interface {
	Unpack() (string, error)
}

// GroupStorage contain all groups.
type GroupStorage struct {
	rgs []repeatGroup
}

// Unpack converts internal state to string.
func (gs GroupStorage) Unpack() (string, error) {
	var result bytes.Buffer
	for _, gr := range gs.rgs {
		if gr.repeatCnt > 0 {
			newChunk := bytes.Repeat(gr.buffer, gr.repeatCnt)
			writedLen, err := result.Write(newChunk)
			if err != nil {
				return "", fmt.Errorf("unable to write chunk %q into buffer, %w", string(newChunk), err)
			}
			if writedLen < len(gr.buffer) {
				return "", fmt.Errorf("buffer writes partially")
			}
		}
	}

	return result.String(), nil
}

// AddRepeatGroup new entry to storage.
func (gs *GroupStorage) addRepeatGroup(rg repeatGroup) (newSize int) {
	gs.rgs = append(gs.rgs, rg)

	return len(gs.rgs)
}

// Add new RepeatGroup entry to storage created it from butes (copy).
func (gs *GroupStorage) add(b []byte, cnt int) (newSize int) {
	chunk := make([]byte, len(b))
	copy(chunk, b)
	gs.addRepeatGroup(repeatGroup{chunk, cnt})

	return len(gs.rgs)
}

// flushBuffer create repeatGroup from buffer and cntm add them to GroupStorage and flush buffer.
func (gs *GroupStorage) flushBuffer(buffer *bytes.Buffer, cnt int) {
	if buffer.Len() > 0 {
		gs.add(buffer.Bytes(), cnt)
		buffer.Reset()
	}
}
