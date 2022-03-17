package cache

import (
	"log"
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

func (self Cache) ReadBytes(cb func() []byte) []byte {
	if !self.exists() {
		return cb()
	}

	data, err := os.ReadFile(self.filePath)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (self Cache) WriteBytes(bytes []byte) {
	if self.exists() {
		return
	}

	err := os.MkdirAll(self.dirPath, 0o777)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(self.filePath, bytes, 0o666)
	if err != nil {
		log.Fatal(err)
	}
}
