# A simple file/memory cache in Golang

It's inspired by PSR-6 and was designed for small pet projects, that need some simple cache solution.

# Installation

```
go get github.com/svetlyi/cache
```

Add the package to dependencies, run for example `go build` and it will download the package.

# Usage

Basically all the cases are covered in tests. An example of usage:

```
package main

import (
	"log"
	"time"

	"github.com/svetlyi/cache"
)

func main() {
	fileForCache := "file.cache"

	f, err := cache.NewFile(fileForCache)
	if err != nil {
		log.Fatalf("could not create file cache: %v", err)
	}
	defer f.Close()

	newItem := cache.CacheItem{
		Key:       "new_item",
		Value:     "test value",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := f.Save(newItem); err != nil {
		log.Fatalf("error saving an item: %v", err)
	}

	storedItem := f.Get("new_item")
	if !storedItem.IsHit() {
		log.Fatalf("could not find an item")
	}

	log.Println("cache value:", storedItem.Value)
}
```
