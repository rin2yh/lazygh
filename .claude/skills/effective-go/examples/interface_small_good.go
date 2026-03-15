package storage

type Reader interface {
	Read(p []byte) (int, error)
}

type File struct{}

func (f *File) Read(p []byte) (int, error) {
	return 0, nil
}
