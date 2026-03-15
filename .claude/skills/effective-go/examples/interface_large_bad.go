package storage

type FileManager interface {
	Read() error
	Write() error
	Delete() error
	List() error
	Copy() error
	Move() error
}
