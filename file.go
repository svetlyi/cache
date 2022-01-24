package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type file struct {
	sync.Mutex
	fileCache   os.File
	memoryCache Cache
}

func NewFile(path string) (Cache, error) {
	fileCache, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0700)

	if err != nil {
		return nil, fmt.Errorf("error opening/creating file: %v", err)
	}

	log.Println("opened file", path)

	jsonBytes, err := ioutil.ReadAll(fileCache)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	memoryCache := NewMemoryCache()

	if len(jsonBytes) != 0 {
		if err := memoryCache.Unserialize(jsonBytes); err != nil {
			return nil, fmt.Errorf("error decoding file: %v", err)
		}
	}

	return &file{
		fileCache:   *fileCache,
		memoryCache: memoryCache,
	}, nil
}

func (f *file) Get(key string) CacheItem {
	return f.memoryCache.Get(key)
}

func (f *file) Has(key string) bool {
	return f.memoryCache.Has(key)
}

func (f *file) Delete(key string) error {
	return f.memoryCache.Delete(key)
}

func (f *file) Save(item CacheItem) error {
	f.Lock()
	defer f.Unlock()

	if err := f.memoryCache.Save(item); err != nil {
		return err
	}

	return f.flush()
}

func (f *file) Close() error {
	if err := f.memoryCache.Close(); err != nil {
		return err
	}

	return f.fileCache.Close()
}

func (f *file) Serialize() ([]byte, error) {
	return f.memoryCache.Serialize()
}

func (f *file) Unserialize(data []byte) error {
	return f.memoryCache.Unserialize(data)
}

func (f *file) flush() error {
	bytes, err := f.memoryCache.Serialize()
	if err != nil {
		return fmt.Errorf("could not serialize cache items: %v", err)
	}
	if err := f.fileCache.Truncate(0); err != nil {
		return fmt.Errorf("error truncating file: %v", err)
	}
	if _, err := f.fileCache.Seek(0, 0); err != nil {
		return fmt.Errorf("error positioning file: %v", err)
	}

	if _, err := f.fileCache.Write(bytes); err != nil {
		return fmt.Errorf("error while saving file to %s: %v", f.fileCache.Name(), err)
	}

	return nil
}
