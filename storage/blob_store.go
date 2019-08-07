package storage

// BlobStore defines the methods that should be exposed by a storage
// implementation
type BlobStore interface {
	Exists(sum string) (exists bool, size int)
	StartUpload() (id string)
	UploadRange(id string, data []byte, start, end int) error
}
