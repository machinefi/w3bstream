package filesystem

type FileSystemOp interface {
	Upload(key string, file []byte) error
	Read(key string) ([]byte, error)
	Delete(key string) error
}
