package summarization

import "sync"

type CircularBuffer struct {
	groupID         int64
	threadID        int
	pinnedMessageID int
	buffer          []string
	size            int
	writePos        int
	count           int
	mutex           sync.Mutex
}

func NewCircularBuffer(size int, groupID int64, threadID, pinnedMessageID int) *CircularBuffer {
	return &CircularBuffer{
		groupID:         groupID,
		threadID:        threadID,
		pinnedMessageID: pinnedMessageID,
		buffer:          make([]string, size),
		size:            size,
		writePos:        0,
		count:           0,
	}
}

func (cb *CircularBuffer) Add(msg string) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.buffer[cb.writePos] = msg

	cb.writePos = (cb.writePos + 1) % cb.size

	if cb.count < cb.size {
		cb.count++
	}
}

func (cb *CircularBuffer) GetAll() []string {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	result := make([]string, cb.count)

	// Copy in correct order (oldest to newest)
	readPos := cb.writePos - cb.count
	if readPos < 0 {
		readPos += cb.size
	}

	for i := 0; i < cb.count; i++ {
		result[i] = cb.buffer[(readPos+i)%cb.size]
	}

	return result
}

func (cb *CircularBuffer) Clear() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.count = 0
}
