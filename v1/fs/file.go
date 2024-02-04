package fs

type File interface {
	GetSize() float64
	GetPath() string
}
