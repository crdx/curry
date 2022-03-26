package cache

import (
	"os"
	"path"
)

type Cache struct {
	filePath string
	dirPath  string
}

func New(filePath string) Cache {
	return Cache{
		filePath: filePath,
		dirPath:  path.Dir(filePath),
	}
}

func (self Cache) exists() bool {
	_, err := os.Stat(self.filePath)
	return err == nil
}

func (self Cache) Delete() bool {
	if !self.exists() {
		return false
	}
	return os.Remove(self.filePath) == nil
}

func (self Cache) ReadBytes(cb func() []byte) (data []byte, err error) {
	if !self.exists() {
		data = cb()
		return
	}

	data, err = os.ReadFile(self.filePath)
	return
}

func (self Cache) WriteBytes(bytes []byte) (err error) {
	if self.exists() {
		return
	}

	err = os.MkdirAll(self.dirPath, 0o777)
	if err != nil {
		return
	}

	err = os.WriteFile(self.filePath, bytes, 0o666)
	return
}
