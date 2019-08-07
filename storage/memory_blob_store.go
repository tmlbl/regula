package storage

import (
	"bytes"

	"github.com/google/uuid"
)

// MemoryBlobStore is an in-memory implementation of BlobStore for testing
type MemoryBlobStore struct {
	uploads map[string]*bytes.Buffer
	blobs   map[string][]byte
}

// NewMemoryBlobStore initializes a MemoryBlobStore
func NewMemoryBlobStore() *MemoryBlobStore {
	return &MemoryBlobStore{
		uploads: map[string]*bytes.Buffer{},
		blobs:   map[string][]byte{},
	}
}

// Exists checks whether a blob with the given sum exists
func (s *MemoryBlobStore) Exists(sum string) (bool, int) {
	if data, ok := s.blobs[sum]; ok {
		return true, len(data)
	}
	return false, 0
}

func (s *MemoryBlobStore) StartUpload() string {
	id := uuid.New().String()
	s.uploads[id] = bytes.NewBuffer([]byte{})
	return id
}

func (s *MemoryBlobStore) UploadRange(id string, data []byte, start, end int) error {
	s.uploads[id].Grow(len(data))
	s.uploads[id].Write(data)
	return nil
}
