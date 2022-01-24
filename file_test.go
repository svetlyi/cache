package cache_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/svetlyi/cache"
)

var tmpFileName = "github_com_svetlyi_file_cache_test"

func TestGet(t *testing.T) {
	fileForCache := filepath.Join(os.TempDir(), tmpFileName)

	f, err := cache.NewFile(fileForCache)
	defer os.Remove(fileForCache)

	if err != nil {
		t.Errorf("could not create file cache: %v", err)
	}
	defer f.Close()

	item := f.Get("test")
	if item.IsHit() {
		t.Error("cache must be empty")
	}
}

func TestHas(t *testing.T) {
	fileForCache := filepath.Join(os.TempDir(), tmpFileName)

	f, err := cache.NewFile(fileForCache)
	defer os.Remove(fileForCache)

	if err != nil {
		t.Errorf("could not create file cache: %v", err)
	}
	defer f.Close()

	if f.Has("test_has") {
		t.Error("cache must be empty")
	}
}

func TestDelete(t *testing.T) {
	fileForCache := filepath.Join(os.TempDir(), tmpFileName)

	f, err := cache.NewFile(fileForCache)
	defer os.Remove(fileForCache)

	if err != nil {
		t.Errorf("could not create file cache: %v", err)
	}
	defer f.Close()

	if err := f.Delete("test_has"); err != nil {
		t.Errorf("error while removing a non existing item: %v", err)
	}
}

func TestSave(t *testing.T) {
	fileForCache := filepath.Join(os.TempDir(), tmpFileName)

	f, err := cache.NewFile(fileForCache)
	defer os.Remove(fileForCache)

	if err != nil {
		t.Errorf("could not create file cache: %v", err)
	}
	defer f.Close()

	newItem := cache.CacheItem{
		Key:       "new_item",
		Value:     "test value",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := f.Save(newItem); err != nil {
		t.Errorf("error saving an item: %v", err)
	}

	storedItem := f.Get("new_item")
	if !storedItem.IsHit() {
		t.Error("could not find an item")
	}

	if storedItem.Value != "test value" {
		t.Errorf("wrong value: %s", storedItem.Value)
	}

	newItem2 := cache.CacheItem{
		Key:       "new_item2",
		Value:     "test value 2",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	expiredItem := cache.CacheItem{
		Key:       "new_item3",
		Value:     "test value 3",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := f.Save(newItem2); err != nil {
		t.Errorf("error saving an item: %v", err)
	}

	if err := f.Save(expiredItem); err != nil {
		t.Errorf("error saving an item: %v", err)
	}

	if f.Has("new_item3") {
		t.Error("new_item3 must be expired")
	}

	if !f.Has("new_item2") {
		t.Error("new_item2 must be present")
	}

	if !f.Has("new_item") {
		t.Error("new_item must still be present")
	}

	// trying to close and then open to make sure, that the data saved
	if err := f.Close(); err != nil {
		t.Fatalf("error closing file cache: %v", err)
	}

	f2, err := cache.NewFile(fileForCache)
	if err != nil {
		t.Fatalf("error creating cache again: %v", err)
	}
	defer f2.Close()

	if !f2.Has("new_item2") {
		t.Error("new_item2 must still be present")
	}
}

func TestOpenExisting(t *testing.T) {
	fileForCache, err := filepath.Abs("test.json")

	if err != nil {
		t.Fatalf("could not get absolute path for file %s: %v", "test.json", err)
	}
	f, err := cache.NewFile(fileForCache)

	if err != nil {
		t.Errorf("could not open file cache: %v", err)
	}
	defer f.Close()

	if !f.Has("test_key_1") {
		t.Error("test_key_1 must be present")
	}
}
