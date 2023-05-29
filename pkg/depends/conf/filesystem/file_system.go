package filesystem

type FileSystemOp interface {
	Upload(key string, file []byte) error
	UploadWithMD5(key, md5 string, file []byte) error
	Read(key string) ([]byte, error)
	ReadWithMD5(key, md5 string) ([]byte, error)
	Delete(key string) error
}
