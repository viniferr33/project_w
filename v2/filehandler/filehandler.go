package filehandler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

var InvalidFilepathIsDir = errors.New("file cannot be a directory")

type File struct {
	Name      string
	Filepath  string
	Extension string
	Size      int64
}

func GetFileFromFilepath(f string) (*File, error) {
	fileInfo, err := os.Stat(f)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, InvalidFilepathIsDir
	}

	return &File{
		Name:      fileInfo.Name(),
		Filepath:  f,
		Extension: filepath.Ext(f),
		Size:      fileInfo.Size(),
	}, nil
}

func (f *File) Copy(d string) (*File, error) {
	source, err := os.Open(f.Filepath)
	if err != nil {
		return nil, err
	}
	defer source.Close()

	destinationFilePath := fmt.Sprintf("%s/%s%s", d, uuid.New().String(), f.Extension)
	destination, err := os.Create(destinationFilePath)
	if err != nil {
		return nil, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return nil, err
	}

	return GetFileFromFilepath(destinationFilePath)
}
